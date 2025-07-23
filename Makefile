BUILD_DIR := build
TARGET := $(BUILD_DIR)/boids

.PHONY: run all clean cache

all: build run

run:
	./$(TARGET)

build: main.go
	@mkdir -p $(BUILD_DIR)
	go build -o $(TARGET) main.go
	@echo "Successfully build target at $(TARGET)"

clean:
	go clean
	rm -rf $(BUILD_DIR)
