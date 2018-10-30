

// Author's notes:
// Gorilla's session half-implements RFC7519 (JWT), 
// ... signing and encryption, not validation.
// So no point in doing signing twice. We'll use setcookie + JWT
// RFC7519 wants us to use the Authorization header and... IDGAF, doesn't seem great.
// Only good for APIs


// In main or handler page: set authContext.
// ***a context(s) can implement ServeHTTP directly, and each Handler defines its credentials

// ***a type that takes a context(s) (via pointer/structure/array) 
// and implements your ServeHTTP function and therefor reuse context.

// ***a funciton that takes a context and returns a function with/ ServerHTTP signature.
// (same as above, different style)

// I mean, either we use what they gave us and do our best

// some jtw notes- claims is a map[string].(JSON) (i don't want to use json)
// It's value is whatever we're using. Needs to make sense tho.

// We have exp - expiresat - Int64
//				 iat - issuedAt - Int64 (Don't use it before it's valid :-p)
//				 nbf - verifynot - Int64 (Don't verify it?)
// being the most common.
// You might give someone a "session id" <-- an associate that with a security level
// Or track a state machine inside their cookie instead of on server
// Or give explicit access to a certain resource, through permission or lookup key
// I don't care what claims you add.
// This package was essentially conceived and architecuted by Keelin "kbw@autopogo.com" and then implemented by Andrew Pikul "ajp@autopogo.com" in a stunning reversal.

package server

import (
    "net/http"
    "time"

    "gopkg.in/dgrijalva/jwt-go.v3"
    log "project/logging"
)


type authContext struct (
    jwt_key string
		cookie_name string
		cookie_persistent bool // close on window (UX hint) 
		cookie_https bool // cookie only over https
		cookie_server_only bool // cookie not accessible to client (prevents XSS in modernt browsers)
		mandatoryTokenRefresh bool // do we refresh the token/cookie if it's below a certain time
		mandatoryTokenRefreshThereshold float32  // whats that time
		lifeTime int64 // seconds jwt+cookie have alive (if cookie persistent)
		// adapter for HTTPServe
		// adapter for HandlerFunc Factory (not a member)
)

// read claims. no cookie, return. claims invalid, return. if claims valid, and exp threshold, set cookie just in case.
func (aC *authContext) ReadJWT(w ResponseWriter, r *http.Request) (claims jwt.Claims, err error) {
	if cookie, err := r.Request.Cookie(aC.cookie_name); err != nil {
		if (err == ErrNoCookie) {
			return nil, ErrNoCookie
		}
		log.Errorf("Weird error trying to find cookie");
		return
	}
	t, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		return [](aC.jwt_key), nil
	})
	if err = t.Valid; err != nil {
		claims = t.Claims
		if val, ok := claims["exp"]; ok && aC.mandatoryTokenRefresh {
			//TODO if exp < mandatoryTokenRefreshThreshold time
			aC.setClaims(w, r, claims)
		} else {
			// 100 reasons why it might not exist, and therefore not our problem
		}
	} else {
		claims = nil
	}
}

// if no passed cookie, try to pull it, or else create it, ultimately update it, encode it, and set it, using the claims you pass. It will update claims. 
func (aC *authContext) setClaims(w ResponseWriter, r *http.Request, claims jwt.Claims){
	if (cookie == nil) {
		if (cookie, err := r.Cookie(aC.cookie_name)); err != nil {
			if (err != ErrNoCookie) {
					log.Errorf("Error having to do with cookie retrieval"); // log flooding?
				else {
					aC.setCookieDefaults(cookie, &claims);
				}
			}
		}
	} else {
		aC.updateExpiries(cookie, claims); //i'd want it to be a pointer but it's a reference type
	}
	token := jwt.NewWithClaims(jwt.SigningMethodsHS256, claims)
	ss, err := token.SignedString(aC.jwt_key); err != nil {
		log.Errorf("Couldn't unsign the encrypted string");
	}
	cookie.Value = ss;
	SetCookie(w, cookie)
}

// create a cookie and update its stuff
func (aC *authContext) setCookieDefaults(cookie *http.Cookie, claims jwt.Claims) {
	if (cookie == nil) {
		cookie = new(http.Cookie);
	}
	aC.updateExpiries(cookie, claims)
	cookie.HttpOnly = cookie_server_only
	cookie.Secure = cookie_https
}

// add the correct expiry to the cookie + jwt
func (aC *authContext) updateExpiries(cookie *http.Cookie, claims jwt.Claims) {
	if ( (aC.cookie_persistent)) ) {
		cookie.Expires = time.Now().Add(lifeTime)
	} else {
		cookie.MaxAge = time.Time(0)
	}
	claims["exp"] = time.Now().Unix()
}

// delete the cookie
func (aC *authContext) DeleteCookie(w ResponseWriter) (err error) {
	http.SetCookie(w, &http.Cookie{Name: aC.cookie_name, MaxAge: -1}
	// I left out things not set as optional but...
	// TODO return error
}

// if you're going to refresh the cookie, place the claims you got here
func (ac *authContext) newClaims(seconds int) (claims Claims, err error) {
	//create basic claims, i guess. not important.
}





func Validate(next http.HandlerFunc) http.HandlerFunc { // I'm pretty sure this is what valid was supposed to do
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if r := recover(); r != nil {
                log.Errorf("***PANIC!*** %v", r)
            }
        }()

        log.Enterf("%v %s", r.Method, r.URL.EscapedPath())
        log.Printf("Claims: %v", *claims)
    })
}

// TODO: writeFail
// if it's a bad token, you should respond with the proper header()... they may not care, so don't 401 everybody, but still. do they need to be authorized? can we treat it like first time? total rejection? let them know its a bad cookie tho
//	send WWW-Authenticate/401 and w/e info you want for the user
// const (
//    ValidationErrorMalformed        uint32 = 1 << iota // Token is malformed
//    ValidationErrorUnverifiable                        // Token could not be verified because of signing problems
//    ValidationErrorSignatureInvalid                    // Signature validation failed

    // Standard Claim validation errors
//    ValidationErrorAudience      // AUD validation failed
//    ValidationErrorExpired       // EXP validation failed
//    ValidationErrorIssuedAt      // IAT validation failed
//    ValidationErrorIssuer        // ISS validation failed
//    ValidationErrorNotValidYet   // NBF validation failed
//    ValidationErrorId            // JTI validation failed
//    ValidationErrorClaimsInvalid // Generic claims validation error
