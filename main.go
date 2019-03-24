package main

import (
	"net/http"
	"github.com/selge/wechat/wx"
	"log"
	"fmt"
	"time"
	"flag"
)

const (
	logLevel = "dev"
)

func get(w http.ResponseWriter, r *http.Request) {
	client, err := wx.NewClient(r, w, *token)
	if err != nil {
		log.Println(err)
		w.WriteHeader(403)
		return
	}

	if len(client.Query.Echostr) > 0 {
		w.Write([]byte(client.Query.Echostr))
		return
	}

	w.WriteHeader(403)
	return
}

func post(w http.ResponseWriter, r *http.Request) {
	client, err := wx.NewClient(r, w, *token)
	if err != nil {
		log.Println(err)
		w.WriteHeader(403)
		return
	}
	client.Run()
	return
}

var port = flag.Int("port", 9099, "http listen port")
var token = flag.String("token", "123456", "token")

func main() {
	flag.Parse()
	server := http.Server{
		Addr:           fmt.Sprintf(":%d", *port),
		Handler:        &httpHandler{},
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 0,
	}

	log.Println(fmt.Sprintf("Listen: %d", *port))
	log.Fatal(server.ListenAndServe())
}
