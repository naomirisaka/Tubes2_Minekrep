package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

// Recipe represents a single recipe from the scraping result
type Recipe struct {
	Element1     string `json:"element1"`
	Element2     string `json:"element2"`
	Result       string `json:"result"`
	IconFilename string `json:"icon_filename"`
}

// Element represents an element with its tier
type Element struct {
	Name string
	Tier int
}

// Node represents a node in the DFS search
type Node struct {
	Element    string
	Path       []string
	Visited    map[string]bool
	Depth      int
	Ingredients map[string][]string // Key: element, Value: slice of elements needed to create it
}

// RecipeTree represents the tree structure for visualization
type RecipeTree struct {
	Element    string      `json:"element"`
	Ingredients []RecipeTree `json:"ingredients,omitempty"`
}

// Global variables
var (
	elements   = make(map[string]Element)
	recipes    = make(map[string][]Recipe)
	baseElements = []string{"Water", "Fire", "Earth", "Air"}
	tiers      = make(map[string]int)
	visited    = 0
)

// Initialize tiers for all elements
func initializeTiers() {
	// Set base elements to tier 1
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
        
        // Find all elements that can be created using the current element
        for result, recipeList := range recipes {
            if processed[result] {
                continue // Already processed this element
            }
            
            for _, recipe := range recipeList {
                if (recipe.Element1 == current || recipe.Element2 == current) {
                    // Only compute tier if both ingredients have tiers
                    if tier1, ok1 := tiers[recipe.Element1]; ok1 {
                        if tier2, ok2 := tiers[recipe.Element2]; ok2 {
                            resultTier := max(tier1, tier2) + 1
                            existingTier, exists := tiers[result]
                            
                            // Update tier if not exists or found a shorter path
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

	// Calculate tiers for all other elements based on their ingredients
	// changed := true
	// for changed {
	// 	changed = false
	// 	for result, recipeList := range recipes {
	// 		if _, exists := tiers[result]; !exists {
	// 			// Element doesn't have a tier yet
	// 			for _, recipe := range recipeList {
	// 				if tier1, ok1 := tiers[recipe.Element1]; ok1 {
	// 					if tier2, ok2 := tiers[recipe.Element2]; ok2 {
	// 						// Both ingredients have tiers, calculate tier for result
	// 						resultTier := max(tier1, tier2) + 1
	// 						tiers[result] = resultTier
	// 						changed = true
	// 						break
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}
	// }
    // Debug: Print tiers for key elements
    fmt.Println("\nDebug - Element Tiers (revised):")
    keysToCheck := []string{"Water", "Fire", "Earth", "Air", "Swamp", "Energy", "Love", "Time", "Life"}
    for _, key := range keysToCheck {
        tier, exists := tiers[key]
        fmt.Printf("- %s: tier %d (exists: %v)\n", key, tier, exists)
    }
}

// Load recipes from JSON file
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

	// Organize recipes by result element
	for _, recipe := range recipeList {
		recipes[recipe.Result] = append(recipes[recipe.Result], recipe)
	}
	
	fmt.Println("\nDebug - Important Recipes:")
    for key, recipeList := range recipes {
        if key == "Life" || key == "life" || key == "LIFE" {
            fmt.Printf("Found '%s' with %d recipes:\n", key, len(recipeList))
            for i, r := range recipeList {
                fmt.Printf("  %d. %s + %s => %s\n", i+1, r.Element1, r.Element2, r.Result)
            }
        }
    }

	// Initialize tiers
	initializeTiers()

	return nil
}

// Check if element is a base element
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

    // DFS to find the recipe
    found := make(map[string][]string) // map[element][]ingredients
    visited = findRecipe(target, found)

    if len(found) > 0 {
        recipeTree := buildRecipeTree(target, found)
        return []RecipeTree{recipeTree}, visited
    }

    return nil, visited
}

// Helper function for DFS
func findRecipe(element string, found map[string][]string) int {
    count := 1 // Count this visit
    
    // If it's a base element or we've already found a recipe for it
    if isBaseElement(element) || found[element] != nil {
        return count
    }
    
    // Get recipes for this element
    recipeList, exists := recipes[element]
    if !exists {
        return count
    }
    
    // Try each recipe
    for _, recipe := range recipeList {
        e1 := recipe.Element1
        e2 := recipe.Element2
        
        // Skip if violating tier rules
        if tiers[e1] >= tiers[element] || tiers[e2] >= tiers[element] {
            continue
        }
        
        // Try to find recipes for the ingredients
        c1 := findRecipe(e1, found)
        c2 := findRecipe(e2, found)
        count += c1 + c2
        
        // If both ingredients have recipes (or are base elements)
        if (isBaseElement(e1) || found[e1] != nil) && 
           (isBaseElement(e2) || found[e2] != nil) {
            // We found a valid recipe
            found[element] = []string{e1, e2}
            return count
        }
    }
    
    return count
}

// DFS search for a recipe
// func dfsSearch(target string, findShortest bool) ([]RecipeTree, int) {
	// visited = 0
	// if isBaseElement(target) {
	// 	tree := RecipeTree{Element: target}
	// 	return []RecipeTree{tree}, visited
	// }

	// if _, exists := recipes[target]; !exists {
	// 	fmt.Printf("Target element '%s' doesn't exist or can't be created\n", target)
	// 	return nil, visited
	// }

	// // Debug recipes for target
    // fmt.Printf("Found %d recipes for '%s':\n", len(recipes[target]), target)
    // for i, recipe := range recipes[target] {
    //     fmt.Printf("Recipe %d: %s + %s => %s\n", i+1, recipe.Element1, recipe.Element2, target)
        
    //     // Debug tier information
    //     e1Tier, e1Exists := tiers[recipe.Element1]
    //     e2Tier, e2Exists := tiers[recipe.Element2]
    //     targetTier, targetExists := tiers[target]
        
    //     fmt.Printf("Tiers: %s (%d, exists: %v), %s (%d, exists: %v), %s (%d, exists: %v)\n",
    //         recipe.Element1, e1Tier, e1Exists,
    //         recipe.Element2, e2Tier, e2Exists,
    //         target, targetTier, targetExists)
        
    //     // Check tier constraint
    //     if e1Tier >= targetTier || e2Tier >= targetTier {
    //         fmt.Printf("Tier constraint violated for this recipe\n")
    //     } else {
    //         fmt.Printf("This recipe passes tier check\n")
    //     }
    // }

	// // Initialize stack with the target element
	// stack := []Node{
	// 	{
	// 		Element:    target,
	// 		Path:       []string{target},
	// 		Visited:    make(map[string]bool),
	// 		Depth:      0,
	// 		Ingredients: make(map[string][]string),
	// 	},
	// }
	// stack[0].Visited[target] = true // Mark target as visited

	// var results []RecipeTree
	// // minDepth := 1000000 // Set to a large value

	// for len(stack) > 0 {
	// 	// Pop the last element from the stack
	// 	currentNode := stack[len(stack)-1]
	// 	stack = stack[:len(stack)-1]
		
	// 	visited++
	// 	fmt.Printf("Visiting: %s (depth: %d)\n", currentNode.Element, currentNode.Depth)

	// 	// Skip if we already found a shorter path and we're looking for the shortest
	// 	// if findShortest && currentNode.Depth > minDepth {
	// 	// 	continue
	// 	// }

	// 	currentElement := currentNode.Element
		
	// 	// If this is a base element, we can't break it down further
	// 	if isBaseElement(currentElement) {
	// 		fmt.Printf("  %s is a base element, continuing\n", currentElement)
	// 		continue
	// 	}

	// 	// Get recipes for current element
	// 	recipeList, exists := recipes[currentElement]
	// 	if !exists {
	// 		fmt.Printf("  No recipes for %s, continuing\n", currentElement)
	// 		continue
	// 	}

	// 	// Try each recipe for the current element
	// 	for _, recipe := range recipeList {
	// 		fmt.Printf("  Trying recipe: %s + %s => %s\n", 
    //              recipe.Element1, recipe.Element2, currentElement)
            
    //         // Debug tier information
    //         e1Tier, e1Exists := tiers[recipe.Element1]
    //         e2Tier, e2Exists := tiers[recipe.Element2]
    //         currTier, currExists := tiers[currentElement]
            
    //         fmt.Printf("  Tiers: %s (%d, exists: %v), %s (%d, exists: %v), %s (%d, exists: %v)\n",
    //             recipe.Element1, e1Tier, e1Exists,
    //             recipe.Element2, e2Tier, e2Exists,
    //             currentElement, currTier, currExists)

	// 		// Skip if we violate tier rules
	// 		 if e1Tier >= currTier || e2Tier >= currTier {
    //             fmt.Printf("  Skipping due to tier constraint\n")
    //             continue
    //         }
	// 		fmt.Printf("  Recipe passes tier check\n")
	// 		// if tiers[recipe.Element1] >= tiers[currentElement] || 
	// 		//    tiers[recipe.Element2] >= tiers[currentElement] {
	// 		// 	continue
	// 		// }

	// 		// Check for cycles
    //         if currentNode.Visited[recipe.Element1] || currentNode.Visited[recipe.Element2] {
    //             fmt.Printf("  Skipping to avoid cycle\n")
    //             continue
    //         }

	// 		// Add ingredient relationship
	// 		currentNode.Ingredients[currentElement] = []string{recipe.Element1, recipe.Element2}
	// 		fmt.Printf("  Added ingredients: %s needs [%s, %s]\n", 
    //             currentElement, recipe.Element1, recipe.Element2)

	// 		// Add both ingredients to the stack
	// 		newNode1 := Node{
	// 			Element:    recipe.Element1,
	// 			Path:       append(currentNode.Path, recipe.Element1),
	// 			Visited:    copyMap(currentNode.Visited),
	// 			Depth:      currentNode.Depth + 1,
	// 			Ingredients: copyIngredientsMap(currentNode.Ingredients),
	// 		}
	// 		newNode1.Visited[recipe.Element1] = true

	// 		newNode2 := Node{
	// 			Element:    recipe.Element2,
	// 			Path:       append(currentNode.Path, recipe.Element2),
	// 			Visited:    copyMap(currentNode.Visited),
	// 			Depth:      currentNode.Depth + 1,
	// 			Ingredients: copyIngredientsMap(currentNode.Ingredients),
	// 		}
	// 		newNode2.Visited[recipe.Element2] = true

	// 		// Check if this is a valid complete recipe (all leaves are base elements)
	// 		// isValid := true
	// 		fmt.Println("  Checking recipe validity:")
	// 		fmt.Println("  Current ingredients map:", newNode1.Ingredients)
	// 		// Debug ingredients map
	// 		for elem, ingList := range newNode1.Ingredients {
	// 			fmt.Printf("    %s needs: %v\n", elem, ingList)
	// 		}

	// 		if currentElement == target {
	// 			// Hanya cek validitas jika kita kembali ke target
	// 			isValid := true
	// 			// missingRecipe := false
				
	// 			fmt.Println("  Checking if recipe complete for target:", target)
				
	// 			// Fungsi rekursif untuk memvalidasi resep
	// 			var validateRecipe func(elem string) bool
	// 			validateRecipe = func(elem string) bool {
	// 				// Base case: jika elemen dasar, selalu valid
	// 				if isBaseElement(elem) {
	// 					fmt.Printf("    %s is a base element - valid\n", elem)
	// 					return true
	// 				}
					
	// 				// Cek apakah elemen memiliki resep dalam ingredients map
	// 				ingredients, exists := newNode1.Ingredients[elem]
	// 				if !exists || len(ingredients) == 0 {
	// 					fmt.Printf("    %s has no ingredients defined - invalid\n", elem)
	// 					return false
	// 				}
					
	// 				// Cek semua ingredients rekursif
	// 				for _, ing := range ingredients {
	// 					if !validateRecipe(ing) {
	// 						return false
	// 					}
	// 				}
					
	// 				return true
	// 			}
				
	// 			// Validasi resep untuk target
	// 			isValid = validateRecipe(target)
				
	// 			if isValid {
	// 				fmt.Printf("  Found valid complete recipe for %s!\n", target)
	// 				recipeTree := buildRecipeTree(target, newNode1.Ingredients)
	// 				results = append(results, recipeTree)
	// 				return results, visited
	// 			}
	// 		}

		// 	for element := range newNode1.Ingredients {
		// 		ingredients := newNode1.Ingredients[element]
		// 		for _, ing := range ingredients {
		// 			fmt.Printf("    Checking ingredient: %s\n", ing)
		// 			if !isBaseElement(ing) && len(newNode1.Ingredients[ing]) == 0 {
		// 				fmt.Printf("    INVALID: %s is not a base element and has no recipe in ingredients map\n", ing)
		// 				isValid = false
		// 				break
		// 			} else if isBaseElement(ing) {
		// 				fmt.Printf("    OK: %s is a base element\n", ing)
		// 			} else {
		// 				fmt.Printf("    OK: %s has a recipe in the ingredients map\n", ing)
		// 			}
		// 		}
		// 		if !isValid {
		// 			break
		// 		}
		// 	}

		// 	if isValid {
		// 		fmt.Printf("  Found valid recipe!\n")
		// 		// Build recipe tree
		// 		recipeTree := buildRecipeTree(target, newNode1.Ingredients)

		// 		results = append(results, recipeTree)
		// 		return results, visited
		// 	}

		// 	// Add to stack for further exploration if not base elements
		// 	if !isBaseElement(recipe.Element1) {
		// 		fmt.Printf("  Adding %s to stack for further exploration\n", recipe.Element1)
		// 		stack = append(stack, newNode1)
		// 	} else {
        //         fmt.Printf("  %s is a base element\n", recipe.Element1)
        //     }

		// 	if !isBaseElement(recipe.Element2) {
		// 		stack = append(stack, newNode2)
		// 	} else {
        //         fmt.Printf("  %s is a base element\n", recipe.Element2)
        //     }
// 		}
// 	}

// 	return results, visited
// }

// Build recipe tree for visualization
func buildRecipeTree(element string, ingredients map[string][]string) RecipeTree {
	tree := RecipeTree{Element: element}
	
	if ingList, exists := ingredients[element]; exists {
		for _, ing := range ingList {
			tree.Ingredients = append(tree.Ingredients, buildRecipeTree(ing, ingredients))
		}
	}
	
	return tree
}

// Calculate depth of recipe tree
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

// Copy a map
func copyMap(original map[string]bool) map[string]bool {
	copy := make(map[string]bool)
	for k, v := range original {
		copy[k] = v
	}
	return copy
}

// Copy ingredients map
func copyIngredientsMap(original map[string][]string) map[string][]string {
	copy := make(map[string][]string)
	for k, v := range original {
		newSlice := make([]string, len(v))
		for i, s := range v {
			newSlice[i] = s
		}
		copy[k] = newSlice
	}
	return copy
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

// Main function
func main() {
	var recipesFile string
	var targetElement string

	if len(os.Args) < 2 {
        recipesFile = "src/data/recipes.json" // Default file path
        fmt.Println("No recipe file specified, using default: src/data/recipes.json")
    } else {
        recipesFile = os.Args[1]
    }

	if len(os.Args) < 3 {
        targetElement = "Brick" // Default target element
        fmt.Println("No target element specified, using default")
    } else {
        targetElement = os.Args[2]
    }

	findShortest := false

	// Parse command line arguments
	// recipesFile := os.Args[1]
	// targetElement := os.Args[2]
	// findShortest := false
	// if len(os.Args) > 3 && os.Args[3] == "true" {
	// 	findShortest = true
	// }

	// Load recipes
	err := loadRecipes(recipesFile)
	if err != nil {
		fmt.Printf("Error loading recipes: %v\n", err)
		return
	}

	fmt.Printf("Searching for recipes to create '%s'...\n", targetElement)
	fmt.Printf("Algorithm: DFS\n")
	// fmt.Printf("Find shortest recipe: %v\n\n", findShortest)

	// Start timer
	startTime := time.Now()

	// Perform DFS search
	results, visitedNodes := dfsSearch(targetElement, findShortest)
	
	// Stop timer
	duration := time.Since(startTime)

	// Print results
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

// Helper function to get max of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}