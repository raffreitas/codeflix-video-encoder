#!/bin/bash

# Set AWS credentials for LocalStack
export AWS_ACCESS_KEY_ID=localstack
export AWS_SECRET_ACCESS_KEY=localstack
export AWS_DEFAULT_REGION=us-east-1

echo "Waiting for LocalStack to start..."
sleep 10

echo "Creating S3 bucket..."
aws --endpoint-url=http://localstack:4566 s3 mb s3://codeflix-local

echo "Listing S3 buckets..."
aws --endpoint-url=http://localstack:4566 s3 ls

echo "Completing LocalStack initialization..."