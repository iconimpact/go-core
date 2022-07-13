package auth

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/iconimpact/go-core/errors"
	"github.com/iconimpact/go-core/respond"
	"go.uber.org/zap"
)

// Request headers required for HMAC authorization.
const (
	HMACHeaderAppID     = "X-Auth-App-ID"
	HMACHeaderSignature = "X-Auth-Signature"
	HMACHeaderNonce     = "X-Auth-Nonce"
	HMACHeaderTimestamp = "X-Auth-Timestamp"
)

// HMACNonceCache is an interface abstracting away the cache implementation
// for caching nonces used for HMAC authorization.
type HMACNonceCache interface {
	Get(string) (interface{}, bool)
	Set(string, interface{}, time.Duration)
}

// HMACMiddleware validates the signature header which is a HEX-encoded SHA512
// HMAC of nonce, timestamp and secret.
// Signature timestamp is considered valid for nonceExpiration duration and
// nonce values must be unique within this timeframe.
func HMACMiddleware(
	secretsPerAppIDs map[string][]byte,
	nonceCache HMACNonceCache,
	nonceExpiration time.Duration,
	requestLogger func(r *http.Request) *zap.Logger,
) func(next http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := requestLogger(r)

			appID, nonce, timestamp, signature := GetHMACHeaders(r)

			if len(appID) == 0 {
				err := fmt.Errorf(
					"invalid authorization: request header %s is missing or empty",
					HMACHeaderAppID)
				err = errors.E(err, errors.Unauthorized, "invalid authorization")
				respond.JSONError(w, log, err)
				return
			}

			sharedSecret, ok := secretsPerAppIDs[appID]
			if !ok {
				err := fmt.Errorf(
					"invalid authorization: request header %s value '%s' is an unknown app ID",
					HMACHeaderAppID, appID)
				err = errors.E(err, errors.Unauthorized, "invalid authorization")
				respond.JSONError(w, log, err)
				return
			}

			_, found := nonceCache.Get(nonce)
			if found {
				err := fmt.Errorf("invalid authorization: nonce was already used")
				err = errors.E(err, errors.Unauthorized, "invalid authorization")
				respond.JSONError(w, log, err)
				return
			}

			ts, err := strconv.ParseInt(timestamp, 10, 64)
			if err != nil {
				err = fmt.Errorf("invalid authorization timestamp: %w", err)
				err = errors.E(err, errors.Unauthorized, "invalid authorization")
				respond.JSONError(w, log, err)
				return
			}
			t := time.Unix(ts, 0)

			age := time.Since(t)
			if age > nonceExpiration {
				err = fmt.Errorf(
					"invalid authorization: timestamp '%s' (unix second %d) has age %s "+
						"older than nonce expiration %s",
					t, ts, age, nonceExpiration)
				err = errors.E(err, errors.Unauthorized, "invalid authorization")
				respond.JSONError(w, log, err)
				return
			}

			err = HMACVerify(sharedSecret, []byte(nonce+timestamp), signature)
			if err != nil {
				err = fmt.Errorf("invalid authorization signature: %v", err)
				err = errors.E(err, errors.Unauthorized, "invalid authorization")
				respond.JSONError(w, log, err)
				return
			}

			nonceCache.Set(nonce, struct{}{}, nonceExpiration)

			next.ServeHTTP(w, r)
		})
	}
}

// HMACSign creates a new hex-encoded SHA512 HMAC signature for the specified
// secret and payload.
func HMACSign(secret, payload []byte) string {
	mac := hmac.New(sha512.New, secret)
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}

// HMACVerify verifies the given hex-encoded SHA512 HMAC signature for the
// specified secret and payload.
func HMACVerify(secret, payload []byte, signature string) error {
	s, err := hex.DecodeString(signature)
	if err != nil {
		return errors.E(err)
	}
	mac := hmac.New(sha512.New, secret)
	mac.Write(payload)
	if !hmac.Equal(s, mac.Sum(nil)) {
		return errors.E(fmt.Errorf("signature mismatch"))
	}
	return nil
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
