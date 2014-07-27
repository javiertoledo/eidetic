package main

import(
    "os"
    "log"
    "bytes"
    "net/url"
    "net/http"
    "net/http/httputil"
    "net/http/httptest"
    //"encoding/json"
    "crypto/sha1"
)

func main() {
    // The remote URL has to be provided as the first command line argument
    remote_host := os.Args[1]
    log.Println("Initializing proxy for: " + remote_host)

    remote, err := url.Parse(remote_host)
    if err != nil {
        panic(err)
    }

    proxy := httputil.NewSingleHostReverseProxy(remote)
    http.HandleFunc("/", handler(proxy))
    err = http.ListenAndServe(":8080", nil)
    if err != nil {
        panic(err)
    }
}

func hashKey(r *http.Request) [20]byte {
    url := []byte(r.URL.String())

    buf := new(bytes.Buffer)
    buf.ReadFrom(r.Body)
    body := buf.Bytes()

    return sha1.Sum(append(url,body...))
}

func handler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        rec := httptest.NewRecorder()

        hash := hashKey(r)
        log.Println(hash)

        // Perform the real request with a recorder
        p.ServeHTTP(rec, r)

        // Copy the recorded data to the actual ResponseWriter
        for k, v := range rec.Header() {
            w.Header()[k] = v
        }
        w.Header().Set("X-Eidetic", "Live request")
        w.Write(rec.Body.Bytes())
        // log.Println(rec.Body)
    }
}