package main

import (
	"log"
	"net/http"
	"time"
	"fmt"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
func writeLog(r *http.Request, t time.Time, match string, pattern string) {
	if logLevel != "prod" {
		d := time.Now().Sub(t)
		l := fmt.Sprintf("[ACCESS] | %-10s | %-40s | %-16s | %-10s | %-40s |", r.Method, r.URL.Path, d.String(), match, pattern)
		log.Println(l)
	}
}
