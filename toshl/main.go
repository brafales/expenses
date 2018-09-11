package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/brafales/expenses/toshl/client"
)

// Handler will handle the lambda function
func Handler(ctx context.Context, snsEvent events.SNSEvent) {
	categories, err := categories()
	if err != nil {
		panic(err)
	}
	toshlClient := client.Client{
		AuthToken:   os.Getenv("token"),
		AccountID:   os.Getenv("accountId"),
		HTTPClient:  http.Client{},
		CategoryMap: categories,
	}

	for _, record := range snsEvent.Records {
		snsRecord := record.SNS

		expense, err := createExpense(snsRecord.Message)
		if err != nil {
			panic(err)
		}
		err = toshlClient.CreateEntry(&expense)
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	lambda.Start(Handler)
}

func categories() (map[string]string, error) {
	rawData := os.Getenv("categoryData")
	var categoryMap map[string]string
	err := json.Unmarshal([]byte(rawData), &categoryMap)
	if err != nil {
		return map[string]string{}, err
	}

	return categoryMap, nil
}

func createExpense(message string) (client.Expense, error) {
	var expense client.Expense
	err := json.Unmarshal([]byte(message), &expense)
	return expense, err
}
