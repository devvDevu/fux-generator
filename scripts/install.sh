#!/bin/bash

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

for cmd in git go sudo; do
    if ! command -v $cmd &> /dev/null; then
        echo -e "${RED}Error: $cmd is not installed.${NC}"
        exit 1
    fi
done

TMP_DIR=$(mktemp -d)

echo "Cloning repository..."
git clone -q -b ver_1 https://github.com/devvDevu/ca-generator.git $TMP_DIR

echo "Building binary..."
cd $TMP_DIR
go install golang.org/x/tools/cmd/goimports@latest
go build -o ca-gen cmd/main.go

echo "Installing to /usr/local/bin..."
sudo mv -f ca-gen /usr/local/bin/

rm -rf $TMP_DIR
echo -e "${GREEN}Done! Use command: ca-gen${NC}"
