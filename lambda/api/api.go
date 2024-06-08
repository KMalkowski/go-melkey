package api

import (
	"encoding/json"
	"fmt"
	"lambda-func/database"
	"lambda-func/types"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type ApiHandler struct {
	dbStore database.UserStore
}

func NewApiHandler(dbStore database.UserStore) ApiHandler {
	return ApiHandler{
		dbStore: dbStore,
	}
}

func (a ApiHandler) RegisterUserHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	registerUser := types.RegisterUser{}
	err := json.Unmarshal([]byte(request.Body), &registerUser)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "invalid request",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	if len(registerUser.Username) < 1 || len(registerUser.Password) < 1 {
		return events.APIGatewayProxyResponse{
			Body:       "invalid request",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	exists, err := a.dbStore.DoesUserExist(registerUser.Username)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "an error occured",
			StatusCode: http.StatusBadRequest,
		}, fmt.Errorf("error while checking if the user exists %w", err)
	}

	if exists {
		return events.APIGatewayProxyResponse{
			Body:       "user already exist",
			StatusCode: http.StatusBadRequest,
		}, fmt.Errorf("user already exists %w", err)
	}

	user, err := types.NewUser(registerUser)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "error registering user, try again later",
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("error while hashing the password for the user %w", err)
	}

	err = a.dbStore.InsertUser(user)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "error registering user, try again later",
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("error while iserting the user %w", err)
	}

	return events.APIGatewayProxyResponse{
		Body:       "user registered",
		StatusCode: http.StatusCreated,
	}, nil
}

func (a ApiHandler) LoginUserHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	loginRequest := LoginRequest{}
	err := json.Unmarshal([]byte(request.Body), &loginRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "invalid login request",
			StatusCode: http.StatusBadRequest,
		}, fmt.Errorf("could not unmarshal login request %w", err)
	}

	user, err := a.dbStore.GetUser(loginRequest.Username)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "internal server error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	valid := types.ValidatePassword(user.PasswordHash, loginRequest.Password)
	if !valid {
		return events.APIGatewayProxyResponse{
			Body:       "wrong username or password",
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	accessToken := types.CreateToken(user)
	successMsg := fmt.Sprintf(`{"access_token": "%s"}`, accessToken)

	return events.APIGatewayProxyResponse{
		Body:       successMsg,
		StatusCode: http.StatusOK,
	}, nil
}
