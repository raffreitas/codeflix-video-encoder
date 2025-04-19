#!/bin/bash

# Set AWS credentials for LocalStack
export AWS_ACCESS_KEY_ID=localstack
export AWS_SECRET_ACCESS_KEY=localstack
export AWS_DEFAULT_REGION=us-east-1

echo "Esperando o LocalStack iniciar..."
sleep 10

echo "Criando bucket S3..."
aws --endpoint-url=http://localstack:4566 s3 mb s3://codeflix-local

echo "Listando buckets S3..."
aws --endpoint-url=http://localstack:4566 s3 ls

echo "Configuração do LocalStack concluída!"