package mm

import (
	"code.google.com/p/go-html-transform/css/selector"
	"code.google.com/p/go-html-transform/h5"
	"code.google.com/p/go.net/html"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const TIME_FORMAT = "2006-01-02 15:04:05"

var numRegexp = regexp.MustCompile("\\d+")
var JST *time.Location

func init() {
	JST, _ = time.LoadLocation("Asia/Tokyo")
}

type Crawler struct {
	client *http.Client
}

func (c *Crawler) CrawlEntryList(url string) ([]*AmebloEntry, error) {
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Server returns %d code.", resp.StatusCode)
	}

	return ParseEntryList(resp.Body)
}

func (c *Crawler) CrawlEntry(url string) (*AmebloEntry, error) {
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Server returns %d code.", resp.StatusCode)
	}

	return ParseEntry(resp.Body)
}

type AmebloEntry struct {
	Url        string    `json:"url"`
	Title      string    `json:"title"`
	Owner      string    `json:"owner"`
	PostAt     time.Time `json:"post_at"`
	Content    string    `json:"content"`
	CrawledAt  string    `json:"crawled_at"`
	AmLikes    int       `json:"am_likes"`
	AmComments int       `json:"am_comments"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func ParseEntry(r io.Reader) (*AmebloEntry, error) {
	root, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	s, _ := selector.Selector(".articleText")
	nodes := s.Find(root)
	if len(nodes) == 0 {
		return nil, nil
	}
	content := h5.RenderNodesToString(nodes)

	s, _ = selector.Selector("title")
	nodes = s.Find(root)
	if len(nodes) == 0 {
		return nil, nil
	}
	title := extractText(nodes[0].FirstChild)

	entry := &AmebloEntry{
		Title:   strings.Split(title, "｜")[0],
		Content: content,
	}
	return entry, nil
}

func ParseEntryList(r io.Reader) ([]*AmebloEntry, error) {
	root, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	s, _ := selector.Selector("ul.contentsList li")
	nodes := s.Find(root)
	entryList := make([]*AmebloEntry, 0)

	for _, listItem := range nodes {
		e := &AmebloEntry{}
		// title & url
		n := findOne("a.contentTitle", listItem)
		e.Title = extractText(n.FirstChild)
		e.Url = getAttributeValue("href", n)
		// postAt
		n = findOne(".contentTime time", listItem)
		e.PostAt, err = time.ParseInLocation(TIME_FORMAT, extractText(n.FirstChild), JST)
		if err != nil {
			continue
		}
		// AmLikes and AmComments
		n = findOne(".contentComment", listItem)
		e.AmComments, _ = strconv.Atoi(
			numRegexp.FindString(extractText(n.FirstChild)),
		)
		n = findOne("a.skinWeakColor", n.Parent)
		e.AmLikes, _ = strconv.Atoi(
			numRegexp.FindString(extractText(n.FirstChild)),
		)
		entryList = append(entryList, e)
	}
	return entryList, nil
}

func extractText(n *html.Node) string {
	return h5.RenderNodesToString([]*html.Node{n})
}

func findOne(sel string, node *html.Node) *html.Node {
	s, _ := selector.Selector(sel)
	n := s.Find(node)
	return n[0]
}

func getAttributeValue(key string, node *html.Node) string {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}
