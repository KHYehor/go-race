package main

import (
	"go-race/internal/server"
	"os"
)

func main() {
	// Check if running in AWS Lambda environment
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		server.StartLambda()
	} else {
		// Local development mode
		srv := server.New()
		srv.Start()
	}
}
