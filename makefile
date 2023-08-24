install:
	@go install

build: 
	@go build -o auth

run: build
	@./auth