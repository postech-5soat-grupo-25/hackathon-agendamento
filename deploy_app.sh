#!/bin/bash

# Variáveis
AWS_REGION="us-east-1"
AWS_ACCOUNT_ID="891376965147"
CLUSTE_NAME="hackathon-ecs-cluster"
SERVICE_NAME="agendamento-service"
ECR_REPOSITORY_URI="${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${SERVICE_NAME}"

# Fazer login no ECR
aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin $ECR_REPOSITORY_URI

# Build da imagem Docker
docker build -t $SERVICE_NAME .

# Taguear a imagem para o ECR
docker tag $SERVICE_NAME:latest $ECR_REPOSITORY_URI:latest

# Push da imagem para o ECR
docker push $ECR_REPOSITORY_URI:latest

# Atualizar o serviço ECS com a nova task definition
aws ecs update-service --cluster $CLUSTE_NAME --service $SERVICE_NAME --force-new-deployment
