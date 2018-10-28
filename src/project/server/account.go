package server

import (
    "net/http"

    "github.com/gorilla/context"

    log "project/logging"
)
// more templates NOPE
func AccountHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        var claims, ok = context.Get(r, "claims").(JWTClaims)
        if !ok {
            log.Errorf("Could not read claims.")
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        if claims.UID == "" {
            log.Errorf("Unauthorized.")
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        log.Check(tmpl.ExecuteTemplate(w, "header", ImageNav(claims.UID, "http:/"+"/"+r.Host)))
        log.Check(tmpl.ExecuteTemplate(w, "account", nil))
        log.Check(tmpl.ExecuteTemplate(w, "footer", struct{Version string}{version}))
        return

    } else {
        log.Warningf("Method not allowed.")
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

}
