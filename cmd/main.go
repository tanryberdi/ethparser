package main

import (
	"log"
	"net/http"
	"os"

	"ethparser/internal/api"
	"ethparser/internal/parser"
)

func main() {
	logger := log.New(os.Stdout, "ethparser: ", log.LstdFlags|log.Lshortfile)

	logger.Printf("Starting Ethereum parser service ...")

	// Initialize the parser
	ethParser := parser.NewEthParser("https://ethereum-rpc.publicnode.com", logger)
	logger.Printf("Initialized parser with endpoint: https://ethereum-rpc.publicnode.com")

	if err := ethParser.Start(); err != nil {
		logger.Fatalf("failed to start parser: %v", err)
	}
	logger.Printf("Parser started successfully")

	// Start the API server
	server := api.NewServer(ethParser)
	server.RegisterRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Printf("starting server on :%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logger.Fatalf("failed to start server: %v", err)
	}
}
