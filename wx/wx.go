package wx

import (
	"net/http"
	"errors"
	"sort"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"github.com/clbanning/mxj"
	"log"
	"encoding/xml"
)

type weixinQuery struct {
	Signature    string `json:"signature"`
	Timestamp    string `json:"timestamp"`
	Nonce        string `json:"nonce"`
	EncryptType  string `json:"encrypt_type"`
	MsgSignature string `json:"msg_signature"`
	Echostr      string `json:"echostr"`
}

type WeixinClient struct {
	Token          string
	Query          weixinQuery
	Message        map[string]interface{}
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	Methods        map[string]func() bool
}

func (client *WeixinClient) initWeixinQuery() {
	var q weixinQuery
	q.Nonce = client.Request.URL.Query().Get("nonce")
	q.Echostr = client.Request.URL.Query().Get("echostr")
	q.Signature = client.Request.URL.Query().Get("signature")
	q.Timestamp = client.Request.URL.Query().Get("timestamp")
	q.EncryptType = client.Request.URL.Query().Get("encrypt_type")
	q.MsgSignature = client.Request.URL.Query().Get("msg_signature")
	client.Query = q
}

func (client *WeixinClient) signature() string {
	strs := sort.StringSlice{client.Token, client.Query.Timestamp, client.Query.Nonce}
	sort.Strings(strs)
	str := ""
	for _, s := range strs {
		str += s
	}
	h := sha1.New()
	h.Write([]byte(str))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func NewClient(r *http.Request, w http.ResponseWriter, token string) (*WeixinClient, error) {
	weixinClient := new(WeixinClient)
	weixinClient.Token = token
	weixinClient.Request = r
	weixinClient.ResponseWriter = w

	weixinClient.initWeixinQuery()

	if weixinClient.Query.Signature != weixinClient.signature() {
		return nil, errors.New("invalid signature")
	}
	return weixinClient, nil
}

func (client *WeixinClient) initMessage() error {
	body, err := ioutil.ReadAll(client.Request.Body)
	if err != nil {
		return err
	}

	m, err := mxj.NewMapXml(body)
	if err != nil {
		return err
	}

	if _, ok := m["xml"]; !ok {
		return errors.New("invalid message")
	}

	message, ok := m["xml"].(map[string]interface{})
	if !ok {
		return errors.New("invalid field `xml` type")
	}

	client.Message = message
	log.Println(client.Message)
	return nil
}

func (client *WeixinClient) text() {
	inMsg, ok := client.Message["Content"].(string)
	if !ok {
		return
	}

	var reply TextMessage
	reply.InitBaseData(client, "text")
	reply.Content = value2CDATA(fmt.Sprintf("receiveMsg: %s", inMsg))

	replyXml, err := xml.Marshal(reply)
	if err != nil {
		log.Println(err)
		client.ResponseWriter.WriteHeader(403)
		return
	}

	client.ResponseWriter.Header().Set("Content-Type", "text/xml")
	client.ResponseWriter.Write(replyXml)
}

func (client *WeixinClient) Run() {
	err := client.initMessage()
	if err != nil {
		log.Println(err)
		client.ResponseWriter.WriteHeader(403)
		return
	}

	MsgType, ok := client.Message["MsgType"].(string)
	if !ok {
		client.ResponseWriter.WriteHeader(403)
		return
	}

	switch MsgType {
	case "text":
		client.text()
		break
	default:
		break
	}
}
