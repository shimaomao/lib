package tv

import (
	"bufio"
	"code.google.com/p/go-html-transform/css/selector"
	"code.google.com/p/go.net/html"
	"code.google.com/p/go.text/encoding/japanese"
	"code.google.com/p/go.text/transform"
	"fmt"
	"github.com/speedland/lib/util"
	"github.com/speedland/wcg"
	v "github.com/speedland/wcg/validation"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type IEpg struct {
	Id           string          `json:"id"`
	StationId    string          `json:"station_id"`
	StationName  string          `json:"station_name"`
	ProgramTitle string          `json:"program_title"`
	ProgramId    int             `json:"program_id"`
	Body         util.ByteString `json:"detail"`
	StartAt      time.Time       `json:"start_at"`
	EndAt        time.Time       `json:"end_at"`
	Category     string          `json:"category"`
	Cid          string          `json:"cid"`
	Sid          string          `json:"sid"`
	Optout       bool            `json:"optout"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

func (iepg *IEpg) ToTvRecord() *TvRecord {
	return &TvRecord{
		Id:        wcg.Must(wcg.UUID()).(string),
		Title:     strings.Replace(iepg.ProgramTitle, "/", "／", -1),
		Category:  iepg.Category,
		StartAt:   iepg.StartAt,
		EndAt:     iepg.EndAt,
		Cid:       iepg.Cid,
		Sid:       iepg.Sid,
		Uid:       "", // for future use.
		IEpgId:    iepg.Id,
		CreatedAt: iepg.CreatedAt,
		UpdatedAt: iepg.UpdatedAt,
	}
}

type Crawler struct {
	client *http.Client
}

func NewCrawler(client *http.Client) *Crawler {
	return &Crawler{
		client,
	}
}

type CrawlerConfig struct {
	Keyword   string    `json:"keyword"`
	Category  string    `json:"category"`
	Scope     int       `json:"scope"`
	CreatedAt time.Time `json:"created_at"`
}

var CrawlerConfigValidator = v.NewObjectValidator()

func init() {
	CrawlerConfigValidator.Field("Keyword").Required()
	CrawlerConfigValidator.Field("Category").Required()
}

const FEED_SCOPE_ALL = 0
const FEED_SCOPE_TERRESTRIAL = 1
const FEED_SCOPE_BS = 2
const FEED_SCOPE_CS = 5
const FEED_SCOPE_CS_PREMIUM = 4
const feed_url_template = "http://tv.so-net.ne.jp/rss/schedulesBySearch.action?stationPlatformId=%d&condition.keyword=%s"
const iepg_url_template = "http://tv.so-net.ne.jp/iepg.tvpid?id=%s"

func (c *Crawler) GetIEpgList(keyword string, scope int) ([]string, error) {
	urlstr := fmt.Sprintf(feed_url_template, scope, url.QueryEscape(keyword))
	resp, err := c.client.Get(urlstr)
	if err != nil {
		return nil, fmt.Errorf("HTTP Error: %v (url = %q)", err, urlstr)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Non-200 code (%d) was returned from %q", resp.StatusCode, urlstr)
	}
	list, err := ParseRss(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Could not parse RSS from %q: %v", urlstr, err)
	}
	return list, nil
}

func (c *Crawler) GetIEpg(id string) (*IEpg, error) {
	urlstr := fmt.Sprintf(iepg_url_template, id)
	resp, err := c.client.Get(urlstr)
	if err != nil {
		return nil, fmt.Errorf("HTTP Error: %v (url = %q)", err, urlstr)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Non-200 code (%d) was returned from %q", resp.StatusCode, urlstr)
	}
	iepg, err := ParseIEpg(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Could not parse RSS from %q: %v", urlstr, err)
	}
	iepg.Id = id
	return iepg, nil
}

const datetime_format = "%s-%s-%s %s:00 +0900"
const datetime_layout = "2006-01-02 15:04:05 -0700"

func ParseIEpg(r io.Reader) (*IEpg, error) {
	var err error
	iepg := &IEpg{}
	iepg.Body = make([]byte, 0)
	rio := transform.NewReader(r, japanese.ShiftJIS.NewDecoder())
	scanner := bufio.NewScanner(rio)
	kv := make(map[string]string)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			for scanner.Scan() {
				line = scanner.Text()
				iepg.Body = append(iepg.Body, []byte(line)...)
			}
			break
		}
		tmp := strings.Split(line, ": ")
		if len(tmp) != 2 {
			fmt.Printf("%q", line)
			continue
		}
		kv[tmp[0]] = tmp[1]
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("IEPG scan error: %v", err)
	}
	iepg.StationId = kv["station"]
	iepg.StationName = kv["station-name"]
	iepg.ProgramTitle = kv["program-title"]
	if kv["program-id"] != "" {
		if iepg.ProgramId, err = strconv.Atoi(kv["program-id"]); err != nil {
			return nil, fmt.Errorf("Could not parse `program-id` - %q in %v: %v ", kv["program-id"], kv, err)
		}
	}
	start_time := fmt.Sprintf(datetime_format, kv["year"], kv["month"], kv["date"], kv["start"])
	if iepg.StartAt, err = time.Parse(datetime_layout, start_time); err != nil {
		return nil, fmt.Errorf("Could not parse start time - %q in %v: %v", start_time, kv, err)
	}
	end_time := fmt.Sprintf(datetime_format, kv["year"], kv["month"], kv["date"], kv["end"])
	if iepg.EndAt, err = time.Parse(datetime_layout, end_time); err != nil {
		return nil, fmt.Errorf("Could not parse end time - %q in %v: %v", end_time, kv, err)
	}
	return iepg, nil
}

// Parse the RSS feed and returns iEPG Ids.
func ParseRss(r io.Reader) ([]string, error) {
	root, err := html.Parse(r)
	if err != nil {
		return nil, fmt.Errorf("Could not parse RSS feed: %v", err)
	}
	s, _ := selector.Selector("item")
	nodes := s.Find(root)
	list := []string{}
	for _, n := range nodes {
		for i := range n.Attr {
			if n.Attr[i].Key == "rdf:about" {
				id := extractIdFromUrl(n.Attr[i].Val)
				if id != "" {
					list = append(list, id)
				}
			}
		}
	}
	return list, nil
}

func extractIdFromUrl(urlstr string) string {
	tmp := strings.Split(urlstr, "/")
	if len(tmp) == 0 {
		return ""
	}
	tmp = strings.Split(tmp[len(tmp)-1], ".")
	if len(tmp) == 0 {
		return ""
	}
	return tmp[0]
}
