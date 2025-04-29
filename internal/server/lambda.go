package server

import (
	"context"
	"go-race/internal/race"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type LambdaHandler struct {
	ctx        context.Context
	carTracker *race.CarTracker
}

func NewLambdaHandler() *LambdaHandler {
	ctx := context.Background()
	return &LambdaHandler{
		ctx:        ctx,
		carTracker: race.NewCarTracker(ctx),
	}
}

func (h *LambdaHandler) HandleRequest(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch request.RequestContext.EventType {
	case "CONNECT":
		return h.handleConnect(request)
	case "DISCONNECT":
		return h.handleDisconnect(request)
	case "MESSAGE":
		return h.handleMessage(request)
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Unsupported event type",
		}, nil
	}
}

func (h *LambdaHandler) handleConnect(request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("New connection established: %s", request.RequestContext.ConnectionID)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

func (h *LambdaHandler) handleDisconnect(request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Connection closed: %s", request.RequestContext.ConnectionID)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

func (h *LambdaHandler) handleMessage(request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Received message: %s", request.Body)

	// Start car race with the received message
	go h.carTracker.StartCarRace([]byte(request.Body))

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

func StartLambda() {
	handler := NewLambdaHandler()
	lambda.Start(handler.HandleRequest)
}
