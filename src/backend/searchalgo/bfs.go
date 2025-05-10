package searchalgo

import (
	"fmt"
	"tubes2/utilities"
)

func BFSSearch(target string, maxRecipes int) ([]utilities.RecipeTree, int) {
	visited := 0

	// Check if target is a base element
	if utilities.IsBaseElement(target) {
		tree := utilities.RecipeTree{Element: target}
		return []utilities.RecipeTree{tree}, visited
	}

	// Check if target exists in recipes
	recipeList, exists := utilities.Recipes[target]
	if !exists || len(recipeList) == 0 {
		fmt.Printf("Target element '%s' doesn't exist or can't be created\n", target)
		return nil, visited
	}

	var allResults []utilities.RecipeTree
	foundCount := 0

	// Try each recipe that can create the target
	for _, recipe := range recipeList {
		// Skip if we've already found enough recipes
		if maxRecipes > 0 && foundCount >= maxRecipes {
			break
		}

		e1 := recipe.Element1
		e2 := recipe.Element2

		// Create a new recipe path
		found := make(map[string][]string)
		found[target] = []string{e1, e2}

		// Check if this recipe can be completed
		visitCount := 0
		if processRecipe(e1, e2, found, &visitCount) {
			visited += visitCount
			recipeTree := utilities.BuildRecipeTree(target, found)
			allResults = append(allResults, recipeTree)
			foundCount++
		} else {
			visited += visitCount
		}
	}

	// Provide information if we couldn't find the requested number of recipes
	if maxRecipes > 0 && foundCount < maxRecipes {
		fmt.Printf("Note: Only found %d recipe(s) for '%s' while %d were requested.\n", 
			foundCount, target, maxRecipes)
	}

	return allResults, visited
}

// processRecipe checks if a recipe can be completed with available elements
func processRecipe(e1 string, e2 string, found map[string][]string, visitCount *int) bool {
	// Initialize a queue for BFS
	queue := []string{}
	
	// Add the ingredients to the queue if they need to be processed
	if !utilities.IsBaseElement(e1) && found[e1] == nil {
		queue = append(queue, e1)
	}
	
	if !utilities.IsBaseElement(e2) && found[e2] == nil {
		queue = append(queue, e2)
	}
	
	// If no ingredients need processing, we're done
	if len(queue) == 0 {
		return true
	}
	
	// Process the queue
	for len(queue) > 0 {
		element := queue[0]
		queue = queue[1:]
		*visitCount++
		
		// Skip if already processed
		if found[element] != nil {
			continue
		}
		
		// Check if this element has recipes
		recipeList, exists := utilities.Recipes[element]
		if !exists {
			// No recipes for this element - can't complete the chain
			return false
		}
		
		// Try each recipe for this element
		elementProcessed := false
		for _, recipe := range recipeList {
			ing1 := recipe.Element1
			ing2 := recipe.Element2
			
			// Record this recipe
			found[element] = []string{ing1, ing2}
			
			// Add non-base, unfound ingredients to queue
			newElements := []string{}
			if !utilities.IsBaseElement(ing1) && found[ing1] == nil {
				newElements = append(newElements, ing1)
			}
			
			if !utilities.IsBaseElement(ing2) && found[ing2] == nil {
				newElements = append(newElements, ing2)
			}
			
			// If this recipe doesn't add new elements, we're done with this element
			if len(newElements) == 0 {
				elementProcessed = true
				break
			}
			
			// Try with this recipe - make a copy of found for backtracking
			foundCopy := make(map[string][]string)
			for k, v := range found {
				foundCopy[k] = v
			}
			
			// Add new elements to queue
			newQueue := append([]string{}, queue...)
			newQueue = append(newQueue, newElements...)
			
			// Try to continue with this recipe
			allProcessed := true
			visitIncrease := 0
			
			// Process all remaining elements
			for len(newQueue) > 0 {
				nextElement := newQueue[0]
				newQueue = newQueue[1:]
				visitIncrease++
				
				// Skip if already processed
				if foundCopy[nextElement] != nil {
					continue
				}
				
				// Find recipes for this element
				nextRecipeList, exists := utilities.Recipes[nextElement]
				if !exists {
					allProcessed = false
					break
				}
				
				// Try to find a recipe that works
				elementResolved := false
				for _, nextRecipe := range nextRecipeList {
					nextIng1 := nextRecipe.Element1
					nextIng2 := nextRecipe.Element2
					
					// Check if we can use this recipe
					if (utilities.IsBaseElement(nextIng1) || foundCopy[nextIng1] != nil) &&
						(utilities.IsBaseElement(nextIng2) || foundCopy[nextIng2] != nil) {
						foundCopy[nextElement] = []string{nextIng1, nextIng2}
						elementResolved = true
						break
					}
					
					// If one ingredient is available, add the other to the queue
					if utilities.IsBaseElement(nextIng1) || foundCopy[nextIng1] != nil {
						if !utilities.IsBaseElement(nextIng2) && foundCopy[nextIng2] == nil {
							foundCopy[nextElement] = []string{nextIng1, nextIng2}
							newQueue = append(newQueue, nextIng2)
							elementResolved = true
							break
						}
					} else if utilities.IsBaseElement(nextIng2) || foundCopy[nextIng2] != nil {
						if !utilities.IsBaseElement(nextIng1) && foundCopy[nextIng1] == nil {
							foundCopy[nextElement] = []string{nextIng1, nextIng2}
							newQueue = append(newQueue, nextIng1)
							elementResolved = true
							break
						}
					}
				}
				
				if !elementResolved {
					allProcessed = false
					break
				}
			}
			
			*visitCount += visitIncrease
			
			if allProcessed {
				// This recipe path works - update the found map
				for k, v := range foundCopy {
					found[k] = v
				}
				elementProcessed = true
				break
			}
			
			// This recipe path didn't work, remove it and try the next one
			delete(found, element)
		}
		
		if !elementProcessed {
			// Couldn't process this element with any recipe
			return false
		}
	}
	
	// All elements processed
	return true
}

// BFSFindRecipeAll is a simpler version that just checks if a recipe can be made from base elements
func BFSFindRecipeAll(e1 string, e2 string, found map[string][]string, visitCount *int) bool {
	// Initialize a queue for BFS
	queue := []string{}
	
	// Add the ingredients to the queue if they need to be processed
	if !utilities.IsBaseElement(e1) && found[e1] == nil {
		queue = append(queue, e1)
	}
	
	if !utilities.IsBaseElement(e2) && found[e2] == nil {
		queue = append(queue, e2)
	}
	
	// If no ingredients need processing, we're done
	if len(queue) == 0 {
		return true
	}
	
	// Process the queue
	for len(queue) > 0 {
		element := queue[0]
		queue = queue[1:]
		*visitCount++
		
		// Skip if already processed
		if found[element] != nil {
			continue
		}
		
		// Check if this element has recipes
		recipeList, exists := utilities.Recipes[element]
		if !exists {
			// No recipes for this element - can't complete the chain
			return false
		}
		
		// Try each recipe for this element
		elementProcessed := false
		for _, recipe := range recipeList {
			ing1 := recipe.Element1
			ing2 := recipe.Element2
			
			// Record this recipe
			found[element] = []string{ing1, ing2}
			
			// Add non-base, unfound ingredients to queue
			if !utilities.IsBaseElement(ing1) && found[ing1] == nil {
				queue = append(queue, ing1)
			}
			
			if !utilities.IsBaseElement(ing2) && found[ing2] == nil {
				queue = append(queue, ing2)
			}
			
			elementProcessed = true
			break
		}
		
		if !elementProcessed {
			// Couldn't process this element with any recipe
			return false
		}
	}
	
	// Verify all elements in the found map can be created
	for elem, ingredients := range found {
		if utilities.IsBaseElement(elem) {
			continue
		}
		
		for _, ing := range ingredients {
			if !utilities.IsBaseElement(ing) && found[ing] == nil {
				return false
			}
		}
	}
	
	return true
}