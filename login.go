// This is a package
package main

import (
		"net/http"
		"fmt"
		"github.com/nu7hatch/gouuid"
		"encoding/json"
		justClaims "github.com/autopogo/justClaims"
		log "github.com/autopogo/justLogging"
)

func NewLoginHandleFunc(jCC *justClaims.Config) (h http.HandlerFunc) {
	h = func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var claims map[string]interface{}
			var err error
			if claims, err = jCC.GetClaims(w,r); err != nil {
				log.Errorf("ReadJWT failed in login with: %v", err)
			} else {
				fmt.Println(claims);
				jCC.SetClaims(w, r, claims)
				w.WriteHeader(http.StatusOK)
			}

			if claims["uid"] != "" { // sure...
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

			//if data["password"].(string) != entry_KEY {
			// THERE IS NO REAL PASSWORD EH?
			if false {
					log.Warningf("Unauthorized.")
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
			}

			if uuid, err := uuid.NewV4(); log.Check(err) {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
			} else {
					claims["uid"] = uuid.String()
			}

/*			if _, err := jCC.SetClaims(w, r, claims); log.Check(err) {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
			}
			*/
			// It should return an error, but it doesn't
			jCC.SetClaims(w, r, claims);

		} else {
				log.Warningf("Method not allowed.")
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
				return
		}
	}
	return h
}

func NewJoinHandleFunc(jCC *justClaims.Config) (h http.HandlerFunc) {
	h = func(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
			claims, _ := jCC.GetClaims(w, r)
        if claims["uid"] != "" {
            log.Errorf("Bad request.")
            http.Error(w, "Bad Request", http.StatusBadRequest)
            return
        }

        var data map[string]interface{}
        if err := json.NewDecoder(r.Body).Decode(&data); log.Check(err) {
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }
        http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
        return
    } else {
        log.Warningf("Method not allowed.")
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }
	}
	return
}

