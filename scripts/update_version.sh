#!/bin/bash
# exit when any command fails
set -e

# Determine the new version from Semantic Release (adjust as needed)
NEW_VERSION=$(semantic-release-cli print-version)

# Update the VERSION constant in main.go using sed
echo $NEW_VERSION > VERSION

git status
cat VERSION

# Commit the updated main.go file
git add VERSION
git status
git commit -m "chore: Update VERSION to $NEW_VERSION in main.go [skip ci]"
git push