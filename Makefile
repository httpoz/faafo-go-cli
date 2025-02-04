ollama-up:
	echo "Setting up Ollama docker container"
	@docker compose up -d
	echo "Pulling Llama 3.2 model"
	@docker exec -it faafo-go-cli-ollama-1 ollama pull llama3.2
	echo "You are ready to go!"

app:
	echo "Running the app"
	@go run main.go
