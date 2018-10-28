package server

import (
    "encoding/json"
    "net/http"

    "github.com/gorilla/context"
    "github.com/nu7hatch/gouuid"

    log "project/logging"
)
// imagine a function that takes a password checker and returns a login handled
// this is also pretty much related to JWT tho.. I guess it has to be
// I suppose abstracting this so it takes a "password checker"
// naw I guess then what do you do right? at the end you need something to definitely make state;
func LoginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        var claims, ok = context.Get(r, "claims").(JWTClaims)
        if !ok {
            log.Errorf("Could not read claims.")
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        if claims.UID != "" {
            log.Errorf("Bad request.")
            http.Error(w, "Bad Request", http.StatusBadRequest)
            return
        }

        var data map[string]interface{}
        if err := json.NewDecoder(r.Body).Decode(&data); log.Check(err) {
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        if data["password"] == nil {
            log.Errorf("Bad request.")
            http.Error(w, "Bad Request", http.StatusBadRequest)
            return
        }

        if data["password"].(string) != entry_KEY {
            log.Warningf("Unauthorized.")
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        if uuid, err := uuid.NewV4(); log.Check(err) {
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        } else {
            claims.UID = uuid.String()
        }

        if _, err := UpdateJWT(claims, w, r); log.Check(err) {
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }
    } else {
        log.Warningf("Method not allowed.")
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }
}

func JoinHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        var claims, ok = context.Get(r, "claims").(JWTClaims)
        if !ok {
            log.Errorf("Could not read claims.")
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        if claims.UID != "" {
            log.Errorf("Bad request.")
            http.Error(w, "Bad Request", http.StatusBadRequest)
            return
        }

        var data map[string]interface{}
        if err := json.NewDecoder(r.Body).Decode(&data); log.Check(err) {
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }
				// okayyyyyyy just fuciing DIE
        http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
        return
    } else {
        log.Warningf("Method not allowed.")
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }
}
