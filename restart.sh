#!/usr/bin/env bash

set -e

# Change into the script's directory.
cd "$(dirname "$0")"

# Stop the demo
./stop.sh
# Start the demo
./run.sh
