package auth

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"net/http"

	"github.com/iconimpact/go-core/errors"
)

// Request headers required for HMAC authorization.
const (
	HMACHeaderAppID     = "X-Auth-App-ID"
	HMACHeaderSignature = "X-Auth-Signature"
	HMACHeaderNonce     = "X-Auth-Nonce"
	HMACHeaderTimestamp = "X-Auth-Timestamp"
)

// HMACSign creates a new hex-encoded SHA512 HMAC signature for the specified
// secret and payload.
func HMACSign(secret, payload []byte) string {
	mac := hmac.New(sha512.New, secret)
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}

// HMACVerify verifies the given hex-encoded SHA512 HMAC signature for the
// specified secret and payload.
func HMACVerify(secret, payload []byte, signature string) (bool, error) {
	s, err := hex.DecodeString(signature)
	if err != nil {
		return false, errors.E(err)
	}
	mac := hmac.New(sha512.New, secret)
	mac.Write(payload)
	return hmac.Equal(s, mac.Sum(nil)), nil
}

// SetHMACHeaders sets the specified HMAC auth headers on an HTTP request.
func SetHMACHeaders(r *http.Request, appID, nonce, timestamp, signature string) {
	r.Header.Set(HMACHeaderAppID, appID)
	r.Header.Set(HMACHeaderNonce, nonce)
	r.Header.Set(HMACHeaderTimestamp, timestamp)
	r.Header.Set(HMACHeaderSignature, signature)
}

// GetHMACHeaders returns the HMAC auth headers from an HTTP request.
func GetHMACHeaders(r *http.Request) (appID, nonce, timestamp, signature string) {
	appID = r.Header.Get(HMACHeaderAppID)
	nonce = r.Header.Get(HMACHeaderNonce)
	timestamp = r.Header.Get(HMACHeaderTimestamp)
	signature = r.Header.Get(HMACHeaderSignature)
	return
}
