#!/bin/bash

GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m'

echo -e "${BLUE}Building StegGo - Image Steganography Tool (CLI & API/Web Server)${NC}"
echo -e "${BLUE}===============================================${NC}"
echo -e "${YELLOW}Last Updated On: 2025-03-31 14:07:35 (UTC)${NC}"
echo -e "${YELLOW}Created By: pranaykumar2${NC}"
echo

START_SERVER=false
if [[ "$1" == "--start-server" ]]; then
    START_SERVER=true
fi

echo -e "${BLUE}Creating required directories...${NC}"
mkdir -p uploads
mkdir -p temp
mkdir -p test-images
mkdir -p web/static/img

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

echo -e "${BLUE}Building API/Web server...${NC}"
go build -v -o steggo-server ./cmd/api

if [ $? -eq 0 ]; then
    echo -e "${GREEN}API/Web server build successful!${NC}"
    chmod +x steggo-server
else
    echo -e "${RED}API/Web server build failed!${NC}"
    exit 1
fi

echo
echo -e "${GREEN}=====================================================${NC}"
echo -e "${GREEN}Build completed successfully!${NC}"
echo -e "${GREEN}=====================================================${NC}"
echo
echo -e "${BLUE}You can now run:${NC}"
echo -e "  ${YELLOW}./stego${NC}               - For the CLI application"
echo -e "  ${YELLOW}./steggo-server${NC}       - For the API/Web server"
echo
echo -e "${BLUE}Web Interface:${NC}"
echo -e "  Start the server and visit: ${YELLOW}http://localhost:8080${NC}"
echo
echo -e "${BLUE}API Documentation:${NC}"
echo -e "  Start the server and visit: ${YELLOW}http://localhost:8080/swagger/index.html${NC}"
echo

if [ ! -f "test-images/sample.jpg" ]; then
    echo -e "${YELLOW}Note: No test images found in ./test-images directory.${NC}"
    echo -e "      Add test images for testing."
fi

if [ "$START_SERVER" = true ]; then
    echo -e "${BLUE}Starting API/Web server...${NC}"
    echo -e "${BLUE}Press Ctrl+C to stop the server${NC}"
    echo
    ./steggo-server
fi
