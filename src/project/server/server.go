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

type JWTClaims struct {
    UID string `json:"uid"`
    Exp int64 `json:"exp"`
}

func (this JWTClaims) Valid() error {
    return nil
}

func UpdateJWT(claims JWTClaims, w http.ResponseWriter, r *http.Request) (string, error) {
    session, _ := store.Get(r, "authorization")
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    token_str, err := token.SignedString([]byte(jwt_key))
    if err != nil {
        return "", err
    }
    session.Values["jwt"] = token_str
    session.Save(r, w)
    return token_str, nil
}

func Validate(next http.HandlerFunc) http.HandlerFunc {
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

        if session.Values["jwt"] != nil {
            token, err := jwt.ParseWithClaims(session.Values["jwt"].(string), &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
                return []byte(jwt_key), nil
            })
            var ok bool
            if claims, ok = token.Claims.(*JWTClaims); ok && token.Valid {
                if session.Values["jwt"] != nil {
                    if claims.Exp < now.Unix() {
                        log.Printf("Expired jwt. Forming new token.")
                        session.Values["jwt"] = nil
                    } else {
                        claims.Exp = now.Add(time.Hour * 24 * 7).Unix()
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
            if log.Check(err) {
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
