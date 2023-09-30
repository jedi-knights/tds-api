#!/bin/bash

# Determine the new version from Semantic Release (adjust as needed)
NEW_VERSION=$(semantic-release-cli print-version)

# Update the VERSION constant in main.go using sed
sed -i "s/const VERSION = \".*\"/const VERSION = \"$NEW_VERSION\"/" main.go

# Commit the updated main.go file
git add main.go
git add .semrel
git commit -m "chore: Update VERSION to $NEW_VERSION in main.go"
