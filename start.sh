#!/bin/bash

DOCKER_COMPOSE_URL="https://raw.githubusercontent.com/akitasoftware/aki-demo/main/docker-compose.yml"

# Display a welcome message
echo "Welcome to the Akita Demo!"
echo "This script will help you download the docker-compose file and set up the required environment variables."
echo "After that, it will start the demo for you."
echo ""

# Download docker-compose.yml
curl -sSL -o docker-compose.yml $DOCKER_COMPOSE_URL

# Prompt the user for environment variables
read -p "Enter your Akita API Key ID: " AKITA_API_KEY_ID
read -p "Enter your Akita API Key Secret: " AKITA_API_KEY_SECRET
read -p "Enter your Akita Project Name: " AKITA_PROJECT_NAME

# Run docker-compose
echo ""
echo "Starting the Akita demo services... Please wait."
DEMO_IMAGE_TAG=latest \
AKITA_API_KEY_ID=$AKITA_API_KEY_ID \
AKITA_API_KEY_SECRET=$AKITA_API_KEY_SECRET \
AKITA_PROJECT_NAME=$AKITA_PROJECT_NAME \
docker-compose up -d

echo ""
echo "The Akita demo services are now up and running!"
echo "View the agent logs with: "
echo "Enjoy!"

# Cleanup: Remove the docker-compose.yml file
rm docker-compose.yml
