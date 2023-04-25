#!/bin/bash

set -e

# Set the default image tag to latest
DEMO_IMAGE_TAG="${DEMO_IMAGE_TAG:-latest}"

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
  echo "Docker is not installed. Please install Docker and try again."
  exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
  echo "Docker Compose is not installed. Please install it and try again."
  exit 1
fi

# Display a welcome message
echo "Welcome to the Akita Demo!"
echo "This script will help you run the demo and prompt you for any required environment variables."
echo "After that, it will start the demo for you."
echo ""

# Check for default credentials in $HOME/.akita/credentials.yaml
CONFIG_FILE="$HOME/.akita/credentials.yaml"
if [[ -f "$CONFIG_FILE" ]]; then
  echo "Found existing Akita credentials in $CONFIG_FILE"
  # Set the API credentials unless a user has provided them via environment variables
  if [[ -z "$AKITA_API_KEY_ID" ]]; then
  	AKITA_API_KEY_ID=$(awk -F ': ' '/api_key_id/{print $2}' "$CONFIG_FILE")
  fi
  if [[ -z "$AKITA_API_KEY_SECRET" ]]; then
  	AKITA_API_KEY_SECRET=$(awk -F ': ' '/api_key_secret/{print $2}' "$CONFIG_FILE")
  fi
fi

# Prompt the user for environment variables
while [[ -z "$AKITA_API_KEY_ID" ]]; do
  read -p "Enter your Akita API Key ID: " AKITA_API_KEY_ID
done

while [[ -z "$AKITA_API_KEY_SECRET" ]]; do
  read -p "Enter your Akita API Key Secret: " AKITA_API_KEY_SECRET
done

while [[ -z "$AKITA_PROJECT_NAME" ]]; do
  read -p "Enter your Akita Project Name: " AKITA_PROJECT_NAME
done

# Pull required images from Docker Hub if they don't exist locally
check_and_pull_image() {
  image="$1"

  docker inspect "${image}" > /dev/null 2>&1 || (echo "Image ${image} not found locally. Pulling..." && docker pull "${image}")
}

echo ""
echo "Pulling the latest Akita demo images..."
# Always pull the latest image of the CLI
docker pull akitasoftware/cli:latest
# Pull the demo images
check_and_pull_image "akitasoftware/demo-server:${DEMO_IMAGE_TAG}"
check_and_pull_image "akitasoftware/demo-client:${DEMO_IMAGE_TAG}"


# Run docker-compose
# Never pull images from Docker Hub so that the local image can be used if it exists
echo ""
echo "Starting the Akita demo..."
AKITA_API_KEY_ID="${AKITA_API_KEY_ID}" \
	AKITA_API_KEY_SECRET="${AKITA_API_KEY_SECRET}" \
	AKITA_PROJECT_NAME="${AKITA_PROJECT_NAME}" \
  DEMO_IMAGE_TAG="${DEMO_IMAGE_TAG}" \
	docker-compose up -d --always-recreate-deps --pull 'never' --remove-orphans

echo ""
echo "The Akita demo is now up and running!"
echo "View the agent logs with: 'docker compose logs akita'"
echo "To stop the demo run: 'make stop-demo'"
echo "Enjoy!"
