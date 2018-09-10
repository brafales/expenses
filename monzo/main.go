package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/brafales/expenses/monzo/handler"
)

func main() {
	h := handler.Handler{
		SnsTopicArn: os.Getenv("snsTopicArn"),
	}
	lambda.Start(h.Handle)
}
