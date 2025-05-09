package main

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
	"log"
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

func BuildRecipeMap(combos []Combination) searchalgo.Recipe {
	recipe := make(searchalgo.Recipe)
	for _, c := range combos {
		ingredients := []string{c.Element1, c.Element2}
		recipe[c.Result] = append(recipe[c.Result], ingredients)
	}
	return recipe
}

func main() {
    combos := LoadCombinations("C:/Users/62812/Stima/Tubes2_Minekrep/src/data/recipes.json")
    recipeMap := BuildRecipeMap(combos)

    startElements := []string{"Water", "Earth", "Fire", "Air"}
    target := "Brick" // Ganti dengan target yang diinginkan

    // Panggil fungsi DFSSingle
    path, steps, found := searchalgo.DFSSingle(startElements, target, recipeMap)
    if found {
        fmt.Printf("Path to target: %v\n", path)
        fmt.Printf("Steps taken: %d\n", steps)
    } else {
        fmt.Println("Target not found.")
    }
}
