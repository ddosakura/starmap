package main

import (
	"log"
	"net/http"
	"os"

	"github.com/kr/pretty"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		pretty.Println(r.Form, r.PostForm, r.URL.Query(), r.URL.Query().Get("a"))
		w.Write([]byte{})
	})
	go func() {
		if e := http.ListenAndServe(":8888", nil); e != nil {
			log.Println(e)
			os.Exit(-1)
		}
	}()

	http.Get("http://localhost:8888?a=1&a=2")

	// Result
	_ = `
	url.Values{} url.Values{} url.Values{
    "a": {"1", "2"},
	} 1
	`
}
