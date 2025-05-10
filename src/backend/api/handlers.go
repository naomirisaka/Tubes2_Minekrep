package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"tubes2/backend/searchalgo"
)

// Recipe represents a recipe from recipes.json
type Recipe struct {
	Element1     string `json:"element1"`
	Element2     string `json:"element2"`
	Result       string `json:"result"`
	IconFilename string `json:"icon_filename"`
}

// SearchRequest represents the search parameters from frontend
type SearchRequest struct {
	Algorithm       string   `json:"algorithm"`
	TargetElement   string   `json:"targetElement"`
	MultipleRecipes bool     `json:"multipleRecipes"`
	RecipeCount     int      `json:"recipeCount"`
	StartElements   []string `json:"startElements,omitempty"`
}

// ResultStep represents a step in the recipe path
type ResultStep struct {
	Element1     string `json:"element1"`
	Element2     string `json:"element2"`
	Result       string `json:"result"`
	IconFilename string `json:"icon_filename"`
}

// RecipeResult represents a single recipe result for the frontend
type RecipeResult struct {
	Path            []string     `json:"path,omitempty"`
	Steps           []ResultStep `json:"steps"`
	TargetElement   string       `json:"targetElement"`
	StartingElement string       `json:"startingElement"`
}

// SearchResult represents the response to the frontend
type SearchResult struct {
	Success bool           `json:"success"`
	Recipes []RecipeResult `json:"recipes"`
	Metrics struct {
		Time         float64 `json:"time"`
		NodesVisited int     `json:"nodesVisited"`
	} `json:"metrics"`
	LiveUpdateSteps []LiveUpdateStep `json:"liveUpdateSteps,omitempty"`
}

// LiveUpdateStep represents a step in the live update visualization
type LiveUpdateStep struct {
	CurrentElement string        `json:"currentElement"`
	Path           []string      `json:"path"`
	Available      []string      `json:"available"`
	Steps          []ResultStep  `json:"steps"`
	NewNodes       []interface{} `json:"newNodes"`
	Step           int           `json:"step"`
}

// LoadRecipesFromJSON loads recipes from the JSON file
func LoadRecipesFromJSON(filePath string) ([]Recipe, error) {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var recipes []Recipe
	err = json.Unmarshal(byteValue, &recipes)
	if err != nil {
		return nil, err
	}

	return recipes, nil
}

// ConvertToAlgoFormat converts recipes to the format used by the search algorithm
func ConvertToAlgoFormat(recipes []Recipe) searchalgo.Recipe {
	result := make(searchalgo.Recipe)

	for _, recipe := range recipes {
		combo := []string{recipe.Element1, recipe.Element2}

		if result[recipe.Result] == nil {
			result[recipe.Result] = [][]string{}
		}

		result[recipe.Result] = append(result[recipe.Result], combo)
	}

	return result
}

// FindRecipe searches for a recipe that matches the given elements
func FindRecipe(recipes []Recipe, element1, element2, result string) *Recipe {
	for _, recipe := range recipes {
		if (recipe.Element1 == element1 && recipe.Element2 == element2 && recipe.Result == result) ||
			(recipe.Element1 == element2 && recipe.Element2 == element1 && recipe.Result == result) {
			return &recipe
		}
	}
	return nil
}

// BuildRecipePath constructs the recipe steps from the path returned by the search algorithm
func BuildRecipePath(path []string, recipes []Recipe, algoRecipes searchalgo.Recipe) []ResultStep {
	var steps []ResultStep

	if len(path) <= 1 {
		return steps // No steps if only one element in path or empty
	}

	// For each step in the path, find the recipe that produces the next element
	for i := 0; i < len(path)-1; i++ {
		currentElement := path[i]
		nextElement := path[i+1]

		// Find the combo that produces nextElement
		if combos, exists := algoRecipes[nextElement]; exists {
			for _, combo := range combos {
				// Check if the current element is in this combo
				if combo[0] == currentElement || combo[1] == currentElement {
					// Found a matching combo
					otherElement := combo[0]
					if otherElement == currentElement {
						otherElement = combo[1]
					}

					// Find the corresponding recipe
					foundRecipe := FindRecipe(recipes, currentElement, otherElement, nextElement)
					if foundRecipe != nil {
						steps = append(steps, ResultStep{
							Element1:     foundRecipe.Element1,
							Element2:     foundRecipe.Element2,
							Result:       foundRecipe.Result,
							IconFilename: foundRecipe.IconFilename,
						})
					} else {
						// If recipe not found in original data, create one with default values
						steps = append(steps, ResultStep{
							Element1:     currentElement,
							Element2:     otherElement,
							Result:       nextElement,
							IconFilename: strings.ToLower(nextElement) + ".png",
						})
					}
					break
				}
			}
		}
	}

	return steps
}

// BuildRecipeFromString builds recipe steps from recipe strings (e.g., "Fire + Water => Steam")
func BuildRecipeFromString(recipeStrings []string, allRecipes []Recipe) []ResultStep {
	var steps []ResultStep

	for _, recipeStr := range recipeStrings {
		var elem1, elem2, result string
		parts := strings.Split(recipeStr, " => ")
		if len(parts) != 2 {
			log.Printf("Error parsing recipe step: invalid format %s", recipeStr)
			continue
		}

		ingredients := strings.Split(parts[0], " + ")
		if len(ingredients) != 2 {
			log.Printf("Error parsing recipe step: invalid ingredients format %s", parts[0])
			continue
		}

		elem1 = ingredients[0]
		elem2 = ingredients[1]
		result = parts[1]

		// Try to find the icon filename from allRecipes
		iconFilename := strings.ToLower(result) + ".png" // Default
		for _, recipe := range allRecipes {
			if (recipe.Element1 == elem1 && recipe.Element2 == elem2) ||
				(recipe.Element1 == elem2 && recipe.Element2 == elem1) {
				if recipe.Result == result {
					iconFilename = recipe.IconFilename
					break
				}
			}
		}

		steps = append(steps, ResultStep{
			Element1:     elem1,
			Element2:     elem2,
			Result:       result,
			IconFilename: iconFilename,
		})
	}

	return steps
}

// SearchHandler handles the /api/search endpoint
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var searchReq SearchRequest
	err := json.NewDecoder(r.Body).Decode(&searchReq)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Load recipes from JSON file
	workDir, _ := os.Getwd()
	recipesPath := filepath.Join(workDir, "data", "recipes.json")
	recipes, err := LoadRecipesFromJSON(recipesPath)
	if err != nil {
		http.Error(w, "Failed to load recipes: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Error loading recipes: %v", err)
		return
	}

	algoRecipes := ConvertToAlgoFormat(recipes)

	// Define basic elements if not provided
	if len(searchReq.StartElements) == 0 {
		searchReq.StartElements = []string{"Air", "Earth", "Fire", "Water"}
	}

	startTime := time.Now()
	var result SearchResult
	result.Success = true

	// Execute the appropriate search algorithm
	if searchReq.Algorithm == "bfs" {
		if searchReq.MultipleRecipes {
			// Multiple recipes search with BFS
			maxWorkers := 4      // Adjust based on your server capacity
			timeoutSeconds := 30 // Maximum execution time
			recipePaths, nodesVisited, success := searchalgo.BFSMultiple(
				searchReq.StartElements,
				searchReq.TargetElement,
				algoRecipes,
				searchReq.RecipeCount,
				maxWorkers,
				timeoutSeconds,
			)

			if !success {
				result.Success = false
			} else {
				for i, path := range recipePaths {
					if i >= searchReq.RecipeCount {
						break
					}

					steps := BuildRecipeFromString(path, recipes)

					startElement := searchReq.StartElements[0]
					if len(steps) > 0 {
						startElement = steps[0].Element1
					}

					recipe := RecipeResult{
						Steps:           steps,
						TargetElement:   searchReq.TargetElement,
						StartingElement: startElement,
					}

					result.Recipes = append(result.Recipes, recipe)
				}

				result.Metrics.NodesVisited = nodesVisited
			}
		} else {
			// Single shortest path search with BFS
			path, nodesVisited, success := searchalgo.BFSSingle(
				searchReq.StartElements,
				searchReq.TargetElement,
				algoRecipes,
			)

			if !success || len(path) == 0 {
				result.Success = false
			} else {
				// Build the recipe steps from the path
				steps := BuildRecipePath(path, recipes, algoRecipes)

				recipe := RecipeResult{
					Path:            path,
					Steps:           steps,
					TargetElement:   searchReq.TargetElement,
					StartingElement: path[0],
				}

				result.Recipes = append(result.Recipes, recipe)
				result.Metrics.NodesVisited = nodesVisited
			}
		}
	} else if searchReq.Algorithm == "dfs" {
		// Add DFS implementation here if available
		http.Error(w, "DFS algorithm not implemented yet", http.StatusNotImplemented)
		return
	} else {
		// Invalid algorithm
		http.Error(w, "Unsupported algorithm", http.StatusBadRequest)
		return
	}

	result.Metrics.Time = float64(time.Since(startTime).Milliseconds())
	result.Success = len(result.Recipes) > 0

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// ElementsHandler handles the /api/elements endpoint
func ElementsHandler(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Load recipes from JSON file to extract all available elements
	workDir, _ := os.Getwd()
	recipesPath := filepath.Join(workDir, "data", "recipes.json")
	recipes, err := LoadRecipesFromJSON(recipesPath)
	if err != nil {
		http.Error(w, "Failed to load recipes", http.StatusInternalServerError)
		log.Printf("Error loading recipes: %v", err)
		return
	}

	// Extract all unique elements
	elementMap := make(map[string]bool)
	for _, recipe := range recipes {
		elementMap[recipe.Element1] = true
		elementMap[recipe.Element2] = true
		elementMap[recipe.Result] = true
	}

	// Convert to slice
	var elements []string
	for element := range elementMap {
		elements = append(elements, element)
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(elements)
}

// BasicElementsHandler handles the /api/elements/basic endpoint
func BasicElementsHandler(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Hard-coded basic elements
	basicElements := []string{"Air", "Earth", "Fire", "Water"}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(basicElements)
}

// RunBFSTest runs a BFS test search
func RunBFSTest(startElements []string, target string, algoRecipes searchalgo.Recipe) ([]string, int, bool) {
	fmt.Printf("Running BFS search for target '%s'...\n", target)
	return searchalgo.BFSSingle(startElements, target, algoRecipes)
}
