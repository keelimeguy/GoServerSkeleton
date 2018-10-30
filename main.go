package main

import (
    "net/http"
    "regexp"
    "time"
    "os"
		"flag"
		"strconv"
		"context"
    "project/server"
    log "project/logging"
)

var ( // should these be checked?
    git_VERSION string
)
func main() {

		// defining pointers. "flag name", "default", "description" -flag_name=whatever etc
		var key_JWT_SIGNING = flag.String("jwt_key",
			"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			"JWT Hash key");
		var port = flag.Int("port",
			8080,
			"Specifies what port to run under")
		var log_dir = flag.Str("log_dir",
			"",
			"Specifies which folder to store logs")
			// empty string should use sys default
		var log_prefix = flag.Str("log_prefix",
			"",
			"Specifies what to prefix logname with")
			// empty string should use sys default
		var template_dir = flagStr("template_dir",
			"templates/",
			"Folder with templates if you're into that kind of thing, I'm not")

		// We could do a lot more?
		// We should also warn
		flag.Parse()

		log.Init(*log_dir, *log_prefix, 524288, true) // is this a proper log? we're using many things that take a proper log
    server.Init(*key_JWT_SIGNING, *template_dir, git_VERSION, "autopogo.com")

    // A file server, if you want it. Serves ./public
		// http.Handle("/public/",  http.FileServer(http.Dir(".")))

		// anything ending '/' is a tree and will match errant paths below it
		// or maybe it's just the rootmost '/'?

		// I want these functions to get access to my functions and my state
    http.HandleFunc("/", server.HomeHandler)
    http.HandleFunc("/login", server.LoginHandler)
    http.HandleFunc("/join", server.JoinHandler)
    http.HandleFunc("/account", server.AccountHandler)
		// Request.Method to make sure it's get/post/whatever

    log.Enterf("Starting test server [%s] at :%s", git_VERSION, *port)
		// ugh that Gorilla hack >:X
    if err := http.ListenAndServe(":"+strconv.Itoa(*port),http.DefaultServeMux);
			err != nil {
        panic(err)
    }
}

// any browser can turn of cache easily, don't turn it on for sever. Open dev tools (usually shift+ctl+c) and look for it.
