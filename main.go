package main

import (
    "net/http"
    "regexp"
    "time"
    "os"

    "github.com/gorilla/context"

    "project/server"
    log "project/logging"
)

var ( // should these be checked?
    git_VERSION string
    key_COOKIE_SIGNING string
    key_JWT_SIGNING string
)
// TODO: server transitions + loadbalancing. Instances should talk to themselves and replace each other gracefully
// TODO: https
func main() {

    if len(os.Args) != 4 {
        panic("Usage: "+os.Args[0]+" <port, e.g 80> <template_regex> <log_dir>")
    } // yah no
    port := os.Args[1]
    template_dir := os.Args[2]
    log_dir := os.Args[3]

    // For the fun of it
    pattern := "^([0-9]{1,4}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$"
    if matched, err := regexp.MatchString(pattern, port); err!=nil || !matched {
        panic("Invalid port number: "+port)
    }

		// probably look into the log manging systems on the platform
    log.Init(log_dir, "myserver_log", 524288, true)
    server.Init(key_JWT_SIGNING, key_COOKIE_SIGNING, template_dir, git_VERSION)

    // TODO: remove NoCache for static files (exists now since static files constantly changing)
    http.Handle("/public/", NoCache(http.StripPrefix("/public/", http.FileServer(http.Dir("./public")))))

    http.HandleFunc("/", server.Validate(server.HomeHandler)) // server was probably a bad module choice name since http already implements Server. it seems like we're not with golang's notation either, which makes it hard to read.
    http.HandleFunc("/login", server.Validate(server.LoginHandler))
    http.HandleFunc("/join", server.Validate(server.JoinHandler))
    http.HandleFunc("/account", server.Validate(server.AccountHandler))

    log.Enterf("Starting test server [%s] at :%s", git_VERSION, port)
    if err := http.ListenAndServe(":"+port, context.ClearHandler(http.DefaultServeMux)); err != nil {
        panic(err)
    }
}

/******************************************************/
// For debugging, stop cacheing of static files:
// https://stackoverflow.com/questions/33880343/go-webserver-dont-cache-files-using-timestamp

var epoch = time.Unix(0, 0).Format(time.RFC1123)

var no_cache_headers = map[string]string{
    "Expires":         epoch,
    "Cache-Control":   "no-cache, private, max-age=0",
    "Pragma":          "no-cache",
    "X-Accel-Expires": "0",
}

var etag_headers = []string{
    "ETag",
    "If-Modified-Since",
    "If-Match",
    "If-None-Match",
    "If-Range",
    "If-Unmodified-Since",
}

func NoCache(h http.Handler) http.Handler {
    fn := func(w http.ResponseWriter, r *http.Request) {
        // Delete any ETag headers that may have been set
        for _, v := range etag_headers {
            if r.Header.Get(v) != "" {
                r.Header.Del(v)
            }
        }

        // Set our NoCache headers
        for k, v := range no_cache_headers {
            w.Header().Set(k, v)
        }

        h.ServeHTTP(w, r)
    }

    return http.HandlerFunc(fn) // it looks like a function call, but its a type, and i hate it
}
/******************************************************/
