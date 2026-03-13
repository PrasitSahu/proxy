package handler

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	conf "github.com/PrasitSahu/proxy/internal"
	"github.com/PrasitSahu/proxy/internal/api"
)

// RootFunc handles the root request
func RootFunc(res http.ResponseWriter, req *http.Request) {
	urlStr := req.URL.Query().Get("url")
	if len(strings.TrimSpace(urlStr)) == 0 {
		http.Error(res, api.ErrNoURL.Error(), http.StatusBadRequest)
		return
	}

	url, err := url.Parse(urlStr)
	if err != nil || url.Host == "" || url.Scheme == "" {
		http.Error(res, api.ErrInvalidURL.Error(), http.StatusBadRequest)
		return
	}

	newRequest, err := http.NewRequestWithContext(
		req.Context(),
		req.Method,
		url.String(),
		req.Body,
	)
	if err != nil {
		http.Error(res, api.ErrReqFail.Error(), http.StatusInternalServerError)
		return
	}

	newRequest.Header = req.Header.Clone()

	resp, err := conf.Config.HttpClient.Do(newRequest)
	if err != nil {
		log.Println(err)
		http.Error(res, api.ErrReqFail.Error(), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	for k, v := range resp.Header {
		for _, vv := range v {
			res.Header().Add(k, vv)
		}
	}

	res.WriteHeader(resp.StatusCode)
	io.Copy(res, resp.Body)
}