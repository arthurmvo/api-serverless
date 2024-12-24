package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

var allowedOrigins = []string{
	"http://localhost:4200",
	"https://analyzer-pxg.vercel.app",
	"https://staging-analyzer.vercel.app",
	"https://analyzer-pxg-git-feature-lambda-arthurs-projects-b50471ed.vercel.app",
}

func cors(request events.LambdaFunctionURLRequest) (string, error) {
	origin := request.Headers["origin"]
	fmt.Println(origin, "origin")
	for _, allowedOrigin := range allowedOrigins {
		if origin == allowedOrigin {
			fmt.Println(allowedOrigin, "allowedOrigin")
			fmt.Println(origin == allowedOrigin, "origin == allowedOrigin")
			return origin, nil
		}
	}
	return "", fmt.Errorf("origin %s not allowed", origin)
}
