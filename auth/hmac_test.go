package auth_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/iconimpact/go-core/auth"
	"github.com/stretchr/testify/require"
)

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
	verifies, err := auth.HMACVerify(secret, payload, signature)
	require.NoError(t, err)
	require.True(t, verifies)
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
