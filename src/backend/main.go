package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
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
	combos := LoadCombinations("../../data/recipes.json")
	recipeMap := BuildRecipeMap(combos)

	startElements := []string{"Water", "Earth", "Fire", "Air"}
	target := "Brick" // change to desired target

	searchalgo.BFSSingle(startElements, target, recipeMap)
	// searchalgo.BFSMultiple(startElements, target, recipeMap, 3)
}
