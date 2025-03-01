#!/bin/bash

GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE}Building Steganography Tool - Initial Setup...${NC}"

echo -e "${BLUE}Tidying modules...${NC}"
go mod tidy

echo -e "${BLUE}Building application...${NC}"
go build -v -o stego ./cmd/stego

if [ $? -eq 0 ]; then
    echo -e "${GREEN}Build successful! You can now run the application using: ./stego${NC}"
    chmod +x stego
else
    echo -e "${RED}Build failed!${NC}"
    exit 1
fi