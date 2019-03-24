package wx

import (
	"encoding/xml"
	"strconv"
	"time"
)

type Base struct {
	FromUserName CDATAText
	ToUserName   CDATAText
	MsgType      CDATAText
	CreateTime   CDATAText
}

type CDATAText struct {
	Text string `xml:",innerxml"`
}

type TextMessage struct {
	XMLName xml.Name `xml:"xml"`
	Base
	Content CDATAText
}

func (b *Base) InitBaseData(w *WeixinClient, msgType string) {
	b.FromUserName = value2CDATA(w.Message["ToUserName"].(string))
	b.ToUserName = value2CDATA(w.Message["FromUserName"].(string))
	b.CreateTime = value2CDATA(strconv.FormatInt(time.Now().Unix(), 10))
	b.MsgType = value2CDATA(msgType)
}
