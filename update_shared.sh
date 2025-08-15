#!/bin/bash

SHARED_PACKAGE="github.com/leetcode-repetition/shared"
VERSION="${1:-latest}"
LAMBDAS_DIR="$(pwd)/lambdas"

GREEN='\033[1;32m'
CYAN='\033[1;36m'
NC='\033[0m'

echo -e "Updating ${GREEN}$SHARED_PACKAGE${NC} to ${CYAN}$VERSION${NC} in all lambdas"

ACTUAL_VERSION=""
VERSION_FOUND=false

for lambda_dir in "$LAMBDAS_DIR"/*; do
  if [ -d "$lambda_dir" ]; then
    lambda_name=$(basename "$lambda_dir")
    echo "Updating $lambda_name..."

    cd "$lambda_dir"

    go get "$SHARED_PACKAGE@$VERSION"
    go mod tidy
    
    if [ "$VERSION_FOUND" = false ]; then
      ACTUAL_VERSION=$(grep "$SHARED_PACKAGE" go.mod | awk '{print $2}')
      if [ -n "$ACTUAL_VERSION" ]; then
        VERSION_FOUND=true
      fi
    fi
    
    cd - > /dev/null
  fi
done

echo -e "Successfully updated all lambdas to ${CYAN}$ACTUAL_VERSION${NC} of ${GREEN}$SHARED_PACKAGE${NC}"

echo -e "Updating git submodule to ${CYAN}$ACTUAL_VERSION${NC}..."
# First update to get the latest refs
git submodule update --remote shared

# Then checkout the specific version we found in go.mod
cd shared
git checkout $ACTUAL_VERSION
cd - > /dev/null

echo -e "${GREEN}Git submodule updated successfully${NC}"