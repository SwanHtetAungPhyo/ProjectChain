.PHONY: full-node-one validator full-node client

full-node-one:
	@echo "Starting Full Node One..."
	@env PORT=8001 HTTP_PORT=3002 PEERS=localhost:8001 VALIDATOR_ADDRESS=9008 go run ./SwanLed/main.go

validator:
	@echo "Starting Validator..."
	@env PORT=8009 VALIDATOR_PORT=9001 PEERS=localhost:8001 go run ./validator/main.go

full-node:
	@echo "Starting Full Node..."
	@env PORT=8002 HTTP_PORT=3003 PEERS=localhost:8001 VALIDATOR_ADDRESS=9002 go run ./SwanLed/main.go

client:
	@echo "Starting Client..."
	@env TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZXhwIjoxNzQyNjYyMTkyLCJuYW1lIjoiSm9obiBEb2UifQ.qCxeF35YvWxgiHbbiQklzwUZ2aNFEqkuLBEWPCVSTmg go run ./client/main.go
