package searchalgo

import (
	"fmt"
	"sync"
	"tubes2/utilities"
)

// BiDirectionalSearch implements the bidirectional search algorithm for recipe finding
// It searches from both the target element and the base elements simultaneously
func BiDirectionalSearch(target string, maxRecipes int) ([]utilities.RecipeTree, int) {
	counter := &SafeCounter{v: 0}
	
	// If target is a base element, return it directly
	if utilities.IsBaseElement(target) {
		tree := utilities.RecipeTree{Element: target}
		return []utilities.RecipeTree{tree}, 0
	}
	
	// Check if recipes exist for the target
	recipeList, exists := utilities.Recipes[target]
	if !exists || len(recipeList) == 0 {
		fmt.Printf("Target element '%s' doesn't exist or can't be created\n", target)
		return nil, 0
	}
	
	var allResults []utilities.RecipeTree
	var mutex sync.Mutex
	
	// Forward search (starting from target) data structures
	forwardQueue := make([]*utilities.Node, 0)
	forwardVisited := make(map[string]map[string][]string) // element -> recipe map
	
	// Backward search (starting from base elements) data structures
	backwardQueue := make([]*utilities.Node, 0)
	backwardVisited := make(map[string]map[string][]string) // element -> recipe map
	
	// Initialize forward search with the target
	startNode := &utilities.Node{
		Element:    target,
		Path:       []string{target},
		Visited:    make(map[string]bool),
		Depth:      0,
		Ingredients: make(map[string][]string),
	}
	startNode.Visited[target] = true
	forwardQueue = append(forwardQueue, startNode)
	
	// Initialize backward search with base elements
	for _, baseElement := range utilities.BaseElements {
		baseNode := &utilities.Node{
			Element:    baseElement,
			Path:       []string{baseElement},
			Visited:    make(map[string]bool),
			Depth:      0,
			Ingredients: make(map[string][]string),
		}
		baseNode.Visited[baseElement] = true
		backwardQueue = append(backwardQueue, baseNode)
		
		// Add base elements to backward visited
		if backwardVisited[baseElement] == nil {
			backwardVisited[baseElement] = make(map[string][]string)
		}
	}
	
	// Track the number of recipes found
	recipesFound := 0
	
	// Continue searching until either queue is empty or max recipes found
	for len(forwardQueue) > 0 && len(backwardQueue) > 0 && (maxRecipes <= 0 || recipesFound < maxRecipes) {
		// Forward search step
		if len(forwardQueue) > 0 {
			currentNode := forwardQueue[0]
			forwardQueue = forwardQueue[1:]
			counter.Inc()
			
			// Check if this element has recipes
			element := currentNode.Element
			if !utilities.IsBaseElement(element) {
				recipeList, exists := utilities.Recipes[element]
				if exists {
					for _, recipe := range recipeList {
						e1 := recipe.Element1
						e2 := recipe.Element2
						
						// Create recipe map for this element
						if forwardVisited[element] == nil {
							forwardVisited[element] = make(map[string][]string)
						}
						forwardVisited[element][fmt.Sprintf("%s+%s", e1, e2)] = []string{e1, e2}
						
						// Track for live update visualization
						currentIngredients := make(map[string][]string)
						currentIngredients[element] = []string{e1, e2}
						utilities.TrackLiveUpdate(element, currentNode.Path, currentIngredients)
						
						// Check if the ingredients are in the backward search
						if checkBackwardConnection(e1, e2, backwardVisited, forwardVisited, element, currentNode, &allResults, &recipesFound, maxRecipes, &mutex) {
							if maxRecipes > 0 && recipesFound >= maxRecipes {
								break
							}
						}
						
						// Add new nodes to the forward queue
						addToForwardQueue(e1, e2, currentNode, &forwardQueue)
					}
				}
			}
		}
		
		// Check if we've found enough recipes
		if maxRecipes > 0 && recipesFound >= maxRecipes {
			break
		}
		
		// Backward search step
		if len(backwardQueue) > 0 {
			currentNode := backwardQueue[0]
			backwardQueue = backwardQueue[1:]
			counter.Inc()
			
			// In backward search, we look for recipes where this element is an ingredient
			element := currentNode.Element
			
			// Search all recipes to find those where this element is an ingredient
			for resultElement, recipes := range utilities.Recipes {
				for _, recipe := range recipes {
					if recipe.Element1 == element || recipe.Element2 == element {
						// We found a recipe where current element is an ingredient
						otherIngredient := recipe.Element1
						if recipe.Element1 == element {
							otherIngredient = recipe.Element2
						}
						
						// Create a new node for the result element
						newPath := append([]string{}, currentNode.Path...)
						newPath = append(newPath, resultElement)
						
						newVisited := make(map[string]bool)
						for k, v := range currentNode.Visited {
							newVisited[k] = v
						}
						newVisited[resultElement] = true
						
						newIngredients := make(map[string][]string)
						for k, v := range currentNode.Ingredients {
							newIngredients[k] = append([]string{}, v...)
						}
						newIngredients[resultElement] = []string{element, otherIngredient}
						
						// Create backward visited entry
						if backwardVisited[resultElement] == nil {
							backwardVisited[resultElement] = make(map[string][]string)
						}
						backwardVisited[resultElement][fmt.Sprintf("%s+%s", element, otherIngredient)] = []string{element, otherIngredient}
						
						// Track for live update visualization
						currentIngredients := make(map[string][]string)
						currentIngredients[resultElement] = []string{element, otherIngredient}
						utilities.TrackLiveUpdate(resultElement, newPath, currentIngredients)
						
						// Check if this result element is in the forward search
						if forwardVisited[resultElement] != nil {
							// We found a connection!
							for _, ingredients := range forwardVisited[resultElement] {
								// Create a complete recipe map
								completeRecipe := buildCompleteRecipe(resultElement, ingredients, element, otherIngredient, forwardVisited, backwardVisited)
								
								// Build recipe tree
								recipeTree := utilities.BuildRecipeTree(target, completeRecipe)
								
								mutex.Lock()
								allResults = append(allResults, recipeTree)
								recipesFound++
								mutex.Unlock()
								
								fmt.Printf("Found valid bidirectional recipe #%d for %s\n", recipesFound, target)
								
								if maxRecipes > 0 && recipesFound >= maxRecipes {
									break
								}
							}
						}
						
						// Check if this element is already visited in backward search
						if !newVisited[otherIngredient] {
							// Add the other ingredient to the backward queue if not a base element
							if !utilities.IsBaseElement(otherIngredient) {
								newNode := &utilities.Node{
									Element:    otherIngredient,
									Path:       append([]string{}, newPath...),
									Visited:    newVisited,
									Depth:      currentNode.Depth + 1,
									Ingredients: newIngredients,
								}
								newNode.Visited[otherIngredient] = true
								backwardQueue = append(backwardQueue, newNode)
							} else {
								// If it's a base element, add it to visited directly
								if backwardVisited[otherIngredient] == nil {
									backwardVisited[otherIngredient] = make(map[string][]string)
								}
							}
						}
						
						// If we've found enough recipes, break out
						if maxRecipes > 0 && recipesFound >= maxRecipes {
							break
						}
					}
				}
				
				// If we've found enough recipes, break out
				if maxRecipes > 0 && recipesFound >= maxRecipes {
					break
				}
			}
		}
	}
	
	fmt.Printf("Bidirectional search complete. Found %d recipes for %s\n", len(allResults), target)
	return allResults, counter.Value()
}

// Helper function to check for connections between forward and backward search
func checkBackwardConnection(e1 string, e2 string, backwardVisited, forwardVisited map[string]map[string][]string, 
                            element string, currentNode *utilities.Node, allResults *[]utilities.RecipeTree, 
                            recipesFound *int, maxRecipes int, mutex *sync.Mutex) bool {
	// Check if e1 is in backward search
	foundConnection := false
	if backwardVisited[e1] != nil {
		// Found a connection through e1
		completeRecipe := buildCompleteRecipe(element, []string{e1, e2}, e1, "", forwardVisited, backwardVisited)
		recipeTree := utilities.BuildRecipeTree(element, completeRecipe)
		
		mutex.Lock()
		*allResults = append(*allResults, recipeTree)
		*recipesFound++
		mutex.Unlock()
		
		fmt.Printf("Found valid bidirectional recipe #%d for %s (through %s)\n", *recipesFound, element, e1)
		foundConnection = true
	}
	
	// Check if e2 is in backward search
	if backwardVisited[e2] != nil && (maxRecipes <= 0 || *recipesFound < maxRecipes) {
		// Found a connection through e2
		completeRecipe := buildCompleteRecipe(element, []string{e1, e2}, e2, "", forwardVisited, backwardVisited)
		recipeTree := utilities.BuildRecipeTree(element, completeRecipe)
		
		mutex.Lock()
		*allResults = append(*allResults, recipeTree)
		*recipesFound++
		mutex.Unlock()
		
		fmt.Printf("Found valid bidirectional recipe #%d for %s (through %s)\n", *recipesFound, element, e2)
		foundConnection = true
	}
	
	return foundConnection
}

// Helper function to add elements to the forward queue
func addToForwardQueue(e1 string, e2 string, currentNode *utilities.Node, forwardQueue *[]*utilities.Node) {
	// Add e1 to queue if not visited and not a base element
	if !currentNode.Visited[e1] && !utilities.IsBaseElement(e1) {
		newVisited := make(map[string]bool)
		for k, v := range currentNode.Visited {
			newVisited[k] = v
		}
		newVisited[e1] = true
		
		newPath := append([]string{}, currentNode.Path...)
		newPath = append(newPath, e1)
		
		newIngredients := make(map[string][]string)
		for k, v := range currentNode.Ingredients {
			newIngredients[k] = append([]string{}, v...)
		}
		
		e1Node := &utilities.Node{
			Element:    e1,
			Path:       newPath,
			Visited:    newVisited,
			Depth:      currentNode.Depth + 1,
			Ingredients: newIngredients,
		}
		*forwardQueue = append(*forwardQueue, e1Node)
	}
	
	// Add e2 to queue if not visited and not a base element
	if !currentNode.Visited[e2] && !utilities.IsBaseElement(e2) {
		newVisited := make(map[string]bool)
		for k, v := range currentNode.Visited {
			newVisited[k] = v
		}
		newVisited[e2] = true
		
		newPath := append([]string{}, currentNode.Path...)
		newPath = append(newPath, e2)
		
		newIngredients := make(map[string][]string)
		for k, v := range currentNode.Ingredients {
			newIngredients[k] = append([]string{}, v...)
		}
		
		e2Node := &utilities.Node{
			Element:    e2,
			Path:       newPath,
			Visited:    newVisited,
			Depth:      currentNode.Depth + 1,
			Ingredients: newIngredients,
		}
		*forwardQueue = append(*forwardQueue, e2Node)
	}
}

// Helper function to build a complete recipe map when a connection is found
func buildCompleteRecipe(element string, forwardIngredients []string, backwardElement string, 
                         otherIngredient string, forwardVisited, backwardVisited map[string]map[string][]string) map[string][]string {
	completeRecipe := make(map[string][]string)
	
	// Add the current element's ingredients
	completeRecipe[element] = forwardIngredients
	
	// Walk forward from backwardElement to base elements
	queue := []string{backwardElement}
	if otherIngredient != "" {
		queue = append(queue, otherIngredient)
	}
	
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		
		// Skip if already processed or is a base element
		if completeRecipe[curr] != nil || utilities.IsBaseElement(curr) {
			continue
		}
		
		// Look for ingredients in backward visited
		if backwardVisited[curr] != nil {
			// Just take the first recipe we find
			for _, ingredients := range backwardVisited[curr] {
				completeRecipe[curr] = ingredients
				
				// Add ingredients to queue if not base elements
				for _, ing := range ingredients {
					if !utilities.IsBaseElement(ing) {
						queue = append(queue, ing)
					}
				}
				
				break
			}
		}
	}
	
	// Walk backward from forwardIngredients to target
	for _, ing := range forwardIngredients {
		if !utilities.IsBaseElement(ing) && completeRecipe[ing] == nil {
			// Breadth-first traversal of forward visited
			ingQueue := []string{ing}
			
			for len(ingQueue) > 0 {
				curr := ingQueue[0]
				ingQueue = ingQueue[1:]
				
				// Skip if already processed
				if completeRecipe[curr] != nil {
					continue
				}
				
				// Look for ingredients in forward visited
				if forwardVisited[curr] != nil {
					// Just take the first recipe we find
					for _, ingredients := range forwardVisited[curr] {
						completeRecipe[curr] = ingredients
						
						// Add ingredients to queue if not base elements
						for _, subIng := range ingredients {
							if !utilities.IsBaseElement(subIng) && completeRecipe[subIng] == nil {
								ingQueue = append(ingQueue, subIng)
							}
						}
						
						break
					}
				}
			}
		}
	}
	
	return completeRecipe
}