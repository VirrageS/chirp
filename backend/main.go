package main

/*
	TODO:
	  - generate a good secret key
	    (useful: http://security.stackexchange.com/questions/95972/what-are-requirements-for-hmac-secret-key,
	    	     https://elithrar.github.io/article/generating-secure-random-numbers-crypto-rand/)

	  - server should not start before index is created in elasticsearch,
	    see: https://github.com/VirrageS/chirp/issues/190
*/

import (
	"os"

	"github.com/Sirupsen/logrus"

	"github.com/VirrageS/chirp/backend/server"
)

func init() {
	logrus.SetOutput(os.Stderr)
}

func main() {
	s := server.New()
	s.Run(":8080")
}
