package models

import (
	"github.com/speedland/wcg"
	"testing"
)

func TestParseEpgJson(t *testing.T) {
	assert := wcg.NewAssert(t)
	list, err := ParseEpgJsonString(testJson)
	assert.Nil(err, "ParseEpgrecJson should not return error")

	p := list[0]
	assert.EqInt(22529, p.EventId, "EventId")
	assert.EqStr("ＡＮＮニュース・あすの空もよう【字】", p.Title, "Title")
	assert.EqStr(
		"正確なニュース・情報をいち早くお伝えするＡＮＮニュース！テレビ朝日系列の放送局２６局が総力をあげ、緻密な取材にもとづいたニュースを最新機材を駆使して放送します。",
		string(p.Detail),
		"Detail",
	)
	assert.EqInt(3, len(p.Categories), "Length of Categories")
	assert.EqStr("定時・総合", p.Categories[0].Middle, "Categories[0].Middle")
	assert.EqStr("天気", p.Categories[1].Middle, "Categories[0].Middle")
	assert.EqStr("ローカル・地域", p.Categories[2].Middle, "Categories[0].Middle")
	assert.EqStr("ニュース／報道", p.Categories[0].Large, "Categories[0].Large")
	assert.EqStr("ニュース／報道", p.Categories[1].Large, "Categories[1].Large")
	assert.EqStr("ニュース／報道", p.Categories[2].Large, "Categories[2].Large")
}

var testJson = `
[{
  "programs": [{
    "event_id": 22529,
    "freeCA": false,
    "audio": [{"extdesc": "","langcode": "jpn","type": "ステレオ"}],
    "video": {"aspect": "16:9","resolution": "HD"},
    "attachinfo": [],
    "channel": "GR5_1064",
    "title": "ＡＮＮニュース・あすの空もよう【字】",
    "detail": "正確なニュース・情報をいち早くお伝えするＡＮＮニュース！テレビ朝日系列の放送局２６局が総力をあげ、緻密な取材にもとづいたニュースを最新機材を駆使して放送します。",
    "extdetail": [],
    "start": 14010548400000,
    "end": 14010550800000,
    "duration": 240,
    "category": [{
        "middle": {"en": "Regular/General","ja_JP": "定時・総合"},
        "large": {"en": "news","ja_JP": "ニュース／報道"}
      },{
        "middle": {"en": "Weather","ja_JP": "天気"},
        "large": {"en": "news","ja_JP": "ニュース／報道"}
      },{
        "middle": {"en": "Local","ja_JP": "ローカル・地域"},
        "large": {"en": "news","ja_JP": "ニュース／報道"}
      }
    ],
    "extdetail": [
          {
            "item": "ヨーロッパ最西端、大西洋に面する国・ポルトガル。首都・リスボンを起点に南へ、アルガルヴェ地方の港町・ファロや大航海時代の要塞が残るサグレスを目指します。",
            "item_description": "◇番組内容"
          }
    ]
  }]
}]
`
