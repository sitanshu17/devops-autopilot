#!/bin/bash
# Format all Go files
echo "Formatting Go files..."
gofmt -w .

echo "Checking formatting..."
unformatted=$(gofmt -l .)
if [ -n "$unformatted" ]; then
    echo "The following files are not properly formatted:"
    echo "$unformatted"
    exit 1
else
    echo "All files are properly formatted!"
fi