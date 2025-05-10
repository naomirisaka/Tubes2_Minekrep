package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"fmt"
	"tubes2_minekrep/src/backend/searchalgo"
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

func BuildRecipeMap(combos []Combination) searchalgo.Recipes {
	recipe := make(searchalgo.Recipes)
	for _, c := range combos {
		ingredients := []string{c.Element1, c.Element2}
		recipe[c.Result] = append(recipe[c.Result], ingredients)
	}
	return recipe
}

func main() {
    combos := LoadCombinations("C:/Users/62812/Stima/Tubes2_Minekrep/src/data/recipes.json")
    recipes := BuildRecipeMap(combos)

    fmt.Printf("Recipes: %v\n", recipes)

    startElements := []string{"Water", "Earth", "Fire", "Air"}
    target := "Brick" 

    tiers := searchalgo.CalculateTiers(recipes, startElements)

    fmt.Println("\n=== DFS Test ===")
    pathDFS, stepsDFS, foundDFS := searchalgo.DFSSingle(startElements, target, recipes, tiers)
    if foundDFS {
        fmt.Printf("Path to target (DFS): %v\n", pathDFS)
        fmt.Printf("Steps taken (DFS): %d\n", stepsDFS)
    } else {
        fmt.Println("Target not found (DFS).")
    }
}