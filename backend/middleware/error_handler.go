package middleware

import (
	"gopkg.in/gin-gonic/gin.v1"
)

func ErrorHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Next()

		var errorMessages []string

		for _, error := range context.Errors {
			errorMessages = append(errorMessages, error.Error())
		}

		if len(errorMessages) > 0 {
			context.JSON(-1, gin.H{
				"errors": errorMessages,
			})
		}
	}
}
