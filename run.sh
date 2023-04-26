#!/usr/bin/env bash

set -e

# Change into the script's directory.
cd "$(dirname "$0")"

# Ensure that the stop and restart scripts are executable
chmod +x stop.sh
chmod +x restart.sh

# Set the default image tag to latest
DEMO_IMAGE_TAG="${DEMO_IMAGE_TAG:-latest}"

if [ "$DEV_MODE" = true ]; then
  echo "Running in dev mode. Building local images..."
  TAG="$DEMO_IMAGE_TAG" make build-images
fi

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
cat <<EOF
Welcome to the Akita Demo!
This script will help you run the demo and prompt you for any required environment variables.
After that, it will start the demo for you.

EOF

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

  if ! docker inspect "${image}" &> /dev/null; then
    echo "Image ${image} not found locally. Pulling..."
    docker pull "${image}"
  fi
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

cat <<EOF

The Akita demo is now up and running!
View the agent logs with: 'docker compose logs akita'
To stop the demo run: './stop.sh' from this directory
Enjoy!
EOF
