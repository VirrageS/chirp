package middleware

import "github.com/gin-gonic/gin"

// ErrorHandler handles api context errors. If there is any error in
// request chain, error handler will catch that and return proper JSON message.
func ErrorHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Next()

		var errorMessages []string

		for _, err := range context.Errors {
			errorMessages = append(errorMessages, err.Error())
		}

		if len(errorMessages) > 0 {
			context.JSON(-1, gin.H{
				"errors": errorMessages,
			})
		}
	}
}
