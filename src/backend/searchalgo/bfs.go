package searchalgo

import (
	"fmt"
	"tubes2/utilities"
)

func BFSSearch(target string, maxRecipes int) ([]utilities.RecipeTree, int) {
	visited := 0

	if utilities.IsBaseElement(target) {
		tree := utilities.RecipeTree{Element: target}
		return []utilities.RecipeTree{tree}, visited
	}

	recipeList, exists := utilities.Recipes[target]
	if !exists || len(recipeList) == 0 {
		fmt.Printf("Target element '%s' doesn't exist or can't be created\n", target)
		return nil, visited
	}

	var allResults []utilities.RecipeTree
	foundCount := 0

	for _, recipe := range recipeList {
		if maxRecipes > 0 && foundCount >= maxRecipes {
			break
		}

		e1 := recipe.Element1
		e2 := recipe.Element2

		found := make(map[string][]string)
		found[target] = []string{e1, e2}

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

	if maxRecipes > 0 && foundCount < maxRecipes {
		fmt.Printf("Note: Only found %d recipe(s) for '%s' while %d were requested.\n", 
			foundCount, target, maxRecipes)
	}

	return allResults, visited
}

func processRecipe(e1 string, e2 string, found map[string][]string, visitCount *int) bool {
	queue := []string{}

	if !utilities.IsBaseElement(e1) && found[e1] == nil {
		queue = append(queue, e1)
	}
	
	if !utilities.IsBaseElement(e2) && found[e2] == nil {
		queue = append(queue, e2)
	}

	if len(queue) == 0 {
		return true
	}
	
	for len(queue) > 0 {
		element := queue[0]
		queue = queue[1:]
		*visitCount++
		
		if found[element] != nil {
			continue
		}

		recipeList, exists := utilities.Recipes[element]
		if !exists {
			return false
		}
		
		elementProcessed := false
		for _, recipe := range recipeList {
			ing1 := recipe.Element1
			ing2 := recipe.Element2
			
			found[element] = []string{ing1, ing2}
			
			newElements := []string{}
			if !utilities.IsBaseElement(ing1) && found[ing1] == nil {
				newElements = append(newElements, ing1)
			}
			
			if !utilities.IsBaseElement(ing2) && found[ing2] == nil {
				newElements = append(newElements, ing2)
			}
			
			if len(newElements) == 0 {
				elementProcessed = true
				break
			}
			
			foundCopy := make(map[string][]string)
			for k, v := range found {
				foundCopy[k] = v
			}
			
			newQueue := append([]string{}, queue...)
			newQueue = append(newQueue, newElements...)
			
			allProcessed := true
			visitIncrease := 0
			
			for len(newQueue) > 0 {
				nextElement := newQueue[0]
				newQueue = newQueue[1:]
				visitIncrease++
				
				if foundCopy[nextElement] != nil {
					continue
				}
				
				nextRecipeList, exists := utilities.Recipes[nextElement]
				if !exists {
					allProcessed = false
					break
				}
				
				elementResolved := false
				for _, nextRecipe := range nextRecipeList {
					nextIng1 := nextRecipe.Element1
					nextIng2 := nextRecipe.Element2
					
					if (utilities.IsBaseElement(nextIng1) || foundCopy[nextIng1] != nil) &&
						(utilities.IsBaseElement(nextIng2) || foundCopy[nextIng2] != nil) {
						foundCopy[nextElement] = []string{nextIng1, nextIng2}
						elementResolved = true
						break
					}
					
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
				for k, v := range foundCopy {
					found[k] = v
				}
				elementProcessed = true
				break
			}
			
			delete(found, element)
		}
		
		if !elementProcessed {
			return false
		}
	}
	
	return true
}

func BFSFindRecipeAll(e1 string, e2 string, found map[string][]string, visitCount *int) bool {
	queue := []string{}
	
	if !utilities.IsBaseElement(e1) && found[e1] == nil {
		queue = append(queue, e1)
	}
	
	if !utilities.IsBaseElement(e2) && found[e2] == nil {
		queue = append(queue, e2)
	}
	
	if len(queue) == 0 {
		return true
	}
	
	for len(queue) > 0 {
		element := queue[0]
		queue = queue[1:]
		*visitCount++
		
		if found[element] != nil {
			continue
		}
		
		recipeList, exists := utilities.Recipes[element]
		if !exists {
			return false
		}
		
		elementProcessed := false
		for _, recipe := range recipeList {
			ing1 := recipe.Element1
			ing2 := recipe.Element2
			
			found[element] = []string{ing1, ing2}
		
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
			return false
		}
	}
	
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