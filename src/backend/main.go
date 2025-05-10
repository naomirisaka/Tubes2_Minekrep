package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

type Recipe struct {
	Element1     string `json:"element1"`
	Element2     string `json:"element2"`
	Result       string `json:"result"`
	IconFilename string `json:"icon_filename"`
}

type Element struct {
	Name string
	Tier int
}

type Node struct {
	Element    string
	Path       []string
	Visited    map[string]bool
	Depth      int
	Ingredients map[string][]string 
}

type RecipeTree struct {
	Element    string      `json:"element"`
	Ingredients []RecipeTree `json:"ingredients,omitempty"`
}

var (
	elements   = make(map[string]Element)
	recipes    = make(map[string][]Recipe)
	baseElements = []string{"Water", "Fire", "Earth", "Air"}
	tiers      = make(map[string]int)
	visited    = 0
)

func initializeTiers() {
	// Set base tier 1
	for _, element := range baseElements {
		tiers[element] = 1
	}
	queue := make([]string, 0)
    queue = append(queue, baseElements...)
    processed := make(map[string]bool)
    
    for _, elem := range baseElements {
        processed[elem] = true
    }
	for len(queue) > 0 {
        current := queue[0]
        queue = queue[1:]
        
        for result, recipeList := range recipes {
            if processed[result] {
                continue 
            }
            
            for _, recipe := range recipeList {
                if (recipe.Element1 == current || recipe.Element2 == current) {
                    // Only compute tier if both ingredients have tiers
                    if tier1, ok1 := tiers[recipe.Element1]; ok1 {
                        if tier2, ok2 := tiers[recipe.Element2]; ok2 {
                            resultTier := max(tier1, tier2) + 1
                            existingTier, exists := tiers[result]
                            
                            // Update tier kalau ada yagn lebih pendek
                            if !exists || resultTier < existingTier {
                                tiers[result] = resultTier
                                if !processed[result] {
                                    queue = append(queue, result)
                                }
                            }
                            
                            processed[result] = true
                        }
                    }
                }
            }
        }
    }
}

func loadRecipes(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	var recipeList []Recipe
	err = json.Unmarshal(data, &recipeList)
	if err != nil {
		return err
	}

	for _, recipe := range recipeList {
		recipes[recipe.Result] = append(recipes[recipe.Result], recipe)
	}

	initializeTiers()

	return nil
}

func isBaseElement(element string) bool {
	for _, base := range baseElements {
		if element == base {
			return true
		}
	}
	return false
}

func dfsSearch(target string, findShortest bool) ([]RecipeTree, int) {
    visited = 0
    if isBaseElement(target) {
        tree := RecipeTree{Element: target}
        return []RecipeTree{tree}, visited
    }

    if _, exists := recipes[target]; !exists {
        fmt.Printf("Target element '%s' doesn't exist or can't be created\n", target)
        return nil, visited
    }

    found := make(map[string][]string) // map[element][]ingredients
    visited = findRecipe(target, found)

    if len(found) > 0 {
        recipeTree := buildRecipeTree(target, found)
        return []RecipeTree{recipeTree}, visited
    }

    return nil, visited
}

func findRecipe(element string, found map[string][]string) int {
    count := 1
    
    if isBaseElement(element) || found[element] != nil {
        return count
    }
    
    recipeList, exists := recipes[element]
    if !exists {
        return count
    }
    
    for _, recipe := range recipeList {
        e1 := recipe.Element1
        e2 := recipe.Element2
        
        // skip tier lebih besar
        if tiers[e1] >= tiers[element] || tiers[e2] >= tiers[element] {
            continue
        }
        
        c1 := findRecipe(e1, found)
        c2 := findRecipe(e2, found)
        count += c1 + c2
        if (isBaseElement(e1) || found[e1] != nil) && 
           (isBaseElement(e2) || found[e2] != nil) {
            found[element] = []string{e1, e2}
            return count
        }
    }
    
    return count
}

func buildRecipeTree(element string, ingredients map[string][]string) RecipeTree {
	tree := RecipeTree{Element: element}
	
	if ingList, exists := ingredients[element]; exists {
		for _, ing := range ingList {
			tree.Ingredients = append(tree.Ingredients, buildRecipeTree(ing, ingredients))
		}
	}
	
	return tree
}

func calculateTreeDepth(tree RecipeTree) int {
	if len(tree.Ingredients) == 0 {
		return 1
	}
	
	maxDepth := 0
	for _, ing := range tree.Ingredients {
		depth := calculateTreeDepth(ing)
		if depth > maxDepth {
			maxDepth = depth
		}
	}
	
	return maxDepth + 1
}

// Print recipe tree
func printRecipeTree(tree RecipeTree, indent string) {
	fmt.Printf("%s%s\n", indent, tree.Element)
	if len(tree.Ingredients) > 0 {
		fmt.Printf("%s└─ combines:\n", indent)
		for i, ing := range tree.Ingredients {
			if i == len(tree.Ingredients)-1 {
				printRecipeTree(ing, indent+"   ")
			} else {
				printRecipeTree(ing, indent+"│  ")
			}
		}
	}
}

func main() {
    recipesFile := "src/data/recipes.json" // Default file path
    targetElement := "Diamond" 
	findShortest := false

	// Load recipes
	err := loadRecipes(recipesFile)
	if err != nil {
		fmt.Printf("Error loading recipes: %v\n", err)
		return
	}

	fmt.Printf("Searching for recipes to create '%s'...\n", targetElement)
	fmt.Printf("Algorithm: DFS\n")
	startTime := time.Now()

	results, visitedNodes := dfsSearch(targetElement, findShortest)
	duration := time.Since(startTime)

	// Print
	if len(results) > 0 {
		fmt.Printf("Found %d recipe(s) for '%s':\n\n", len(results), targetElement)
		for i, result := range results {
			fmt.Printf("Recipe %d:\n", i+1)
			printRecipeTree(result, "")
			fmt.Println()
		}
	} else {
		fmt.Printf("No recipes found for '%s'\n", targetElement)
	}

	// Print stats
	fmt.Printf("Stats:\n")
	fmt.Printf("- Execution time: %v\n", duration)
	fmt.Printf("- Nodes visited: %d\n", visitedNodes)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}