// main.go
package main

import (
	"basic_lambda/client"
	"basic_lambda/config"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

func basicHandler() (client.Response, error) {
	conf, err := config.LoadConf()
	if err != nil {
		return errResponse(err), nil
	}

	if conf.DryRun {
		return client.Response{
			StatusCode: 200,
			Headers:    client.Headers,
			Body:       "dry run test",
		}, nil
	}

	c, err := client.NewClient(
		client.AsService("lambda"),
		client.WithApiSubId(conf.ApiSubKey),
		client.WithHTTPHost(conf.LimiterHost),
		client.WithInsecureSkipVerify(),
	)
	if err != nil {
		return errResponse(err), nil
	}

	resp, err := c.GetLimitResponse(context.Background())
	if err != nil {
		return errResponse(err), nil
	}

	// log success
	log.Print(fmt.Sprintf("status code: %d, body: %s", resp.StatusCode, resp.Body))

	return resp, nil
}

func main() {
	log.SetOutput(os.Stdout)

	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(basicHandler)
}

func errResponse(err error) client.Response {
	return client.Response{
		StatusCode: 500,
		Headers:    client.Headers,
		Body:       err.Error(),
	}
}
