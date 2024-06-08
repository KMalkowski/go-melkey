package main

import (
	"context"
	"fmt"
	"lambda-func/app"
	"lambda-func/middleware"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const (
	invalidRequest = "Invalid Request"
)

type MyEvent struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

func HandleRequest(ctx context.Context, event MyEvent) (string, error) {
	if event.Name == "" {
		return invalidRequest, fmt.Errorf("name cannot be empty in request")
	}

	return fmt.Sprintf("Successful call - thank you %s", event.Name), nil
}

func ProtectedHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       "This is a secret path",
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	myApp := app.NewApp()
	lambda.Start(func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		switch request.Path {
		case "/register":
			return myApp.ApiHandler.RegisterUserHandler(request)
		case "/login":
			return myApp.ApiHandler.LoginUserHandler(request)
		case "/protected":
			return middleware.ValidateJWTMiddleware(ProtectedHandler)(request)
		default:
			return events.APIGatewayProxyResponse{
				Body:       "not found",
				StatusCode: http.StatusNotFound,
			}, nil
		}
	})
}
