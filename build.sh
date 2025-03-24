#!/bin/bash

GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m'

echo -e "${BLUE}Building StegGo - Image Steganography Tool (CLI & API)${NC}"
echo -e "${BLUE}===============================================${NC}"
echo -e "${YELLOW}Current Date: 2025-03-24 03:08:08 (UTC)${NC}"
echo -e "${YELLOW}User: pranaykumar2${NC}"
echo

# Check command-line arguments
START_SERVER=false
if [[ "$1" == "--start-api" ]]; then
    START_SERVER=true
fi

# Create necessary directories
echo -e "${BLUE}Creating required directories...${NC}"
mkdir -p uploads
mkdir -p temp
mkdir -p test-images

# Update dependencies
echo -e "${BLUE}Tidying modules...${NC}"
go mod tidy

# Build CLI application
echo -e "${BLUE}Building CLI application...${NC}"
go build -v -o stego ./cmd/stego

if [ $? -eq 0 ]; then
    echo -e "${GREEN}CLI build successful!${NC}"
    chmod +x stego
else
    echo -e "${RED}CLI build failed!${NC}"
    exit 1
fi

# Build API server
echo -e "${BLUE}Building API server...${NC}"
go build -v -o steggo-api ./cmd/api

if [ $? -eq 0 ]; then
    echo -e "${GREEN}API server build successful!${NC}"
    chmod +x steggo-api
else
    echo -e "${RED}API server build failed!${NC}"
    exit 1
fi

echo
echo -e "${GREEN}=====================================================${NC}"
echo -e "${GREEN}Build completed successfully!${NC}"
echo -e "${GREEN}=====================================================${NC}"
echo
echo -e "${BLUE}You can now run:${NC}"
echo -e "  ${YELLOW}./stego${NC}               - For the CLI application"
echo -e "  ${YELLOW}./steggo-api${NC}          - For the API server"
echo
echo -e "${BLUE}API Documentation:${NC}"
echo -e "  Start the API server and visit: ${YELLOW}http://localhost:8080/swagger/index.html${NC}"
echo

# Check for test images
if [ ! -f "./test-images/sample.png" ]; then
    echo -e "${YELLOW}Note: No test images found in ./test-images directory.${NC}"
    echo -e "      Add test images for API testing."
fi

# Optionally start the API server
if [ "$START_SERVER" = true ]; then
    echo -e "${BLUE}Starting API server...${NC}"
    echo -e "${BLUE}Press Ctrl+C to stop the server${NC}"
    echo
    ./steggo-api
fi
