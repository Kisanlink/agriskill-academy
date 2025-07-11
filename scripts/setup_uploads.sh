#!/bin/bash

# Setup uploads directory structure for ASA Backend
# This script creates the necessary directories for file uploads

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}✓${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_info() {
    echo -e "${CYAN}ℹ${NC} $1"
}

echo "=== ASA Backend Uploads Directory Setup ==="

# Create base uploads directory
if [ ! -d "uploads" ]; then
    mkdir -p uploads
    print_status "Created uploads directory"
else
    print_info "Uploads directory already exists"
fi

# Create subdirectories
directories=("resumes" "certificates" "images" "documents")

for dir in "${directories[@]}"; do
    if [ ! -d "uploads/$dir" ]; then
        mkdir -p "uploads/$dir"
        print_status "Created uploads/$dir directory"
    else
        print_info "uploads/$dir directory already exists"
    fi
done

# Set proper permissions
chmod 755 uploads
chmod 755 uploads/*

print_status "Set proper permissions on upload directories"

echo ""
echo "=== Upload Directory Structure ==="
tree uploads 2>/dev/null || ls -la uploads/

echo ""
print_status "Uploads directory setup completed!"
echo ""
echo "Directory structure:"
echo "  uploads/"
echo "  ├── resumes/      (for resume files)"
echo "  ├── certificates/ (for certificate files)"
echo "  ├── images/       (for profile photos)"
echo "  └── documents/    (for other documents)"
echo ""
echo "Files will be automatically saved to these directories when uploaded." 