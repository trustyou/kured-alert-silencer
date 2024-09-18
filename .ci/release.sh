#!/bin/bash
set -e

print_help() {
    echo "Usage: $0 <MAJOR>.<MINOR>.<PATCH>"
    echo "Example: $0 1.2.3"
    exit 1
}

if [ -z "$1" ]; then
    echo "Error: No version provided."
    print_help
fi

if [[ "$1" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "Version is valid: $1"
else
    echo "Error: Version format is incorrect."
    print_help
fi

VERSION=$1 make update-changelog

git add .
git commit -m "release: Version v$1"

echo "After merging the PR, tag and release are automatically done"
