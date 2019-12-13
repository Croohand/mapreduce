package httputil

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	httpClient ClientWithLogging
	avChecker  ClientWithLogging
)

type ClientWithLogging struct {
	SenderName string
	*http.Client
}

func NewClient(name string) ClientWithLogging {
	avChecker = ClientWithLogging{name, &http.Client{Timeout: time.Duration(100 * time.Millisecond)}}
	httpClient = ClientWithLogging{name, &http.Client{Timeout: time.Duration(5 * time.Second)}}
	return ClientWithLogging{name, &http.Client{Timeout: time.Duration(5 * time.Second)}}
}

func (c ClientWithLogging) Do(r *http.Request) (*http.Response, error) {
	if c.SenderName != "" {
		r.Header.Set("Sender-Name", c.SenderName)
	}
	return c.Client.Do(r)
}

func (c ClientWithLogging) Get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c ClientWithLogging) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return c.Do(req)
}

func (c ClientWithLogging) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	return c.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}

type DefaultMuxWithLogging struct {
	SelfName   string
	LoggerAddr string
}

func (m DefaultMuxWithLogging) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if m.LoggerAddr != "" {
		sender := r.Header.Get("Sender-Name")
		if sender == "" {
			sender = "unknown"
		}
		e := fmt.Sprintf("%v %v %v %v", m.SelfName, r.URL.Path, sender, r.ContentLength)
		httpClient.PostForm(m.LoggerAddr+"/LogEntry", url.Values{"Entry": {e}})
	}
	http.DefaultServeMux.ServeHTTP(w, r)
}
