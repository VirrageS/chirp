package middleware

import (
	"errors"
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
)

// TODO: (if needed) support more than only JSON content type
func ContentTypeChecker() gin.HandlerFunc {
	return func(context *gin.Context) {
		contentType := context.Request.Header.Get("Content-Type")
		if contentType != "application/json" {
			context.AbortWithError(
				http.StatusUnsupportedMediaType,
				errors.New("Required content-type: application/json."),
			)
			return
		}

		context.Next()
	}
}
