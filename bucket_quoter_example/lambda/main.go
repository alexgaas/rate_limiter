// main.go
package main

import (
	"basic_lambda/client"
	"basic_lambda/config"
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

func basicHandler() (client.Response, error) {
	conf, err := config.LoadConf()
	if err != nil {
		return client.Response{
			StatusCode: 500,
			Headers:    map[string]string{"content-type": "application/json"},
			Body:       err.Error(),
		}, nil
	}

	if conf.DryRun {
		return client.Response{
			StatusCode: 200,
			Headers:    map[string]string{"content-type": "application/json"},
			Body:       "dry run test",
		}, nil
	}

	c, err := client.NewClient(
		client.AsService("lambda"),
		client.WithApiSubId(conf.ApiSubKey),
		client.WithHTTPHost(conf.LimiterHost),
		client.WithInsecureSkipVerify(),
		//client.WithTLSCert(conf.CaCertPath),
	)
	if err != nil {
		return client.Response{
			StatusCode: 500,
			Headers:    map[string]string{"content-type": "application/json"},
			Body:       err.Error(),
		}, nil
	}

	resp, respErr := c.GetLimitResponse(context.Background())
	if respErr != nil {
		return client.Response{
			StatusCode: 500,
			Headers:    map[string]string{"content-type": "application/json"},
			Body:       respErr.Error(),
		}, nil
	}

	return *resp, nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(basicHandler)
}
