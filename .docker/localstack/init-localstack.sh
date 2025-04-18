#!/bin/bash

echo "Esperando o LocalStack iniciar..."
sleep 10

echo "Criando bucket S3..."
aws --endpoint-url=http://localhost:4566 s3 mb s3://codeflix-videos --region us-east-1

echo "Listando buckets S3..."
aws --endpoint-url=http://localhost:4566 s3 ls

echo "Configuração do LocalStack concluída!"