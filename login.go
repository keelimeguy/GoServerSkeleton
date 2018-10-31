package main

import (
//    "encoding/json"
    "net/http"
		"fmt"

//    "github.com/nu7hatch/gouuid"

    log "project/logging"
)
// pass a copy of AuthContext os you can't mess it up
// this probably woulda been better with an structure taking a pointer to AuthContext
// at runtime and this as HTTPServe
func LoginHandler(aC AuthContext) (h http.HandlerFunc) {
	h = func(w http.ResponseWriter, r *http.Request) {
		var claims map[string]interface{}
		if r.Method == http.MethodPost {
			var err error
			if claims, err = aC.ReadJWT(w,r); err != nil {
				log.Errorf("ReadJWT failed in login with: %v", err)
			} else {
				fmt.Println(claims);
				aC.SetClaims(w, r, claims)
				w.WriteHeader(http.StatusOK)
			}

			/*if !ok {
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
			}*/
		} else {
				log.Warningf("Method not allowed.")
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
				return
		}
	}
	return h
}
/*
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
*/
