package api

import (
	"errors"
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
)

func (api *API) Search(context *gin.Context) {
	requestingUserID := context.MustGet("userID").(int64)
	queryString := context.Query("querystring")

	if queryString == "" {
		context.AbortWithError(http.StatusBadRequest, errors.New("Invalid querystring. Expected non-empty."))
		return
	}

	ftsResult, err := api.service.FullTextSearch(queryString, requestingUserID)
	if err != nil {
		statusCode := getStatusCodeFromError(err)
		context.AbortWithError(statusCode, err)
		return
	}

	context.IndentedJSON(http.StatusOK, ftsResult)
}
