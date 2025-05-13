package main

import (
<<<<<<< HEAD
	//"encoding/json"
	"fmt"
	//"io/ioutil"
	"time"
	"tubes2_minekrep/src/backend/utilities"
	// "tubes2_minekrep/src/backend/scraper"
	"tubes2_minekrep/src/backend/searchalgo"
)

func main() {
    recipesFile := "C:/Users/62812/Stima/Tubes2_Minekrep/data/recipes.json" // Default file path
    targetElement := "Steam" 
	maxRecipes := 1
	startTime := time.Now()
	// Load recipes
	err := utilities.LoadRecipes(recipesFile)
	if err != nil {
		fmt.Printf("Error loading recipes: %v\n", err)
		return
	}

	fmt.Printf("Searching for recipes to create '%s'...\n", targetElement)
	fmt.Printf("Algorithm: DFS\n")

	results, visitedNodes := searchalgo.DFSSearch(targetElement, maxRecipes)
	duration := time.Since(startTime)

	// Print
	if len(results) > 0 {
		fmt.Printf("Found %d recipe(s) for '%s':\n\n", len(results), targetElement)
		for i, result := range results {
			fmt.Printf("Recipe %d:\n", i+1)
			utilities.PrintRecipeTree(result, "")
			fmt.Println()
		}
	} else {
		fmt.Printf("No recipes found for '%s'\n", targetElement)
	}

	// Print stats
	fmt.Printf("Stats:\n")
	fmt.Printf("- Execution time: %.2f ms\n", float64(duration)/float64(time.Millisecond))
	fmt.Printf("- Nodes visited: %d\n", visitedNodes)
=======
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"tubes2/api"
	"tubes2/scraper"
	"tubes2/utilities"
)

func main() {
	// Ensure data directory exists
	workDir, _ := os.Getwd()
	recipesPath := filepath.Join(workDir, "data", "recipes.json")
	scraper.ScrapeIfNeeded(recipesPath)

	// Command line flags
	utilities.LoadRecipes("data/recipes.json")
	portPtr := flag.String("port", "8080", "Port for the server to listen on")
	modePtr := flag.String("mode", "server", "Mode to run (server or test)")
	flag.Parse()

	// if _, err := os.Stat(recipesPath); os.IsNotExist(err) {
	// 	log.Fatalf("Recipes file not found: %s\nMake sure the 'data' directory with 'recipes.json' exists", recipesPath)
	// }

	// Run in the specified mode
	if *modePtr == "server" {
		// Start the server
		runServer(*portPtr)
	} else {
		log.Fatalf("Invalid mode: %s. Use 'server'", *modePtr)
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
>>>>>>> dd6ca3248ae7b2d1452d4e3847b539e901bdfce9
}
