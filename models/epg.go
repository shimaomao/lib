package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/speedland/lib/util"
	"io"
	"strings"
	"time"
)

type Epg struct {
	EventId    int             `json:"evenrt_id"`
	Title      string          `json:"title"`
	Detail     util.ByteString `json:"detail"`
	StartAt    time.Time       `json:"start_at"`
	EndAt      time.Time       `json:"end_at"`
	Cid        string          `json:"cid"`
	Sid        string          `json:"sid"`
	Categories []Category      `json:"category"`
}

type Category struct {
	Middle string `json:"middle"`
	Large  string `json:"large"`
}

func ParseEpgJsonString(jsonstr string) ([]*Epg, []error) {
	return ParseEpgJson(bytes.NewBuffer([]byte(jsonstr)))
}

func ParseEpgJson(jsonio io.Reader) ([]*Epg, []error) {
	var v [](map[string]interface{})
	err := json.NewDecoder(jsonio).Decode(&v)
	if err != nil {
		return nil, []error{err}
	}
	programs := make([]*Epg, 0)
	elist := make([]error, 0)
	for _, v := range v[0]["programs"].([]interface{}) {
		if epg := newEpgFromMap(v.(map[string]interface{}), &elist); epg != nil {
			programs = append(programs, epg)
		}
	}
	if len(elist) > 0 {
		return programs, elist
	} else {
		return programs, nil
	}
}

type ErrEpgParseFailed struct {
	err    error
	Source map[string]interface{}
}

func (e *ErrEpgParseFailed) Error() string {
	return e.err.Error()
}

func newEpgFromMap(m map[string]interface{}, elist *[]error) *Epg {
	var epg *Epg
	defer func() {
		if v := recover(); v != nil {
			*elist = append(*elist, &ErrEpgParseFailed{
				err: fmt.Errorf("%v: %v", v, m),
			})
		}
	}()
	epg = new(Epg)
	epg.EventId = int(m["event_id"].(float64))
	epg.Cid, epg.Sid = parseChannel(m["channel"].(string))
	epg.Title = m["title"].(string)
	epg.Detail = []byte(m["detail"].(string))
	epg.StartAt = time.Unix(int64(m["start"].(float64))/10000, 0)
	epg.EndAt = time.Unix(int64(m["end"].(float64))/10000, 0)
	epg.Categories = make([]Category, 0)
	for _, v := range m["category"].([]interface{}) {
		c := Category{}
		if middle, ok := v.(map[string]interface{})["middle"]; ok {
			c.Middle = middle.(map[string]interface{})["ja_JP"].(string)
		}
		if large, ok := v.(map[string]interface{})["large"]; ok {
			c.Large = large.(map[string]interface{})["ja_JP"].(string)
		}
		epg.Categories = append(epg.Categories, c)
	}
	return epg
}

func parseChannel(channel string) (string, string) {
	s := strings.Split(channel, "_")
	return s[0], s[1]
}
