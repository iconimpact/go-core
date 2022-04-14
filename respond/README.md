# Respond

Package respond provides idiomatic way for API responses.
 - `respond.JSON` - for success responses.
 - `respond.JSONError` - for fail responses.

`respond.JSONError` response depends on [go-core/errors](https://github.com/iconimpact/go-core/tree/master/errors) pkg for HTTP status and Msg message.

Feel free to add new functions or improve the existing code.

## Install

```bash
go get github.com/iconmobile-dev/go-core/respond
```

## Usage and Examples

```go

// @Summary [get] handleRoute
// @Description Swagger doc for GET handleRoute
// @Description
// @Tags session-endpoints
// @Accept  json
// @Produce json
// @Param Authorization header string true "Example: Bearer token"
// @Param orderid path string true "Order ID"
// @Success 200 {object} respond.JSONMsg{data=<data struct type sent to respond.JSON} "Success"
// @Failure 400 {object} respond.JSONMsg "Invalid request JSON"
// @Failure 403 {object} respond.JSONMsg "Forbidden"
// @Failure 422 {object} respond.JSONMsg "Params validation error"
// @Failure 404 {object} respond.JSONMsg "Order not found"
// @Failure 500 {object} respond.JSONMsg "Internal server error"
// @Router /v1/handleRoute [get]
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