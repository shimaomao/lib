package mm

import (
	"strings"
)

type Member struct {
	Name       string   `json:"name"`
	PictureUrl string   `json:"profile_url"`
	BlogUrl    string   `json:"blog_url"`
	Nicknames  []string `json:'nicknames'`
}

func (m *Member) IsMentionedIn(content string) bool {
	for _, name := range m.Nicknames {
		if strings.Index(content, name) >= 0 {
			return true
		}
	}
	return false
}
