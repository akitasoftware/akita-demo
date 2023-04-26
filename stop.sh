#!/usr/bin/env bash

set -e

# Change into the script's directory.
cd "$(dirname "$0")"

# Call docker compose down to stop the demo
# The -v flag removes all volumes
# The --rmi flag removes all images
DEMO_IMAGE_TAG="${DEMO_IMAGE_TAG:-latest}" docker compose down -v --rmi all
