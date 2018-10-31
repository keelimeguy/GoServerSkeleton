package main

import (
    "encoding/json"
    "net/http"

    "github.com/gorilla/context"

    log "project/logging"
)
// just don't care about go generated HTML yet.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
    // Handle unknown path with 404
    if r.URL.Path != "/" {
        log.Warningf("404. Not found.")
        w.WriteHeader(http.StatusNotFound)
        log.Check(tmpl.ExecuteTemplate(w, "404", nil))
        return
    }

    var claims, ok = context.Get(r, "claims").(JWTClaims)
    if !ok {
        log.Errorf("Could not read claims.")
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    // GET the main page
    if r.Method == http.MethodGet {
        log.Check(tmpl.ExecuteTemplate(w, "header", ImageNav(claims.UID, "http:/"+"/"+r.Host)))
        if claims.UID == "" {
            log.Check(tmpl.ExecuteTemplate(w, "join_us", nil))
        } else {
            log.Check(tmpl.ExecuteTemplate(w, "welcome", nil))
        }
        log.Check(tmpl.ExecuteTemplate(w, "footer", struct{Version string}{version}))
        return

    } else if r.Method == http.MethodPut {
        var data map[string]interface{}
        if err := json.NewDecoder(r.Body).Decode(&data); log.Check(err) {
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        if data["type"] == nil {
            log.Errorf("Bad request.")
            http.Error(w, "Bad Request", http.StatusBadRequest)
            return
        }

        log.Printf(data["type"].(string))

        switch data_type := data["type"].(string); data_type {
            default:
                log.Errorf("Bad request.")
                http.Error(w, "Bad Request", http.StatusBadRequest)
                return
        }

    } else {
        log.Warningf("Method not allowed.")
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }
}*/
