package models

import (
	"fmt"
	"github.com/speedland/lib/util"
	"github.com/speedland/wcg"
	v "github.com/speedland/wcg/validation"
	"time"
)

var MaxRecordTime = 24 * time.Hour
var MinRecordTime = 1 * time.Minute
var RecordValidator = v.NewObjectValidator()
var InvalidChars = "[\x00-\x1F\x22-\x27\x2a-\x2c\x2f]"
var ValidUUID = "^[a-z0-9]{8}-[a-z0-9]{4}-[a-z0-9]{4}-[a-z0-9]{4}-[a-z0-9]{12}$"

type Record interface {
	Key() string
	Start() error
	Stop() error
	IsRunning() bool
	StartAt() time.Time
	EndAt() time.Time
	CheckInterval() time.Duration
}

type RecordState int

var (
	RSWaiting   = RecordState(1)
	RSRecording = RecordState(2)
	RSCanceled  = RecordState(3)
	RSSucceeded = RecordState(4)
	RSFailed    = RecordState(5)
)

type RecordCtrl struct {
	record        Record
	state         RecordState
	stopCh        chan bool
	logger        wcg.Logger
	checkInterval time.Duration
}

func NewRecordCtrl(rec Record) *RecordCtrl {
	return &RecordCtrl{
		record: rec,
		state:  RSWaiting,
		stopCh: make(chan bool, 1),
		logger: util.GetLogger(),
	}
}

func (ctrl *RecordCtrl) String() string {
	return fmt.Sprintf("[%v / %d]", ctrl.record, ctrl.state)
}

func (ctrl *RecordCtrl) Start() {
	go func() {
		ctrl.control()
		c := time.Tick(ctrl.record.CheckInterval())
		for {
			select {
			case <-ctrl.stopCh:
				ctrl.record.Stop()
				ctrl.state = RSCanceled
				break
			case <-c:
				ctrl.control()
				break
				// default:
				// 	ctrl.logger.Info("Controled %v", ctrl)
				// 	break
			}

			if ctrl.state == RSSucceeded || ctrl.state == RSFailed || ctrl.state == RSCanceled {
				break
			}
		}
	}()
}

func (ctrl *RecordCtrl) Cancel() {
	ctrl.stopCh <- true
}

func (ctrl *RecordCtrl) control() {
	now := time.Now()
	record := ctrl.record
	switch ctrl.state {
	case RSWaiting:
		if now.After(record.StartAt()) && now.Before(record.EndAt()) {
			err := record.Start()
			if err == nil {
				ctrl.state = RSRecording
				ctrl.logger.Info("%v: Recording started.", record)
			} else {
				ctrl.logger.Error("%v: Could not start recording - %v", record, err)
			}
		} else if now.After(record.EndAt()) {
			// invalid.
			ctrl.state = RSCanceled
		}
		break
	case RSRecording:
		if now.After(record.StartAt()) && now.Before(record.EndAt()) {
			if !record.IsRunning() {
				err := record.Start()
				if err == nil {
					ctrl.logger.Info("%v: Recording restarted.", record)
				} else {
					ctrl.logger.Error("%v: Could not restart recording - %v", record, err)
				}
			}
		} else if now.After(record.EndAt()) {
			if record.IsRunning() {
				ctrl.record.Stop()
				ctrl.state = RSSucceeded
				ctrl.logger.Info("%v: Recording stopped.", record)
			} else {
				ctrl.logger.Warn("%v: Recording is not running, marked as failure.", record)
				ctrl.state = RSFailed
			}
		}
		break
	default:
		// nothing for RSSucceeded, RSFailed, and RSCanceled
		break
	}
}

type RecordResult struct {
	record Record
	err    error
}

var ErrCanceled = fmt.Errorf("Canceled")

// Recorder is a control component to manage all recorder objects.
type Recorder struct {
	receiver   <-chan []Record
	recorder   chan *RecordResult
	controls   map[string]*RecordCtrl
	iSec       int
	shouldStop bool
	logger     wcg.Logger
	// stats
	waiting   int
	recording int
	succeeded int
	canceled  int
	failed    int
}

// Create a new Recorder instance.
// receiver should be a channel to register the record
func NewRecorder(receiver <-chan []Record) *Recorder {
	return &Recorder{
		receiver:   receiver,
		recorder:   make(chan *RecordResult, 128),
		controls:   make(map[string]*RecordCtrl),
		shouldStop: false,
		logger:     util.GetLogger(),
	}
}

func (r *Recorder) Start(interval time.Duration) {
	c := time.Tick(interval)
	for !r.shouldStop {
		select {
		case records := <-r.receiver:
			r.merge(records)
			break
		case <-c:
			break
		}
		r.updateStats()
	}
}

func (r *Recorder) Stop() {
	r.shouldStop = true
}

func (r *Recorder) Upcomming() time.Time {
	var min time.Time
	for _, v := range r.controls {
		t := v.record.StartAt()
		if v.state == RSWaiting || v.state == RSRecording {
			if min.IsZero() {
				min = t
			} else {
				if t.Before(min) {
					min = t
				}
			}
		}
	}
	return min
}

func (r *Recorder) merge(newrecords []Record) {
	// convert newrecords to map
	now := time.Now()
	newmap := make(map[string]Record)
	for _, rec := range newrecords {
		if rec.EndAt().After(now) {
			newmap[rec.Key()] = rec
		}
	}
	// compare existing records with newrecords
	if len(r.controls) > 0 {
		r.logger.Trace("Check %d record(s) under controls.", len(r.controls))
		for key, ctrl := range r.controls {
			r.logger.Trace("Check %s key in controls.", key)
			if _, ok := newmap[key]; ok {
				r.logger.Debug("%s already exists, skipping.", key)
				delete(newmap, key)
			} else {
				// no longer exists so cancel.
				r.logger.Debug("%s is no longer exists, canceled.", key)
				ctrl.Cancel()
			}
		}
	}

	if len(newmap) > 0 {
		r.logger.Info("New %d record(s) gets under controls.", len(newmap))
		for key, rec := range newmap {
			ctrl := NewRecordCtrl(rec)
			ctrl.Start()
			r.controls[key] = ctrl
		}
	}
}

func (r *Recorder) updateStats() {
	removals := make([]string, 0)
	var waiting, recording, succeeded, canceled, failed int
	for key, rctrl := range r.controls {
		switch rctrl.state {
		case RSWaiting:
			waiting += 1
			break
		case RSRecording:
			recording += 1
			break
		case RSSucceeded:
			removals = append(removals, key)
			succeeded += 1
		case RSFailed:
			removals = append(removals, key)
			failed += 1
		case RSCanceled:
			removals = append(removals, key)
			canceled += 1
			break
		}
	}
	for _, key := range removals {
		delete(r.controls, key)
	}
	r.waiting = waiting
	r.recording = recording
	r.succeeded = r.succeeded + succeeded
	r.failed = r.failed + failed
	r.canceled = r.canceled + canceled
}

type TvRecord struct {
	Id        string    `json:"id"`
	Title     string    `json:"title"`
	Category  string    `json:"category"` // Category Name
	StartAt   time.Time `json:"start_at"`
	EndAt     time.Time `json:"end_at"`
	Cid       string    `json:"cid"`       // Channel ID
	Sid       string    `json:"sid"`       // Signal ID
	Uid       string    `json:"uid"`       // user id
	InputIdx  int       `json:"input_idx"` // PT2 input channel index
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewTvRecord(title string, category string, start time.Time, end time.Time,
	cid string, sid string, uid string) *TvRecord {
	r := &TvRecord{
		Id:       wcg.Must(wcg.UUID()).(string),
		Title:    title,
		Category: category,
		StartAt:  start,
		EndAt:    end,
		Cid:      cid,
		Sid:      sid,
		Uid:      uid,
	}
	now := time.Now()
	r.CreatedAt = now
	r.UpdatedAt = now
	return r
}

func (r *TvRecord) Key() string {
	return r.Id
}

func (r *TvRecord) String() string {
	return "Record:" + r.Key()
}

func (r *TvRecord) RecordTime() time.Duration {
	return r.EndAt.Sub(r.StartAt)
}

func init() {
	RecordValidator.Field("Id").Required().Match(ValidUUID)
	RecordValidator.Field("Title").Required().
		Unmatch(InvalidChars).Min(1).Max(64)
	RecordValidator.Field("Category").Required().
		Unmatch(InvalidChars).Min(1).Max(64)

	RecordValidator.Field("StartAt").Required()
	RecordValidator.Field("EndAt").Required()
	RecordValidator.Field("Cid").Required()
	RecordValidator.Field("Sid").Required()
	RecordValidator.Field("Uid").Required()
	RecordValidator.Func(func(r interface{}) *v.FieldValidationError {
		rec := r.(*TvRecord)
		rt := rec.RecordTime()
		if rt < MinRecordTime {
			return v.NewFieldValidationError("録画時間が短すぎます。", nil)
		}
		if rt > MaxRecordTime {
			return v.NewFieldValidationError("録画時間が長すぎます。", nil)
		}
		return nil
	})
	// Note:
	//   we don't validate overlaps since it should be noticed by
	//   other way and fixed by manually.
}
