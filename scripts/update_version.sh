#!/bin/bash

# Determine the new version from Semantic Release (adjust as needed)
NEW_VERSION=$(semantic-release-cli print-version)

# Update the VERSION constant in main.go using sed
sed -i "s#.*#$NEW_VERSION#" > a.out
rm -f VERSION
mv a.out VERSION

# Commit the updated main.go file
git add .
git commit -m "chore: Update VERSION to $NEW_VERSION in main.go"
