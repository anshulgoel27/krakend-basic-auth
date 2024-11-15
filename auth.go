package auth

// import (
// 	"crypto/subtle"
// 	"encoding/base64"
// )

// // Validator defines the interface for all the possible validation processes
// type Validator interface {
// 	IsValid(subject string) bool
// }

// // NewCredentialsValidator creates a validator for a given credentials pair
// func NewCredentialsValidator(credentials Credentials) Validator {
// 	base := credentials.User + ":" + credentials.Pass
// 	header := "Basic " + base64.StdEncoding.EncodeToString([]byte(base))
// 	return authHeader{int32(len(header)), []byte(header)}
// }

// type authHeader struct {
// 	lenght  int32
// 	content []byte
// }

// // IsValid implements the Validator interface
// func (a authHeader) IsValid(subject string) bool {
// 	if subtle.ConstantTimeEq(int32(len(subject)), a.lenght) == 1 {
// 		return subtle.ConstantTimeCompare([]byte(subject), a.content) == 1
// 	}
// 	// Securely compare actual to itself to keep constant time, but always return false.
// 	return subtle.ConstantTimeCompare(a.content, a.content) == 1 && false
// }

import (
	"encoding/base64"
	"net/http"
	"strings"
)

type AuthFunc func(r *http.Request) bool

func New(cfg Credentials) AuthFunc {
	return cfg.Auth
}

func (d *Credentials) Auth(r *http.Request) bool {
	aHeader := r.Header.Get("Authorization")

	if aHeader == "" || !strings.HasPrefix(aHeader, "Basic ") {
		return false
	}

	if aHeader != "Basic "+base64.StdEncoding.EncodeToString([]byte(d.User+":"+d.Pass)) {
		return false
	}

	return true
}
