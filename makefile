# Go parameters for backend
BACKEND_DIR=.
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_NAME=kolibra

# Frontend parameters
FRONTEND_DIR=frontend
NPMCMD=yarn
NPMINSTALL=$(NPMCMD) install
NPMBUILD=$(NPMCMD) build --outDir ../static/dist

# Default target executed when no arguments are given to make.
# It will run the 'build_all' target
default: build_all

# Build both backend and frontend
build_all: build_frontend build_backend

# Builds the backend
build_backend:
	cd $(BACKEND_DIR) && $(GOBUILD) -o $(BINARY_NAME) -v

# Cleans the backend
clean_backend:
	cd $(BACKEND_DIR) && $(GOCLEAN)
	rm -f $(BACKEND_DIR)/$(BINARY_NAME)

# Runs backend tests
test_backend:
	cd $(BACKEND_DIR) && $(GOTEST) -v ./...

# Builds the frontend
build_frontend:
	cd $(FRONTEND_DIR) && $(NPMINSTALL) && $(NPMBUILD)

# Cleans the frontend
clean_frontend:
	rm -rf $(FRONTEND_DIR)/node_modules
	rm -rf $(FRONTEND_DIR)/dist
	rm -rf ${BACKEND_DIR}/static/dist

# Clean both backend and frontend
clean_all: clean_backend clean_frontend

.PHONY: default build_all build_backend clean_backend test_backend build_frontend clean_frontend clean_all