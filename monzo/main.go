package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
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
	}
	lambda.Start(h.Handle)
}

func getCategories() ([]string, error) {
	rawData := os.Getenv("categoryData")
	fmt.Printf("categoryData: %s \n", rawData)
	var categoryMap map[string]string
	err := json.Unmarshal([]byte(rawData), categoryMap)
	if err != nil {
		return []string{}, err
	}
	categories := make([]string, 0, len(categoryMap))
	for k := range categoryMap {
		categories = append(categories, k)
	}

	return categories, nil
}
