package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"tubes2/backend/api"
)

func main() {
	// Command line flags
	portPtr := flag.String("port", "8081", "Port for the server to listen on")
	modePtr := flag.String("mode", "server", "Mode to run (server or test)")
	targetPtr := flag.String("target", "Brick", "Target element to search for in test mode")
	algoPtr := flag.String("algo", "bfs", "Algorithm to use in test mode (bfs or dfs)")
	flag.Parse()

	// Ensure data directory exists
	workDir, _ := os.Getwd()
	recipesPath := filepath.Join(workDir, "data", "recipes.json")
	if _, err := os.Stat(recipesPath); os.IsNotExist(err) {
		log.Fatalf("Recipes file not found: %s\nMake sure the 'data' directory with 'recipes.json' exists", recipesPath)
	}

	// Run in the specified mode
	if *modePtr == "server" {
		// Start the server
		runServer(*portPtr)
	} else if *modePtr == "test" {
		// Run a test search
		runTest(*targetPtr, *algoPtr)
	} else {
		log.Fatalf("Invalid mode: %s. Use 'server' or 'test'", *modePtr)
	}
}

func runServer(port string) {
	// Set up API routes
	mux := http.NewServeMux()
	mux.HandleFunc("/api/search", api.SearchHandler)
	mux.HandleFunc("/api/elements", api.ElementsHandler)
	mux.HandleFunc("/api/elements/basic", api.BasicElementsHandler)

	// Serve static files for the frontend
	workDir, _ := os.Getwd()
	staticDir := filepath.Join(workDir, "static")
	fileServer := http.FileServer(http.Dir(staticDir))
	mux.Handle("/", fileServer)

	// Start the server
	addr := ":" + port
	log.Printf("Server started on http://localhost%s", addr)
	log.Printf("API endpoints: /api/search, /api/elements, /api/elements/basic")
	log.Fatal(http.ListenAndServe(addr, mux))
}

func runTest(target, algorithm string) {
	log.Printf("Running test search for target '%s' using algorithm '%s'", target, algorithm)

	// Load recipes from JSON file
	workDir, _ := os.Getwd()
	recipesPath := filepath.Join(workDir, "data", "recipes.json")
	recipes, err := api.LoadRecipesFromJSON(recipesPath)
	if err != nil {
		log.Fatalf("Failed to load recipes: %v", err)
	}

	log.Printf("Loaded %d recipes from %s", len(recipes), recipesPath)

	// Convert to algorithm format
	algoRecipes := api.ConvertToAlgoFormat(recipes)

	// Define start elements
	startElements := []string{"Air", "Earth", "Fire", "Water"}

	fmt.Printf("Searching for '%s' starting with %v...\n", target, startElements)

	if algorithm == "bfs" {
		path, nodesVisited, success := api.RunBFSTest(startElements, target, algoRecipes)
		if success {
			fmt.Printf("Path found: %v\n", path)
			fmt.Printf("Nodes visited: %d\n", nodesVisited)
		} else {
			fmt.Printf("No path found for target '%s'\n", target)
		}
	} else if algorithm == "dfs" {
		fmt.Println("DFS algorithm not implemented yet")
	} else {
		fmt.Printf("Unknown algorithm: %s\n", algorithm)
	}
}
