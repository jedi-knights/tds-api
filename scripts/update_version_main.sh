#!/bin/bash
# exit when any command fails
set -e

# Determine the new version from Semantic Release (adjust as needed)
NEW_VERSION=$(semantic-release-cli print-version)

# Update the version in main.go using sed
# sed -i "s/@version .*/@version $NEW_VERSION/" main.go
cat main.go | sed "s#\/\/ @version .*#\/\/ @version $NEW_VERSION#" > a.out
rm -f main.go
mv a.out main.go

go install github.com/swaggo/swag/cmd/swag@latest
swag init -g main.go

git status
cat main.go

# Commit the updated main.go file
git add main.go
git status
git commit -m "chore: Update version to $NEW_VERSION in main.go [skip ci]"
git push
