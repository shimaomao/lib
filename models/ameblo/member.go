package ameblo

import (
	"strings"
	"time"
)

type Member struct {
	Name       string    `json:"name"`
	PictureUrl string    `json:"profile_url"`
	BlogUrl    string    `json:"blog_url"`
	Nicknames  []string  `json:"nicknames"`
	Color      int       `json:"color"` // RGB
	Generation int       `json:"generation"`
	Birthday   time.Time `json:"birthday"`
}

func (m *Member) IsMentionedIn(content string) bool {
	for _, name := range m.Nicknames {
		if strings.Index(content, name) >= 0 {
			return true
		}
	}
	return false
}

func (m *Member) ValidateEntryListUrl(url string) bool {
	return strings.HasPrefix(url, strings.Split(m.BlogUrl, ".html")[0])
}
