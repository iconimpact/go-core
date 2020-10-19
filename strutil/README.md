# Strutil

A collection of helper functions to deal with string encoding/decoding, generation.

Feel free to add new functions or improve the existing code.

## Install

```bash
go get github.com/iconmobile-dev/go-core/strutil
```

## Usage and Examples

`strutil.Encrypt` encrypts data using 256-bit AES-GCM.

`strutil.Decrypt` decrypts data using 256-bit AES-GCM.

`strutil.Random` generates a `NOT SECURE!` random string of defined length.

`strutil.RandomSecure` generates a SECURELY random string of defined length and type: `alpha`, `number`, `alpha-numeric`

`strutil.Hash` generates a hash of data using HMAC-SHA-512/256.
