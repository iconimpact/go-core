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

### How to sign HTTP requests

:information_source: This section is especially useful for services written in languages other than Go.

To sign an HTTP request (so that it passes the HMAC authorization checks) the following 4 headers need to be set on it:

- `X-Auth-App-ID`

  - This is the ID of the application sending the request. Needs to be configured also on the receiving server, together with it's corresponding shared secret.
  - Example value: `Dispoman`

- `X-Auth-Nonce`

  - This is some random value (e.g. a UUID or a number) that must be **unique among all requests** that the server application receives within a certain duration (e.g. for Abfallpass API server this duration is **2 minutes**).

- `X-Auth-Timestamp`

  - The time at which the request is sent as number of seconds since UNIX epoch start time (i.e. since January 1st, 1970 at 00:00:00 UTC). It must not be older than a certain duration (the same duration that is used for checking the validity of the `X-Auth-Nonce` header mentioned above - e.g. for Abfallpass API server this duration is **2 minutes**).

- `X-Auth-Signature`

  - This is the signature itself.
  - It's value needs to be computed like this (pseudocode): ***`HEX( HMAC( SHA512, nonce+timestamp, shared-secret ) )`***.
    - Or, to put it in words, it must be the **hexadecimal encoding** of an **SHA 512 HMAC hash** of the **concatenated nonce and timestamp** (in this order - nonce immediately followed by the timestamp, without any other character between them) created using the **shared secret**.
