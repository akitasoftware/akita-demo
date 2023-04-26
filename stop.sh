#!/bin/bash

set -e

# Set the default image tag to latest
DEMO_IMAGE_TAG="${DEMO_IMAGE_TAG:-latest}"

DEMO_IMAGE_TAG="DEMO_IMAGE_TAG" docker compose down -v --rmi all
