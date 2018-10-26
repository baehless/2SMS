package main

import (
	"net/http"
	"log"
	"io/ioutil"
	"fmt"
	"errors"
	"strings"
	"strconv"
)

type scraperProxyHandler struct {
	httpsClient		*http.Client
	scionClient		*SCIONClient
}

func (sph *scraperProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var resp *http.Response
	var err error
	// Try tunneling over SCION
	resp, err = sph.scionClient.TunnelRequest(r)
	if err != nil {
		// If that doesn't work use HTTPS
		log.Println("Target unreachable via SCION:", err, "Fallback to HTTPS.")
		ip := strings.Split(r.URL.Host, ":")[0]
		httpPort, err := strconv.Atoi(strings.Split(r.URL.Host, ":")[1])
		if err != nil {
			log.Println("Cannot compute HTTPS port from SCION one:", err)
			w.WriteHeader(500)
			return
		}
		host := ip + ":" + fmt.Sprint(httpPort + 1)
		if r.Method == http.MethodGet {
			// Make scrape GET request
			resp, err = sph.httpsClient.Get("https://" + host + r.URL.Path)
		} else if r.Method == http.MethodPost {
			resp, err = sph.httpsClient.Post("https://"+ host +r.URL.Path, "application/x-protobuf", r.Body)
		} else {
			err = errors.New("Unsupported method: " + r.Method)
		}
		if err != nil {
			log.Println("Failed processing HTTPS request:", err)
			w.WriteHeader(400)
			return
		}
	}

	// Copy response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading response body:", err)
		return
	}
	w.Write(body)
	defer resp.Body.Close()
	// Copy headers
	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}

	// Print proxy traffic infos
	//log.Printf("%v\n", resp.Status)
	//log.Println(string(body))
}