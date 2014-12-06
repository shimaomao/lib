package tv

import (
	"github.com/speedland/wcg"
	"net/http"
	"os"
	"testing"
)

func TestGetIEpgList(t *testing.T) {
	assert := wcg.NewAssert(t)
	client := NewCrawler(http.DefaultClient)
	list, err := client.GetIEpgList("今井絵理子", FEED_SCOPE_ALL)
	assert.Nil(err, "GetIEpgList should not return an error.")
	assert.NotNil(list, "GetIEpgList should return list of ids")
}

func TestGetIEpg(t *testing.T) {
	t.Skipf("TestGetIEpg needs http server mock since the real side would remoe iEPG link after the program is broadcasted.")
	// assert := wcg.NewAssert(t)
	// client := NewCrawler(http.DefaultClient)
	// list, err := client.GetIEpg("200171201412080100")
	// assert.Nil(err, "GetIEpg should not return an error.")
	// assert.NotNil(list, "GetIEpg should return list of ids")
}

func TestParseRss(t *testing.T) {
	assert := wcg.NewAssert(t)
	file, _ := os.Open("./iepg-feed-sample.html")
	defer file.Close()
	list, err := ParseRss(file)
	assert.Nil(err, "ParseRss should not return an error.")
	assert.EqStr("101072201412050100", list[0], "Id Match")
	assert.EqStr("101056201412052300", list[1], "Id Match")
	assert.EqStr("200171201412080100", list[2], "Id Match")
	assert.EqStr("400639201412081030", list[3], "Id Match")
	assert.EqStr("101040201412110214", list[4], "Id Match")
	assert.EqStr("400639201412120930", list[5], "Id Match")
}

func TestParseIEpg(t *testing.T) {
	assert := wcg.NewAssert(t)
	file, _ := os.Open("./iepg-sample.iepg")
	defer file.Close()
	iepg, err := ParseIEpg(file)
	assert.Nil(err, "ParseIEPG should not return an error")
	assert.EqStr("The　Girls　Live　▽道重さゆみ卒業ライブに密着▽LoVendoЯスタジオライブ", iepg.ProgramTitle, "ProgramTitle")
	assert.EqStr("テレビ東京", iepg.StationName, "StationName")
	assert.EqStr("DFS00430", iepg.StationId, "StationId")
	assert.EqInt(7852, iepg.ProgramId, "ProgramId")
	assert.EqInt(2014, iepg.StartAt.Year(), "StartAt.Year")
	assert.EqInt(12, int(iepg.StartAt.Month()), "StartAt.Date")
	assert.EqInt(5, iepg.StartAt.Day(), "StartAt.Date")
	assert.EqInt(1, iepg.StartAt.Hour(), "StartAt.Hour")
	assert.EqInt(0, iepg.StartAt.Minute(), "StartAt.Mininute")
	zone, _ := iepg.StartAt.Zone()
	assert.EqStr("JST", zone, "StartAt.Zone")

	assert.EqInt(2014, iepg.EndAt.Year(), "EndAt.Year")
	assert.EqInt(12, int(iepg.EndAt.Month()), "EndAt.Date")
	assert.EqInt(5, iepg.EndAt.Day(), "EndAt.Date")
	assert.EqInt(1, iepg.EndAt.Hour(), "EndAt.Hour")
	assert.EqInt(30, iepg.EndAt.Minute(), "EndAt.Mininute")
	zone, _ = iepg.EndAt.Zone()
	assert.EqStr("JST", zone, "EndAt.Zone")

}
