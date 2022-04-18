# Respond

Package respond provides idiomatic way for API responses.
 - `respond.JSON` - for success responses.
 - `respond.JSONError` - for fail responses.
 - `respond.SetJSONErrorResponse` - useful for handling errors differently, define custom response.

`respond.JSONError` response depends on [go-core/errors](https://github.com/iconimpact/go-core/tree/master/errors) pkg for HTTP status and Msg message.

Feel free to add new functions or improve the existing code.

## Install

```bash
go get github.com/iconimpact/go-core/respond
```

## Usage and Examples

```go

// handle errors differently, define custom response.
errorRsp := func(err error) interface{} {
	var status int
	var errMsg string

	// set custom app err Message
	appErr, ok := err.(*errors.Error)
	if !ok {
		status = http.StatusInternalServerError
		errMsg = "Internal Server Error"
	} else {
		status = errors.ToHTTPStatus(appErr)
		errMsg = errors.ToHTTPResponse(appErr)
	}

	rsp := struct {
		Msg    string `json:"msg"`
		Status int    `json:"status"`
	}{
		Msg:    errMsg,
		Status: status,
	}

	return rsp
}
respond.SetJSONErrorResponse(errorRsp)

// usage in handler
func handleRoute(w http.ResponseWriter, r *http.Request) {

	data, err := loadFromDB()
	if err != nil {

	    // respond with an error
		respond.JSONError(w, logger, errors.E(err, errors.NotFound, "Data not found"))
		return // always return after responding

	}

	// respond with OK, and the data
	respond.JSON(w, logger, http.StatusOK, data)

}
```