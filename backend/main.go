package main

/*
	TODO:
	  - generate a good secret key
	    (useful: http://security.stackexchange.com/questions/95972/what-are-requirements-for-hmac-secret-key,
	    	     https://elithrar.github.io/article/generating-secure-random-numbers-crypto-rand/)
*/

import "github.com/VirrageS/chirp/backend/server"

func main() {
	server := server.CreateNew()
	server.Run(":8080")
}
