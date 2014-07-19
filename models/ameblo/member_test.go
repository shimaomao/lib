package ameblo

import (
	"github.com/speedland/wcg"
	"testing"
	"time"
)

func TestIsMentionedIn(t *testing.T) {
	assert := wcg.NewAssert(t)
	member := &Member{
		"foo",
		"http://example.com/picture.png",
		"http://example.com/blog/",
		[]string{"test", "テスト", "てすと"},
		123456,
		1,
		time.Time{},
	}
	assert.Ok(member.IsMentionedIn("This is a test."), "English word")
	assert.Ok(member.IsMentionedIn("This is a テスト."), "Katakana word")
	assert.Ok(member.IsMentionedIn("This is a てすと."), "Hiragana word")
	assert.Ok(!member.IsMentionedIn("This is a."), "Not mentioned")
}
