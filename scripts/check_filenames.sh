#!/bin/bash

# Find all Go files in the project directory and its subdirectories, except in the vendor directory
for file in $(find . -name "*.go" -not -path "./vendor/*"); do

  # Check if the file name contains uppercase letters
  if [[ "$file" =~ [A-Z] ]]; then
    echo "Error: $file contains uppercase letters. All Go files in the project (excluding vendor directory) must use snake_case"
    exit 1
  fi
done

echo "All Go files in the project (excluding vendor directory) use lowercase letters"
exit 0
