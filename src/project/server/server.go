package server

import (
    "html/template"
    "net/http"
    "reflect"
    "time"

    "github.com/gorilla/context"
    "github.com/gorilla/sessions"
    "gopkg.in/dgrijalva/jwt-go.v3"

    log "project/logging"
)
// so after reading this, I don't know what server is. Is it a reimplementation of httpserver? Actually, so it's something that takes a handler and binds it tightly with hard-coded JWT operations and expectations? i feel like we're not fully utilizing the libraries we've adopted, nor are we fully shimming them. so it's like, we're dependent but hateful of them.
var (
    version string
    jwt_key string
    store sessions.Store
    tmpl *template.Template
    entry_KEY string
)

func Init(_jwt_key string, _cookie_key string, template_dir string, _version string) {
    store = sessions.NewCookieStore([]byte(_cookie_key))
    jwt_key = _jwt_key
    version = _version
    entry_KEY = "password"

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
}

type JWTClaims struct { // how JWT is implemented probably very specific so maybe uh, maybe this should be an interface? with defaults?
    UID string `json:"uid"`
    Exp int64 `json:"exp"` // default claims is best, just a map yes
}

func (this JWTClaims) Valid() error {
    return nil // TODO
}

func UpdateJWT(claims JWTClaims, w http.ResponseWriter, r *http.Request) (string, error) {
    session, _ := store.Get(r, "authorization") // grabbing something from the header obviously?
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    token_str, err := token.SignedString([]byte(jwt_key))
    if err != nil {
        return "", err
    }
    session.Values["jwt"] = token_str
    session.Save(r, w) // okay so you're writing the session structure to the r/w... not writing the r
    return token_str, nil
}

// what a way to use interfaces 
func Validate(next http.HandlerFunc) http.HandlerFunc { // I'm pretty sure this is what valid was supposed to do
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if r := recover(); r != nil {
                log.Errorf("***PANIC!*** %v", r)
            }
        }()

        log.Enterf("%v %s", r.Method, r.URL.EscapedPath())

        session, _ := store.Get(r, "authorization")

        now := time.Now()
        var claims *JWTClaims

        if session.Values["jwt"] != nil { // still don't know if this server/client side
            token, err := jwt.ParseWithClaims(session.Values["jwt"].(string), &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
                return []byte(jwt_key), nil
            })
            var ok bool
						// okay so, token.Valid does nothing...,, yet we call it, and then do all the things it was supposed to do inside of an if block, but it was supposed to do it during "ParseWithClaims"
            if claims, ok = token.Claims.(*JWTClaims); ok && token.Valid { 
                if session.Values["jwt"] != nil {
                    if claims.Exp < now.Unix() {
                        log.Printf("Expired jwt. Forming new token.")
                        session.Values["jwt"] = nil
                    } else {
                        claims.Exp = now.Add(time.Hour * 24 * 7).Unix() // wow that's a lot
                        if _, err := UpdateJWT(*claims, w, r); log.Check(err) { 
                            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                            return
                        }
                    }
                }
            } else {
                log.Check(err)
                log.Errorf("Malformed JWT.")
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }
        }

        if session.Values["jwt"] == nil {
            claims = &JWTClaims{UID: "", Exp: now.Add(time.Hour * 24 * 7).Unix()}
            tokenString, err := UpdateJWT(*claims, w, r)
            if log.Check(err) { // oh now I get it and I kinda like it
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
                return
            }
            log.Printf("New token: %s", tokenString)
        }

        log.Printf("Claims: %v", *claims)
        context.Set(r, "claims", *claims)
        next(w, r)
    })
}
