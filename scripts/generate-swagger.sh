#!/bin/sh

# Generate Swagger docs using swag
swag init -g cmd/main.go -o docs
