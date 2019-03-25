package handler_test

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"

	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/brafales/expenses/monzo/handler"
)

const monzoRequest = `{
	"type": "transaction.created",
	"data": {
	  "id": "txn_id",
	  "created": "2018-09-11T01:19:22.692Z",
	  "description": "OYSTER AUTO TOPUP      LONDON        GBR",
	  "amount": -2000,
	  "fees": {},
	  "currency": "GBP",
	  "merchant": {
		"id": "id",
		"group_id": "group_id",
		"created": "2018-07-04T01:03:28.487Z",
		"name": "Oyster Auto Topup",
		"logo": "https://mondo-logo-cache.appspot.com/twitter/tfl/?size=large",
		"emoji": "🚅",
		"category": "transport",
		"online": false,
		"atm": false,
		"address": {
		  "short_formatted": "Somewhere in 9th Fl Albany House, 0343 2221234 SW1H 0BD",
		  "formatted": "9th Fl Albany House, 0343 2221234 SW1H 0BD, United Kingdom",
		  "address": "9TH FL ALBANY HOUSE",
		  "city": "0343 2221234",
		  "region": "",
		  "country": "GBR",
		  "postcode": "SW1H0BD",
		  "latitude": 51.4994181,
		  "longitude": -0.1342755,
		  "zoom_level": 4,
		  "approximate": true
		},
		"updated": "2018-09-05T15:21:03.56Z",
		"metadata": {
		  "created_for_transaction": "tx",
		  "enriched_from_settlement": "tx",
		  "twitter_id": "tfl"
		},
		"disable_feedback": false
	  },
	  "notes": "",
	  "metadata": {
		"ledger_insertion_id": "a",
		"mastercard_auth_message_id": "s",
		"mastercard_lifecycle_id": "r"
	  },
	  "labels": null,
	  "account_balance": 0,
	  "attachments": null,
	  "international": null,
	  "category": "transport",
	  "is_load": false,
	  "settled": "2018-07-04T01:03:28.487Z",
	  "local_amount": -2000,
	  "local_currency": "GBP",
	  "updated": "2018-09-11T01:19:22.853Z",
	  "account_id": "tr",
	  "user_id": "ee",
	  "counterparty": {},
	  "scheme": "mastercard",
	  "dedupe_id": "dfjhdsjf",
	  "originator": false,
	  "include_in_spending": true,
	  "can_be_excluded_from_breakdown": true,
	  "can_be_made_subscription": true,
	  "can_split_the_bill": true
	}
	}`

const SNSMessage = `{"amount":-2000,"created":"2018-09-11T01:19:22.692Z","currency":"GBP","description":"OYSTER AUTO TOPUP      LONDON        GBR","category":"transport"}`

var categories = []string{"transport"}

type testSNSClient struct {
	snsiface.SNSAPI
	T      *testing.T
	Called bool
}

func (t *testSNSClient) Publish(input *sns.PublishInput) (*sns.PublishOutput, error) {
	if *input.Message != SNSMessage {
		t.T.Errorf("Unexpected message published to SNS. Expected %s, got %s", SNSMessage, *input.Message)
	}
	id := "id"
	output := sns.PublishOutput{
		MessageId: &id,
	}
	t.Called = true
	return &output, nil
}

func TestHappyPath(t *testing.T) {
	client := testSNSClient{
		T:      t,
		Called: false,
	}
	handler := handler.Handler{
		SnsTopicArn: "topic",
		Categories:  categories,
		SNSClient:   &client,
	}

	context := context.Background()
	proxyRequest := events.APIGatewayProxyRequest{
		Body: monzoRequest,
	}
	_, err := handler.Handle(context, proxyRequest)
	if err != nil {
		t.Error(err)
	}
	if !client.Called {
		t.Error("Expected SNS Client to be called")
	}
}
