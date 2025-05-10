package main

import (
	"fmt"
	"time"
	"tubes2/utilities"
	"tubes2/scraper"
	"tubes2/searchalgo"
)

func main() {
	recipesFile := "../../data/recipes.json" // Default file path
	targetElement := "Brick" 
	maxRecipes := 10
	startTime := time.Now()

	// Load recipes
	err := scraper.LoadRecipes(recipesFile)
	if err != nil {
		fmt.Printf("Error loading recipes: %v\n", err)
		return
	}

	fmt.Printf("Searching for recipes to create '%s'...\n", targetElement)
	fmt.Printf("Algorithm: BFS\n") 

	results, visitedNodes := searchalgo.BFSSearch(targetElement, maxRecipes)
	// results, visitedNodes := searchalgo.DFSSearch(targetElement, maxRecipes)

	duration := time.Since(startTime)

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

	fmt.Printf("Stats:\n")
	fmt.Printf("- Execution time: %.2f ms\n", float64(duration)/float64(time.Millisecond))
	fmt.Printf("- Nodes visited: %d\n", visitedNodes)
}
