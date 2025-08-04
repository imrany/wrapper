proto:
	bash proto/generate_proto.sh

run:
	go run main.go --port=8080 --gemini-api-key=$(GEMINI_API_KEY)
