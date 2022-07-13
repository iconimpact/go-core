# Auth

Package auth provides utility functions for authorization.

## HMAC authorization (server to server)

[HMAC](https://www.wolfe.id.au/2012/10/20/what-is-hmac-authentication-and-why-is-it-useful/)

The following functions are available:

- `HMACSign` and `HMACVerify` functions for creating and verifying hex-encoded
sha-512 HMAC signatures for a specified secret and a payload.

- `SetHMACHeaders` and `GetHMACHeaders` functions for setting and getting custom
HTTP requests headers to be used for authorization.

- `HTTPMiddleware` function which creates an HTTP middleware for authorizing
requests using the signatures and headers mentioned above.

:bulb: See [hmac_test.go](./hmac_test.go) for examples on how to use these.
