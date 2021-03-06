package main

import (
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/brafales/expenses/monzo/handler"
)

func main() {
	categories, err := getCategories()
	if err != nil {
		panic(err)
	}
	h := handler.Handler{
		SnsTopicArn: os.Getenv("snsTopicArn"),
		Categories:  categories,
		SNSClient:   sns.New(session.New()),
	}
	lambda.Start(h.Handle)
}

func getCategories() ([]string, error) {
	rawData := os.Getenv("categoryData")
	var categoryMap map[string]string
	err := json.Unmarshal([]byte(rawData), &categoryMap)
	if err != nil {
		return []string{}, err
	}
	categories := make([]string, 0, len(categoryMap))
	for k := range categoryMap {
		categories = append(categories, k)
	}

	return categories, nil
}
