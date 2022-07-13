package auth_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/iconimpact/go-core/auth"
	"github.com/iconimpact/go-core/testhelpers"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestHMACMiddleware(t *testing.T) {
	logUnsugared, err := zap.NewDevelopment()
	require.NoError(t, err)
	log := logUnsugared.Sugar()

	log, logs := testhelpers.ObserveLogs(t, log)

	appIDFromConfig := "some-app-id"
	appSecretFromConfig := []byte("some-app-secret")
	secretsPerAppIDs := map[string][]byte{
		appIDFromConfig: appSecretFromConfig,
	}

	nonceExpiration := 2 * time.Second
	nonceCache := cache.New(nonceExpiration, nonceExpiration)

	type loggerContextKeyType struct{}
	loggerContextKey := loggerContextKeyType{}

	hmacMiddleware := auth.HMACMiddleware(
		secretsPerAppIDs,
		nonceCache,
		nonceExpiration,
		func(r *http.Request) *zap.Logger {
			return r.Context().Value(loggerContextKey).(*zap.SugaredLogger).Desugar()
		},
	)
	require.NotNil(t, hmacMiddleware)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(
		"GET", "http://some-service-to-service-request-url", nil)
	req = req.WithContext(context.WithValue(req.Context(), loggerContextKey, log))

	var appID, nonce, timestamp, signature string

	setRequestHeaders := func() {
		auth.SetHMACHeaders(req, appID, nonce, timestamp, signature)
	}

	requireNonceCache := func(nonceCacheSize int, nonceMustBeFoundInCache bool) {
		require.Equal(t, nonceCacheSize, nonceCache.ItemCount())
		_, nonceFoundInCache := nonceCache.Get(nonce)
		require.Equal(t, nonceMustBeFoundInCache, nonceFoundInCache)
	}

	requireOK := func(nonceCacheSize int, nonceMustBeFoundInCache bool) {
		respRecorder := httptest.NewRecorder()
		hmacMiddleware(handler).ServeHTTP(respRecorder, req)
		respBody := respRecorder.Body.String()
		require.Equal(t, http.StatusOK, respRecorder.Code, "body:\n%s", respBody)
		requireNonceCache(nonceCacheSize, nonceMustBeFoundInCache)
	}

	requireUnauthorized := func(
		nonceCacheSize int,
		nonceMustBeFoundInCache bool,
		expectedErrorFieldContaining string) {

		respRecorder := httptest.NewRecorder()
		hmacMiddleware(handler).ServeHTTP(respRecorder, req)
		respBody := respRecorder.Body.String()
		require.Equal(
			t, http.StatusUnauthorized, respRecorder.Code, "body:\n%s", respBody)
		require.Equal(t, `{"msg":"invalid authorization"}`, respBody)
		testhelpers.RequireLastLogEntry(t, logs, zapcore.ErrorLevel, "",
			map[string]string{
				"error": expectedErrorFieldContaining,
			})
		requireNonceCache(nonceCacheSize, nonceMustBeFoundInCache)
	}

	// happy path
	appID = appIDFromConfig
	nonce = uuid.NewString()
	timestamp = fmt.Sprintf("%d", time.Now().Unix())
	signature = auth.HMACSign(appSecretFromConfig, []byte(nonce+timestamp))
	setRequestHeaders()
	requireOK(1, true)

	// empty app ID header
	appID = ""
	nonce = uuid.NewString()
	timestamp = fmt.Sprintf("%d", time.Now().Unix())
	signature = auth.HMACSign(appSecretFromConfig, []byte(nonce+timestamp))
	setRequestHeaders()
	requireUnauthorized(1, false, "unauthorized: invalid authorization: "+
		"request header X-Auth-App-ID is missing or empty")

	// unknown app ID in header
	appID = "some-unknown-app-ID"
	nonce = uuid.NewString()
	timestamp = fmt.Sprintf("%d", time.Now().Unix())
	signature = auth.HMACSign(appSecretFromConfig, []byte(nonce+timestamp))
	setRequestHeaders()
	requireUnauthorized(1, false, "unauthorized: invalid authorization: request "+
		"header X-Auth-App-ID value 'some-unknown-app-ID' is an unknown app ID")

	// happy path again
	appID = appIDFromConfig
	nonce = uuid.NewString()
	timestamp = fmt.Sprintf("%d", time.Now().Unix())
	signature = auth.HMACSign(appSecretFromConfig, []byte(nonce+timestamp))
	setRequestHeaders()
	requireOK(2, true)

	// same nonce
	requireUnauthorized(
		2, true, "unauthorized: invalid authorization: nonce was already used")

	// happy path again: wait for expiration and reuse nonce
	time.Sleep(nonceExpiration + 500*time.Millisecond)
	timestamp = fmt.Sprintf("%d", time.Now().Unix())
	signature = auth.HMACSign(appSecretFromConfig, []byte(nonce+timestamp))
	setRequestHeaders()
	requireOK(1, true)

	// nonce OK, but timestamp is invalid (not a number)
	nonce = uuid.NewString()
	timestamp = ""
	signature = auth.HMACSign(appSecretFromConfig, []byte(nonce+timestamp))
	setRequestHeaders()
	requireUnauthorized(1, false, "invalid authorization timestamp")

	// nonce OK, but timestamp expired
	nonce = uuid.NewString()
	timestamp = fmt.Sprintf(
		"%d", time.Now().Add(-nonceExpiration-1*time.Second).Unix())
	signature = auth.HMACSign(appSecretFromConfig, []byte(nonce+timestamp))
	setRequestHeaders()
	requireUnauthorized(
		1, false, fmt.Sprintf("older than nonce expiration %s", nonceExpiration))

	// invalid signature: different secret
	nonce = uuid.NewString()
	timestamp = fmt.Sprintf("%d", time.Now().Unix())
	signature = auth.HMACSign(append(appSecretFromConfig, 'X'), []byte(nonce+timestamp))
	setRequestHeaders()
	requireUnauthorized(1, false, "invalid authorization signature")

	// invalid signature: different payload
	nonce = uuid.NewString()
	timestamp = fmt.Sprintf("%d", time.Now().Unix())
	signature = auth.HMACSign(append(appSecretFromConfig, 'X'), []byte(nonce+timestamp+"X"))
	setRequestHeaders()
	requireUnauthorized(1, false, "invalid authorization signature")

	// invalid signature
	nonce = uuid.NewString()
	timestamp = fmt.Sprintf("%d", time.Now().Unix())
	signature = auth.HMACSign(appSecretFromConfig, []byte(nonce+timestamp)) + "X"
	setRequestHeaders()
	requireUnauthorized(1, false, "invalid authorization signature")
}

func TestHMACSignAndVerify(t *testing.T) {
	secret := []byte("some-secret")
	nonce := "1"
	now, err := time.Parse("2006-01-02", "2022-07-13")
	require.NoError(t, err)
	timestamp := fmt.Sprintf("%d", now.Unix())
	payload := []byte(nonce + timestamp)

	// sign
	signature := auth.HMACSign(secret, payload)
	require.Equal(
		t,
		"1e18f0e0aca2ad42e180b3a70ef96cabce58aa05d6ab10ecc7e449a82029f7c7df"+
			"c8fc52f26c6c3cc669198c901036d4c1f72ee7047cdba3d176641ee81ea4f5",
		signature)

	// verify
	require.NoError(t, auth.HMACVerify(secret, payload, signature))
}

func TestSetAndGetHMACHeaders(t *testing.T) {
	appID := "some-app-id"
	nonce := "1"
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	signature := "some-signature"

	r, err := http.NewRequest("GET", "/some-url", nil)
	require.NoError(t, err)

	// set headers
	auth.SetHMACHeaders(r, appID, nonce, timestamp, signature)
	require.Equal(t, appID, r.Header.Get(auth.HMACHeaderAppID))
	require.Equal(t, nonce, r.Header.Get(auth.HMACHeaderNonce))
	require.Equal(t, timestamp, r.Header.Get(auth.HMACHeaderTimestamp))
	require.Equal(t, signature, r.Header.Get(auth.HMACHeaderSignature))

	// get headers
	gotAppID, gotNonce, gotTimestamp, gotSignature := auth.GetHMACHeaders(r)
	require.Equal(t, appID, gotAppID)
	require.Equal(t, nonce, gotNonce)
	require.Equal(t, timestamp, gotTimestamp)
	require.Equal(t, signature, gotSignature)
}
