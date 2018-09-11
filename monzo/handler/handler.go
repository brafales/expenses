package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

var categories = [...]string{"eating_out", "transport"}

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// {
//     "type": "transaction.created",
//     "data": {
//         "account_id": "acc_00008gju41AHyfLUzBUk8A",
//         "amount": -350,
//         "created": "2015-09-04T14:28:40Z",
//         "currency": "GBP",
//         "description": "Ozone Coffee Roasters",
//         "id": "tx_00008zjky19HyFLAzlUk7t",
//         "category": "eating_out",
//         "is_load": false,
//         "settled": true,
//         "merchant": {
//             "address": {
//                 "address": "98 Southgate Road",
//                 "city": "London",
//                 "country": "GB",
//                 "latitude": 51.54151,
//                 "longitude": -0.08482400000002599,
//                 "postcode": "N1 3JD",
//                 "region": "Greater London"
//             },
//             "created": "2015-08-22T12:20:18Z",
//             "group_id": "grp_00008zIcpbBOaAr7TTP3sv",
//             "id": "merch_00008zIcpbAKe8shBxXUtl",
//             "logo": "https://pbs.twimg.com/profile_images/527043602623389696/68_SgUWJ.jpeg",
//             "emoji": "🍞",
//             "name": "The De Beauvoir Deli Co.",
//             "category": "eating_out"
//         }
//     }
// }
type monzoEvent struct {
	Type string `json:"type"`
	Data struct {
		AccountID   string    `json:"account_id"`
		Amount      int       `json:"amount"`
		Created     time.Time `json:"created"`
		Currency    string    `json:"currency"`
		Description string    `json:"description"`
		ID          string    `json:"id"`
		Category    string    `json:"category"`
		IsLoad      bool      `json:"is_load"`
		Merchant    struct {
			Address struct {
				Address   string  `json:"address"`
				City      string  `json:"city"`
				Country   string  `json:"country"`
				Latitude  float64 `json:"latitude"`
				Longitude float64 `json:"longitude"`
				Postcode  string  `json:"postcode"`
				Region    string  `json:"region"`
			} `json:"address"`
			Created  time.Time `json:"created"`
			GroupID  string    `json:"group_id"`
			ID       string    `json:"id"`
			Logo     string    `json:"logo"`
			Emoji    string    `json:"emoji"`
			Name     string    `json:"name"`
			Category string    `json:"category"`
		} `json:"merchant"`
	} `json:"data"`
}

type expenseEvent struct {
	Amount      int       `json:"amount"`
	Created     time.Time `json:"created"`
	Currency    string    `json:"currency"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
}

// Handler is a type that will handle a Monzo request coming through an API Gateway Proxy Request
type Handler struct {
	SnsTopicArn string
}

// Handle handles a Monzo request coming through an API Gateway Proxy Request
func (h *Handler) Handle(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	event := &monzoEvent{}

	if err := json.Unmarshal([]byte(request.Body), event); err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	if h.interested(*event) {
		err := h.publishEvent(*event)
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}
	}

	return events.APIGatewayProxyResponse{Body: "OK", StatusCode: 200}, nil
}

func (h *Handler) publishEvent(event monzoEvent) error {
	expense := expenseEvent{
		Amount:      event.Data.Amount,
		Created:     event.Data.Created,
		Description: event.Data.Description,
		Category:    event.Data.Category,
		Currency:    event.Data.Currency,
	}
	expenseBytes, err := json.Marshal(expense)
	if err != nil {
		return err
	}

	svc := sns.New(session.New())
	params := &sns.PublishInput{
		Message:  aws.String(string(expenseBytes)),
		TopicArn: aws.String(h.SnsTopicArn),
	}
	resp, err := svc.Publish(params)

	if err != nil {
		return err
	}
	fmt.Println(resp)
	return nil
}

func (h *Handler) interested(event monzoEvent) bool {
	for _, v := range categories {
		if v == event.Data.Category {
			return true
		}
	}
	return false
}
