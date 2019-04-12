package main

import (
    "net/http"
		"html/template"
		JustClaims "github.com/autopogo/justClaims"
    log "github.com/autopogo/justLogging"
)
// more templates NOPE

type AccountHandlerConfig struct {
	jc *JustClaims.Config
	version string
	tmpl *template.Template
}
func NewAccountHandler(_jc *JustClaims.Config, _version string, _tmpl *template.Template) (*AccountHandlerConfig){
	// I think it's okay to return a pointer
	return &AccountHandlerConfig{jc : _jc, version: _version, tmpl: _tmpl}
}

func (h *AccountHandlerConfig) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
			claims, _ := h.jc.GetClaims(w, r)
        if claims["uid"] == "" {
            log.Errorf("Unauthorized.")
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        log.Check(h.tmpl.ExecuteTemplate(w, "header", ImageNav(claims["uid"].(string), "http:/"+"/"+r.Host)))
        log.Check(h.tmpl.ExecuteTemplate(w, "account", nil))
        log.Check(h.tmpl.ExecuteTemplate(w, "footer", struct{Version string}{h.version}))
        return

    } else {
        log.Warningf("Method not allowed.")
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

}
