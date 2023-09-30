#!/bin/bash

# Determine the new version from Semantic Release (adjust as needed)
NEW_VERSION=$(semantic-release-cli print-version)

# Update the version in main.go using sed
sed -i "s/@version .*/@version $NEW_VERSION/" main.go

go install github.com/swaggo/swag/cmd/swag@latest
swag init -g main.go


# Commit the updated main.go file
git add main.go
git add docs
git commit -m "chore: Update version to $NEW_VERSION in main.go"
