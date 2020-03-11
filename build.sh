#!/usr/bin/env bash

IMAGE_NAME=alirezafaghani/surge
IMAGE_TAG=latest

echo "About to build the $IMAGE_NAME:$IMAGE_TAG image"
docker build -t "$IMAGE_NAME:$IMAGE_TAG" .

echo "Signing into registry!"
docker login -u "$REGISTRY_USER" -p "$REGISTRY_PASSWORD"

echo "Pushing the $IMAGE_NAME:$IMAGE_TAG image"
docker push "$IMAGE_NAME:$IMAGE_TAG"

echo "Removing the $IMAGE_NAME:$IMAGE_TAG from local build"
docker rmi -f "$IMAGE_NAME:$IMAGE_TAG"