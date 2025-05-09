package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"runtime"
	"time"
	"tubes2/searchalgo"
)

type Combination struct {
	Element1     string `json:"element1"`
	Element2     string `json:"element2"`
	Result       string `json:"result"`
	IconFilename string `json:"icon_filename"`
}

func LoadCombinations(filepath string) []Combination {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Failed to read JSON: %v", err)
	}
	var combos []Combination
	if err := json.Unmarshal(data, &combos); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}
	return combos
}

func BuildRecipeMap(combos []Combination) searchalgo.Recipe {
	recipe := make(searchalgo.Recipe)
	for _, c := range combos {
		ingredients := []string{c.Element1, c.Element2}
		recipe[c.Result] = append(recipe[c.Result], ingredients)
	}
	return recipe
}

func main() {
	// flag for testing
	targetPtr := flag.String("target", "Vinegar", "Target element to search for") // change target here
	algorithmPtr := flag.String("algo", "bfs", "Search algorithm to use (bfs or dfs)")
	modePtr := flag.String("mode", "multiple", "Search mode (single or multiple)")
	maxRecipesPtr := flag.Int("max", 3, "Maximum number of recipes to find in multiple mode")
	workersPtr := flag.Int("workers", 0, "Number of worker goroutines") // 0 for auto
	timeoutPtr := flag.Int("timeout", 60, "Search timeout in seconds")
	flag.Parse()

	combos := LoadCombinations("../../data/recipes.json")
	recipeMap := BuildRecipeMap(combos)
	startElements := []string{"Water", "Earth", "Fire", "Air"}
	target := *targetPtr

	workers := *workersPtr
	if workers <= 0 {
		workers = runtime.NumCPU()
	}

	fmt.Printf("Searching for '%s' using %s algorithm in %s mode\n",
		target, *algorithmPtr, *modePtr)

	startTime := time.Now()

	if *algorithmPtr == "bfs" {
		if *modePtr == "single" {
			path, steps, found := searchalgo.BFSSingle(startElements, target, recipeMap)
			if found {
				fmt.Printf("Found path in %d steps\n", steps)
				for i, step := range path {
					fmt.Printf("%d. %s\n", i+1, step)
				}
			} else {
				fmt.Printf("No path found to %s\n", target)
			}
		} else {
			recipes, visited, success := searchalgo.BFSMultiple(
				startElements, target, recipeMap, *maxRecipesPtr,
				workers, *timeoutPtr)

			if success {
				fmt.Printf("\nRecipes found (%d):\n", len(recipes))
				for i, recipe := range recipes {
					fmt.Printf("\nRecipe %d:\n", i+1)
					for j, step := range recipe {
						fmt.Printf("  %d. %s\n", j+1, step)
					}
				}
				fmt.Printf("\nTotal nodes visited: %d\n", visited)
			} else {
				fmt.Println("Search timed out")
			}
		}
	} else if *algorithmPtr == "dfs" {
		fmt.Println("work in progress")
	}

	duration := time.Since(startTime)
	fmt.Printf("\nTotal execution time: %v\n", duration)
}
