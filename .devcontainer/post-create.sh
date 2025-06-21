#!/bin/bash
set -e

GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Check if required environment variables are set
if [[ -z "${GOLANGCI_LINT_VERSION}" ]]; then
    echo -e "${RED}Error: GOLANGCI_LINT_VERSION environment variable is not set.${NC}"
    exit 1
fi

if [[ -z "${JUST_INSTALL_PATH}" ]]; then
    echo -e "${RED}Error: JUST_INSTALL_PATH environment variable is not set.${NC}"
    exit 1
fi

# Install golangci-lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | \
    sh -s -- -b "$(go env GOPATH)/bin" "${GOLANGCI_LINT_VERSION}"

# Install just
mkdir -p "${JUST_INSTALL_PATH}"
curl --proto '=https' --tlsv1.2 -sSf https://just.systems/install.sh | \
    bash -s -- --to "${JUST_INSTALL_PATH}"

echo -e "${GREEN}postCreateCommand script completed successfully.${NC}"