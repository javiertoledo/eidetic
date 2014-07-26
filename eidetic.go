package main

import(
    "os"
    "log"
    "net/url"
    "net/http"
    "net/http/httputil"
)

func main() {
    // The remote URL has to be provided as the first command line argument
    remote_host = os.Args[1]
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

func handler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("X-Eidetic", "Live request")
        log.Println(r.URL)
        p.ServeHTTP(w, r)
    }
}