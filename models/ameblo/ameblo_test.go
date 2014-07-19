package ameblo

import (
	"github.com/speedland/wcg"
	"net/http"
	"os"
	"testing"
)

func TestCrawlEntryList(t *testing.T) {
	assert := wcg.NewAssert(t)
	c := &Crawler{
		http.DefaultClient,
	}
	list, err := c.CrawlEntryList("http://ameblo.jp/sayumimichishige-blog/entrylist.html")
	assert.Nil(err, "CrawlEntryList should not return an error")
	assert.Ok(len(list) > 0, "CrawlEntryList should return some entries.")
	assert.Ok(list[0].Title != "", "An entry in EntryList should have title.")
}

func TestCrawlEntry(t *testing.T) {
	assert := wcg.NewAssert(t)
	c := &Crawler{
		http.DefaultClient,
	}
	entry, err := c.CrawlEntry("http://ameblo.jp/sayumimichishige-blog/entry-11874881676.html")
	assert.Nil(err, "CrawlEntryList should not return an error")
	assert.EqStr("春ツアー思い出", entry.Title, "An entry in EntryList should have title.")
}

func TestParseEntry(t *testing.T) {
	assert := wcg.NewAssert(t)
	file, _ := os.Open("./entry-sample.html")
	defer file.Close()
	entry, err := ParseEntry(file)
	assert.Nil(err, "ParseEntryContent should not return an error")
	assert.EqStr("早くヴァンプになりた〜い！工藤 遥", entry.Title, "Title")
}

func TestParseEntryList(t *testing.T) {
	assert := wcg.NewAssert(t)
	file, _ := os.Open("./entrylist-sample.html")
	defer file.Close()
	entryList, err := ParseEntryList(file)
	assert.Nil(err, "ParseEntryContent should not return an error")
	assert.EqInt(20, len(entryList), "# of entries")
	assert.EqStr("http://ameblo.jp/morningmusume-10ki/entry-11874661134.html", entryList[1].Url, "Url")
	assert.EqStr("早くヴァンプになりた〜い！工藤 遥", entryList[1].Title, "Title")
	assert.EqInt(126, entryList[1].AmComments, "AmComments")
	assert.EqInt(551, entryList[1].AmLikes, "AmLikes")
}
