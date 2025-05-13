package main

import (
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
}
