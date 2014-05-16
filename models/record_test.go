package models

import (
	"fmt"
	"github.com/speedland/lib/util"
	"github.com/speedland/wcg"
	"testing"
	"time"
)

type DummyRecord struct {
	key       string
	startAt   time.Time
	endAt     time.Time
	stopCh    chan bool
	done      bool
	isRunning bool
}

func NewDummyRecord(key string) *DummyRecord {
	return &DummyRecord{
		key:    key,
		stopCh: make(chan bool, 1),
	}
}

func (dr *DummyRecord) Key() string {
	return dr.key
}

func (dr *DummyRecord) StartAt() time.Time {
	return dr.startAt
}

func (dr *DummyRecord) EndAt() time.Time {
	return dr.endAt
}

func (dr *DummyRecord) String() string {
	return fmt.Sprintf("<Dummy %s>", dr.Key())
}

func (dr *DummyRecord) IsRunning() bool {
	return dr.isRunning
}

func (dr *DummyRecord) Start() error {
	dr.isRunning = true
	go func() {
		<-dr.stopCh
		dr.isRunning = false
		dr.done = true
	}()
	return nil
}

func (dr *DummyRecord) Stop() error {
	dr.stopCh <- true
	return nil
}

func (dr *DummyRecord) CheckInterval() time.Duration {
	return time.Duration(1 * time.Millisecond)
}

func TestRecorderStart(t *testing.T) {
	assert := wcg.NewAssert(t)
	receiver := make(chan []Record)
	recorder := NewRecorder((<-chan []Record)(receiver))

	interval := time.Duration(30) * time.Millisecond
	wait := 10
	go recorder.Start(interval)

	now := time.Now()
	r1 := NewDummyRecord("r1") // cancel before starting
	r1.startAt = now
	r1.endAt = now.Add(time.Duration(5 * time.Second))
	receiver <- []Record{r1}
	err := util.WaitFor(func() bool {
		return len(recorder.controls) == 1
	}, wait)
	assert.Nil(err, "check r1 in controls.")
	ctrl := recorder.controls[r1.Key()]

	err = util.WaitFor(func() bool {
		return r1.done == true
	}, wait)
	assert.Nil(err, "check r1 has been done.")
	assert.EqInt(int(RSSucceeded), int(ctrl.state), "Record state should be RSSucceeded")
	recorder.Stop()
}

func TestRecorderStart_CancelBeforeStart(t *testing.T) {
	assert := wcg.NewAssert(t)
	receiver := make(chan []Record)
	recorder := NewRecorder((<-chan []Record)(receiver))

	interval := time.Duration(30) * time.Millisecond
	wait := 10
	go recorder.Start(interval)

	now := time.Now()
	r1 := NewDummyRecord("r1") // cancel before starting
	r1.startAt = now.Add(time.Duration(30 * time.Minute))
	r1.endAt = now.Add(time.Duration(60 * time.Minute))
	// r2 := NewDummyRecord("r2") // cancel after starting
	// r3 := NewDummyRecord("r3") // succeeded
	// r4 := NewDummyRecord("r4") // failed

	receiver <- []Record{r1}
	err := util.WaitFor(func() bool {
		return len(recorder.controls) == 1
	}, wait)
	assert.Nil(err, "check r1 in controls.")
	assert.EqStr(r1.Key(), recorder.controls[r1.Key()].record.Key(), "check r1 in controls.")

	receiver <- []Record{}
	err = util.WaitFor(func() bool {
		return len(recorder.controls) == 0
	}, wait)
	assert.Nil(err, "check r1 removed from controls.")

	recorder.Stop()
}

func TestRecorderStart_CancelWhileRecording(t *testing.T) {
	assert := wcg.NewAssert(t)
	receiver := make(chan []Record)
	recorder := NewRecorder((<-chan []Record)(receiver))
	interval := time.Duration(30) * time.Millisecond
	wait := 5

	go recorder.Start(interval)

	now := time.Now()
	r1 := NewDummyRecord("r1") // cancel before starting
	r1.startAt = now
	r1.endAt = now.Add(time.Duration(30 * time.Minute))

	receiver <- []Record{r1}
	err := util.WaitFor(func() bool {
		return len(recorder.controls) == 1
	}, wait)
	assert.Nil(err, "check r1 in controls.")
	assert.EqStr(r1.Key(), recorder.controls[r1.Key()].record.Key(), "check r1 in controls.")

	err = util.WaitFor(func() bool {
		return recorder.controls[r1.Key()].state == RSRecording
	}, wait)
	assert.Nil(err, "check r1 in RSRecording state.")

	receiver <- []Record{}
	err = util.WaitFor(func() bool {
		return len(recorder.controls) == 0
	}, wait)
	assert.Nil(err, "check r1 removed from controls.")

	recorder.Stop()
}

func genTestRecord() *TvRecord {
	start := time.Now()
	end := start.Add(50 * time.Minute)
	return NewTvRecord(
		"title", "category",
		start, end,
		"20", "hd", "me",
	)
}

func TestRecordValidator(t *testing.T) {
	var err error
	assert := wcg.NewAssert(t)
	r := genTestRecord()

	r.Title = ""
	err = RecordValidator.Eval(r)
	assert.NotNil(err, "Empty Title Error")

	r.Title = "../foo"
	err = RecordValidator.Eval(r)
	assert.NotNil(err, "Special Character Validation")

	r.Title = "NormalTitle"
	err = RecordValidator.Eval(r)
	assert.Nil(err, "Normal Character Validation")

	r.Title = "日本語はつかえる"
	err = RecordValidator.Eval(r)
	assert.Nil(err, "日本語 Validation")

	r.Title = "モーニング娘。"
	err = RecordValidator.Eval(r)
	assert.Nil(err, "日本語 特殊文字 Valiation")
}
