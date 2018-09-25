package main

import (
    "net/http"
    "regexp"
    "os"

    "github.com/gorilla/context"

    "project/server"
    log "project/logging"
)

var (
    git_VERSION string
    key_COOKIE_SIGNING string
    key_JWT_SIGNING string
)

func main() {

    if len(os.Args) != 4 {
        panic("Usage: "+os.Args[0]+" <port, e.g 80> <template_regex> <log_dir>")
    }
    port := os.Args[1]
    template_dir := os.Args[2]
    log_dir := os.Args[3]

    // For the fun of it
    pattern := "^([0-9]{1,4}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$"
    if matched, err := regexp.MatchString(pattern, port); err!=nil || !matched {
        panic("Invalid port number: "+port)
    }

    log.Init(log_dir, "myserver_log", 524288, true)
    server.Init(key_JWT_SIGNING, key_COOKIE_SIGNING, template_dir, git_VERSION)

    http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

    http.HandleFunc("/", server.Validate(server.HomeHandler))
    http.HandleFunc("/login", server.Validate(server.LoginHandler))

    log.Enterf("Starting test server at :%s", port)
    if err := http.ListenAndServe(":"+port, context.ClearHandler(http.DefaultServeMux)); err != nil {
        panic(err)
    }
}
