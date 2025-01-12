package main

import (
	"context"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	origin, err := cors(request)

	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 403,
			Body:       err.Error(),
		}, nil
	}
	var response events.LambdaFunctionURLResponse

	switch {
	case (request.RequestContext.HTTP.Method == "OPTIONS"):
		return events.LambdaFunctionURLResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  origin,
				"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE",
				"Access-Control-Allow-Headers": "Content-Type",
			},
		}, nil

	case strings.HasPrefix(request.RequestContext.HTTP.Path, "/clans"):
		if request.RequestContext.HTTP.Method == "GET" {
			response, err = getClans()
		}
	// case strings.HasPrefix(request.RequestContext.HTTP.Path, "/drops"):
	// 	if request.RequestContext.HTTP.Method == "GET" {
	// 		response, err = getDrops()
	// 	} else if request.RequestContext.HTTP.Method == "POST" {
	// 		response, err = createDrop(request)
	// 	}
	// case strings.HasPrefix(request.RequestContext.HTTP.Path, "/hunts"):
	// 	if request.RequestContext.HTTP.Method == "GET" {
	// 		response, err = getHunts()
	// 	}
	// case strings.HasPrefix(request.RequestContext.HTTP.Path, "/pokemon"):
	// 	if request.RequestContext.HTTP.Method == "GET" {
	// 		response, err = getPokemons()
	// 	}
	// case strings.HasPrefix(request.RequestContext.HTTP.Path, "/player"):
	// 	if request.RequestContext.HTTP.Method == "PUT" {
	// 		response, err = updatePlayer(request)
	// 	} else if request.RequestContext.HTTP.Method == "GET" {
	// 		response, err = getPlayer(request)
	// 	} else if request.RequestContext.HTTP.Method == "DELETE" {
	// 		response, err = deletePlayer(request)
	// 	} else if request.RequestContext.HTTP.Method == "POST" {
	// 		response, err = createPlayer(request)
	// 	}

	default:
		response = events.LambdaFunctionURLResponse{
			StatusCode: 404,
			Body:       "Not Found",
		}
	}

	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, nil
	}

	// Set CORS headers
	if response.Headers == nil {
		response.Headers = make(map[string]string)
	}
	response.Headers["Access-Control-Allow-Origin"] = origin
	response.Headers["Access-Control-Allow-Headers"] = "Content-Type"
	response.Headers["Access-Control-Allow-Methods"] = "OPTIONS,POST,GET,PUT,DELETE"
	response.Headers["Content-Type"] = "application/json"

	return response, nil

}

func main() {
	lambda.Start(handler)
	// router := lambdahandler.NewRouter()

	// router.Origins = allowedOrigins
	// router.Methods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
	// router.Headers = []string{"Content-Type", "Authorization"}

	// router.Get("/clans", getClans)
	// router.Get("/drops", getDrops)
	// router.Post("/drops", createDrop)
	// router.Get("/hunts", getHunts)
	// router.Get("/pokemon", getPokemons)
	// router.Put("/player/:uid", updatePlayer)
	// router.Get("/player/:uid", getPlayer)
	// router.Delete("/player/:uid", deletePlayer)
	// router.Post("/player", createPlayer)

	// lambda.Start(router.HandleRequest)
}
