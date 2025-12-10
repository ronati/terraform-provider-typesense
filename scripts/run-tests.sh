#!/bin/bash
# Helper script to run acceptance tests locally with Typesense

set -e

CONTAINER_NAME="typesense-test"
TYPESENSE_IMAGE="typesense/typesense:29.0"
TYPESENSE_PORT=8108
API_KEY="test-api-key"

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== Terraform Typesense Provider Test Runner ===${NC}"
echo ""

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo -e "${RED}Error: Docker is not installed${NC}"
    echo "Please install Docker from https://docs.docker.com/get-docker/"
    exit 1
fi

# Function to check if Typesense is running
check_typesense() {
    curl -s -f -H "X-TYPESENSE-API-KEY: $API_KEY" "http://localhost:$TYPESENSE_PORT/health" > /dev/null 2>&1
}

# Function to cleanup Typesense container
cleanup() {
    if [ "$(docker ps -q -f name=$CONTAINER_NAME)" ]; then
        echo -e "${YELLOW}Stopping Typesense container...${NC}"
        docker stop $CONTAINER_NAME > /dev/null 2>&1
    fi
    if [ "$(docker ps -aq -f name=$CONTAINER_NAME)" ]; then
        docker rm $CONTAINER_NAME > /dev/null 2>&1
    fi
}

# Check if Typesense is already running
if check_typesense; then
    echo -e "${GREEN}✓ Typesense is already running at http://localhost:$TYPESENSE_PORT${NC}"
    CLEANUP_AFTER=false
else
    # Clean up any existing container
    cleanup

    echo -e "${YELLOW}Starting Typesense container...${NC}"
    docker run -d --name $CONTAINER_NAME \
        -p $TYPESENSE_PORT:8108 \
        -e TYPESENSE_DATA_DIR=/tmp \
        -e TYPESENSE_API_KEY=$API_KEY \
        $TYPESENSE_IMAGE > /dev/null

    # Wait for Typesense to be ready
    echo -e "${YELLOW}Waiting for Typesense to be ready...${NC}"
    MAX_WAIT=30
    WAITED=0
    while ! check_typesense; do
        if [ $WAITED -ge $MAX_WAIT ]; then
            echo -e "${RED}Error: Typesense did not start within $MAX_WAIT seconds${NC}"
            cleanup
            exit 1
        fi
        sleep 1
        WAITED=$((WAITED + 1))
    done

    echo -e "${GREEN}✓ Typesense is ready!${NC}"
    CLEANUP_AFTER=true
fi

echo ""
echo -e "${GREEN}Running acceptance tests...${NC}"
echo ""

# Set environment variables (will override defaults if needed)
export TYPESENSE_API_KEY=$API_KEY
export TYPESENSE_API_ADDRESS="http://localhost:$TYPESENSE_PORT"

# Run the tests
set +e
make testacc
TEST_EXIT_CODE=$?
set -e

echo ""

# Cleanup if we started the container
if [ "$CLEANUP_AFTER" = true ]; then
    echo -e "${YELLOW}Cleaning up Typesense container...${NC}"
    cleanup
    echo -e "${GREEN}✓ Cleanup complete${NC}"
fi

# Exit with the test exit code
if [ $TEST_EXIT_CODE -eq 0 ]; then
    echo ""
    echo -e "${GREEN}=== All tests passed! ===${NC}"
else
    echo ""
    echo -e "${RED}=== Tests failed ===${NC}"
fi

exit $TEST_EXIT_CODE
