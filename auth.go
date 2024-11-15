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
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
)

type AuthFunc func(r *http.Request) (bool, error)

func New(cfg Credentials) AuthFunc {
	return cfg.Auth
}

func extractCredentials(authHeader string) (username, password string, err error) {
	// Check if the header starts with "Basic "
	if !strings.HasPrefix(authHeader, "Basic ") {
		return "", "", fmt.Errorf("invalid authorization header")
	}

	// Extract the Base64 part of the header
	encoded := authHeader[len("Basic "):]

	// Decode the Base64 string
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", "", fmt.Errorf("failed to decode Base64: %v", err)
	}

	// Split the decoded string into username and password
	parts := strings.SplitN(string(decoded), ":", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid format for credentials")
	}

	// Hash the username and password separately using SHA-256
	username = sha256ToHex(parts[0])
	password = sha256ToHex(parts[1])

	return username, password, nil
}

func sha256ToHex(input string) string {
	hash := sha256.New()
	hash.Write([]byte(input))
	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes)
}

func (d *Credentials) Auth(r *http.Request) (bool, error) {
	aHeader := r.Header.Get("Authorization")

	username, password, err := extractCredentials(aHeader)
	if err != nil {
		return false, err
	}

	if d.User != username || d.Pass != password {
		return false, fmt.Errorf("authorization failed")
	}

	return true, nil
}
