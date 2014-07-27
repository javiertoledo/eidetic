package main

import (
	"os"
	"log"
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"encoding/json"
	"crypto/sha1"
	"menteslibres.net/gosexy/redis"
	"encoding/base64"
)

type SerializableResponse struct {
    Header    map[string][]string   // the HTTP response headers
    Body      []byte                // the body string
}

func main() {
	// The remote URL has to be provided as the first command line argument
	remote_host := os.Args[1]
	log.Println("Initializing proxy for: " + remote_host)

	remote, err := url.Parse(remote_host)
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)

	redisClient := redis.New()
	err = redisClient.Connect("127.0.0.1", uint(6379))
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", handler(proxy, redisClient))
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func hashKey(r *http.Request) string {
	url := []byte(r.URL.String())
	// log.Println("\n\n\nRequest URL: " + r.URL.String())

	// Read the body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	// log.Println("\nRequest Body:\n" + string(body))

	// Hack to pretend the reader has ben resetted
	reader := bytes.NewReader(body)
	r.Body = ioutil.NopCloser(reader)

	hash := sha1.Sum(append(url,body...))
	return "eidetic#" + base64.URLEncoding.EncodeToString(hash[:])
}

func handler(p *httputil.ReverseProxy, redisClient *redis.Client) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		hash := hashKey(r)

		serializableResponse := SerializableResponse{make(map[string][]string), nil}

		s, err := redisClient.Get(hash)
		if err != nil { // The request is not cached
			rec := httptest.NewRecorder()
			log.Println("Non cached: " + hash)

			// Perform the real request and cache it
			p.ServeHTTP(rec, r)

			for k, v := range rec.Header() {
				serializableResponse.Header[k] = v
			}
			serializableResponse.Body = rec.Body.Bytes()

			jsonResponse, err := json.Marshal(serializableResponse)
			if err != nil {
				panic(err)
			}
			redisClient.Set(hash, jsonResponse)
			w.Header().Set("X-Eidetic", "Live request")

		} else { // The request is cached
			log.Println("Cached!: " + hash)

			// Load the cached request
			err = json.Unmarshal([]byte(s), &serializableResponse)
			if err != nil {
				panic(err)
			}

			w.Header().Set("X-Eidetic", "Cached request")

		}

		//Copy the data to the actual ResponseWriter
		// log.Println("\n\n\nResponse Headers:")
		for k, v := range serializableResponse.Header {
			w.Header()[k] = v
			// log.Println(k + ": ")
			// for _, str := range v {
			// 	log.Println("  " + str)
			// }
		}
		w.Write([]byte(serializableResponse.Body))
		// log.Println("\nResponse body:\n" + string(serializableResponse.Body))

	}
}
