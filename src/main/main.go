package main

import (
	"net/http"
	"flag"
	"log"
	"errors"
	"sync"
	"net/url"
	"io/ioutil"
	"strings"
)

var (
	listen = flag.String("listen", ":80", "HTTP listen address.")
	path = flag.String("path", "/xhprof.io/import", "Import endpoint.")
	xhProfIoUrl = flag.String("xhProfIoUrl", "http://localhost/xhprof.io", "Url to xhprof.io installation.")
)

type importRequest struct {
	HttpHost string
	RequestMethod string
	RequestUri string
	XHProfData string
}

type importHandler interface {
	// interface to consume requests
	consume(*importRequest, *sync.WaitGroup)
}

type importHandle struct {
	// import queue
	queue chan importRequest
	// channel that receive one close event
	closeSignal chan string

	count int
}

func main() {
	flag.Parse()

	handle := new(importHandle)
	handle.count = 0
	handle.queue = make(chan importRequest)
	var wg *sync.WaitGroup = new(sync.WaitGroup)
	go importHandlerLoop(handle, wg)

	// provide http endpoint
	http.HandleFunc(*path, func(w http.ResponseWriter, r *http.Request) {
		importRequest, err := createImportRequest(r)
		if err == nil {
			handle.queue <- importRequest
			handle.count++
			w.WriteHeader(http.StatusAccepted)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Import error: %v", err)
	})

	// startup http server
	err := http.ListenAndServe(*listen, nil)
	if err != nil {
		panic("HTTP server err: " + err.Error())
	}

}

func importHandlerLoop(handle *importHandle, wg *sync.WaitGroup) {
	for {
		select {
		case importRequest := <-handle.queue:
			wg.Add(1)
			handle.consume(importRequest, wg)
		}
	}
}

func (this importHandle) consume(r importRequest, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Printf("Executing import request for '%s %s%s'.", r.RequestMethod, r.HttpHost, r.RequestUri)

	values := make(url.Values)
	values.Set("http_host", r.HttpHost)
	values.Set("request_method", r.RequestMethod)
	values.Set("request_uri", r.RequestUri)
	values.Set("xhprof_data", r.XHProfData)
	url := []string{*xhProfIoUrl, "/import.php"};
	resp, err := http.PostForm(strings.Join(url, ""), values)
	if err != nil {
		log.Printf("ERROR: %s", err)
		return
	}

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	log.Printf("%s", body)
}


func createImportRequest(r *http.Request) (importRequest importRequest, err error) {
	httpHost := r.FormValue("http_host")
	requestMethod := r.FormValue("request_method")
	requestUri := r.FormValue("request_uri")
	xhProfData := r.FormValue("xhprof_data")

	if httpHost != "" && requestMethod != "" && requestUri != "" && xhProfData != "" {
		importRequest.HttpHost = httpHost
		importRequest.RequestMethod = requestMethod
		importRequest.RequestUri = requestUri
		importRequest.XHProfData = xhProfData
	} else {
		err = errors.New("Missing host, method, uri or xhprof data.")
	}

	return importRequest, err
}
