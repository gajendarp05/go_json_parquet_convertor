package main

import (
	"context"
	"go_json_parquet_convertor/convertor"
	"go_json_parquet_convertor/infra"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func main() {
	// Intialize with environment variables
	poller := infra.CheckEnv("Poller")
	worker := infra.CheckEnv("Worker")
	awsSQSName := infra.CheckEnv("AWS_SQS")
	awsS3BucketName := infra.CheckEnv("AWS_S3")

	// These environment variables are credentials for accesing AWS infrastructure
	// Want to read more, checkout this: https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html
	infra.CheckEnv("AWS_ACCESS_KEY_ID")
	infra.CheckEnv("AWS_SECRET_ACCESS_KEY")
	infra.CheckEnv("AWS_DEFAULT_REGION")

	numPoller, err := infra.Str_Int(poller)
	if err != nil {
		log.Fatalf("%s", err)
	}
	numWorker, err := infra.Str_Int(worker)
	if err != nil {
		log.Fatalf("%s", err)
	}

	c := &convertor.Convertor{
		Poller:       numPoller,
		Worker:       numWorker,
		SQSName:      awsSQSName,
		S3BucketName: awsS3BucketName,
	}

	// Using AWS SDK v2 - github.com/aws/aws-sdk-go-v2/config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Panic("configuration error: ", err)
	}

	// Get URL of SQS queue
	c.SQS = sqs.NewFromConfig(cfg)
	urlResult, err := c.SQS.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{
		QueueName: &c.SQSName,
	})
	if err != nil {
		log.Panic("Got an error getting the SQS URL: ", err)
	}
	c.SQSUrl = *urlResult.QueueUrl

	c.Start()

}