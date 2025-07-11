#!/bin/bash

INTERNAL_MOCKS_DIR=internal/mocks
INTERNAL_FILES=(
  repositories/application.go
  repositories/offer.go
  banks/bank.go
  controllers/ws/ws.go
)

generate () {
  echo "Generating mock for ${1}"
  mockgen -source="${1}" -destination="${2}" || exit 1
}

command -v mockgen > /dev/null || go install github.com/golang/mock/mockgen@v1.6.0
rm -rf "${INTERNAL_MOCKS_DIR}"

for f in ${INTERNAL_FILES[@]}; do
  generate "internal/${f}" "${INTERNAL_MOCKS_DIR}/${f}"
done
