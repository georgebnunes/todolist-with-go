# Makefile para build e deploy do Lambda em Go
# Uso: make build | make deploy | make zip

BINARY_NAME  = bootstrap
LAMBDA_ZIP   = lambda.zip
BUILD_DIR    = ./bin
LAMBDA_ENTRY = ./cmd/lambda

# GOOS=linux GOARCH=amd64 compila para o ambiente do Lambda (Amazon Linux 2)
# CGO_ENABLED=0 gera binário estático — sem dependências de sistema
.PHONY: build
build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
		go build -tags lambda.norpc -ldflags="-s -w" \
		-o $(BUILD_DIR)/$(BINARY_NAME) $(LAMBDA_ENTRY)
	@echo "✅ Build OK: $(BUILD_DIR)/$(BINARY_NAME)"

# Empacota o binário no zip que o Lambda espera
.PHONY: zip
zip: build
	cd $(BUILD_DIR) && zip ../$(LAMBDA_ZIP) $(BINARY_NAME)
	@echo "✅ Zip OK: $(LAMBDA_ZIP)"

# Faz upload direto para o Lambda (para dev — em produção use CI/CD)
.PHONY: deploy
deploy: zip
	aws lambda update-function-code \
		--function-name todo-lambda \
		--zip-file fileb://$(LAMBDA_ZIP) \
		--region us-east-1
	@echo "✅ Deploy OK"

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR) $(LAMBDA_ZIP)

.PHONY: test
test:
	go test ./... -v -count=1

.PHONY: tidy
tidy:
	go mod tidy
