package main

import (
	"net/http"
	"flag"
	"strconv"
  "html/template"
  "reflect"
	JustClaims "github.com/autopogo/justClaims"
	log "github.com/autopogo/justLogging"
)

var (
	git_VERSION string
)

var (
	// defining pointers. "flag name", "default", "description" -flag_name=whatever etc
	key_JWT_SIGNING = flag.String("Jwt_key",
		"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		"JWT Hash key");
	port = flag.Int("port",
		8080,
		"Specifies what port to run under")
	log_dir = flag.String("log_dir",
		"",
		"Specifies which folder to store logs")
		// empty string should use sys default
	log_prefix = flag.String("log_prefix",
		"",
		"Specifies what to prefix logname with")
		// empty string should use sys default
	certFile = flag.String("https certFile",
		"",
		"If you fill out a certFile and keyFile, we'll use https")
		// empty string should use sys default
	keyFile = flag.String("https keyFile",
		"",
		"If you fill out a certFile and keyFile, we'll use https")
		// empty string should use sys default 
	template_dir = flag.String("template_dir",
		"templates/",
		"Folder with templates if you're into that kind of thing, I'm not")
	_ = template_dir
)

// prepPublic() sets up the template datastructure by preparing custom template functions and finding the folder of templates

func main() {

	// We should also warn
	flag.Parse()

	// I want handle functions that can use or create contexts, with configuration
	// structures in main
	var myJWTContext = &JustClaims.Config {
		Jwt_key : *key_JWT_SIGNING,
		Cookie_name : "example_cookie",
		Cookie_persistent : true,
		Cookie_https : false,
		Cookie_server_only : false,
		MandatoryTokenRefresh : true,
		MandatoryTokenRefreshThreshold : 0.5,
		LifeSpanNano : 1e9 * 60,
	}

	log.Init(*log_dir, *log_prefix, 524288, true)
	tmpl := prepPublic(*template_dir)

	// A file server, if you want it. Serves ./public
	http.Handle("/public/",  http.FileServer(http.Dir(".")))

	// anything ending '/' is a tree and will match errant paths below it
	// or maybe it's just the rootmost '/'?
	// I demonstrate http.Handle and http.HandleFunc

	http.HandleFunc("/", NewLoginHandleFunc(myJWTContext));
	//http.HandleFunc("/", HomeHandler) // I didn't port this
	http.HandleFunc("/join", NewJoinHandler(myJWTContext))
	http.Handle("/account", NewAccountHandler(myJWTContext, git_VERSION, tmpl))

	log.Enterf("Starting test server [%v] at :%v", git_VERSION, *port)


	if ( (*certFile != "") && (*keyFile != "") ) {
		if err := http.ListenAndServeTLS(":"+strconv.Itoa(*port),*certFile, *keyFile, nil); err != nil {
			panic(err)
		}
	} else {
		if err := http.ListenAndServe(":"+strconv.Itoa(*port),nil); err != nil {
			panic(err)
		}
	}
}


// prepPublic() sets up the template datastructure by preparing custom template functions and finding the folder of templates
func prepPublic(template_dir string) (tmpl *template.Template){

	tmpl = template.Must(template.New("main").Funcs(template.FuncMap{
            "len": func(l interface{}) int {
                v := reflect.ValueOf(l)
                switch v.Kind() {
                    case reflect.Array, reflect.Slice, reflect.Map:
                        return v.Len()
                    default:
                        return 0
                }
            },
        }).ParseGlob(template_dir))
	return
}
