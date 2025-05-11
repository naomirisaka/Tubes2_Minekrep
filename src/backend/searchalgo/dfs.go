package searchalgo

import (
	"fmt"
	"sync"
	"tubes2/backend/utilities"
)

// Atomic counter untuk visited
type SafeCounter struct {
	v   int
	mux sync.Mutex
}

func (c *SafeCounter) Inc() {
	c.mux.Lock()
	c.v++
	c.mux.Unlock()
}

func (c *SafeCounter) Value() int {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.v
}

// DFSSearch implements the depth-first search algorithm for recipe finding
func DFSSearch(target string, maxRecipes int) ([]utilities.RecipeTree, int) {
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
    
    // Try each recipe for the target element
    for i, recipe := range recipeList {
        if maxRecipes > 0 && len(allResults) >= maxRecipes {
            break
        }
        
        e1 := recipe.Element1
        e2 := recipe.Element2
        
        fmt.Printf("Exploring DFS recipe #%d: %s + %s => %s\n", i+1, e1, e2, target)
        
        // Initialize recipe map for tracking ingredients
        found := make(map[string][]string)
        found[target] = []string{e1, e2}
        
        // Track for live update visualization
        utilities.TrackLiveUpdate(target, []string{target}, found)
        
        counter.Inc()
        if processDFSRecipe(e1, e2, found, counter) {
            recipeTree := utilities.BuildRecipeTree(target, found)
            
            mutex.Lock()
            allResults = append(allResults, recipeTree)
            mutex.Unlock()
            
            fmt.Printf("  Found valid DFS recipe #%d for %s\n", len(allResults), target)
        }
    }
    
    fmt.Printf("DFS complete. Found %d recipes for %s\n", len(allResults), target)
    return allResults, counter.Value()
}

func processDFSRecipe(e1 string, e2 string, found map[string][]string, counter *SafeCounter) bool {
    // Process first element if it's not a base element and not processed yet
    if !utilities.IsBaseElement(e1) && found[e1] == nil {
        counter.Inc()
        
        recipeList, exists := utilities.Recipes[e1]
        if !exists || len(recipeList) == 0 {
            fmt.Printf("  Can't find recipes for '%s'\n", e1)
            return false
        }
        
        // Try the first recipe for this element
        recipe := recipeList[0]
        ing1 := recipe.Element1
        ing2 := recipe.Element2
        
        found[e1] = []string{ing1, ing2}
        
        // Track for live update visualization
        path := []string{e1}
        for k := range found {
            if k != e1 {
                path = append(path, k)
            }
        }
        currentIngredients := make(map[string][]string)
        currentIngredients[e1] = []string{ing1, ing2}
        utilities.TrackLiveUpdate(e1, path, currentIngredients)
        
        // Recursively process the ingredients
        if !processDFSRecipe(ing1, ing2, found, counter) {
            fmt.Printf("  Failed to process recipe for '%s'\n", e1)
            return false
        }
    }
    
    // Process second element if it's not a base element and not processed yet
    if !utilities.IsBaseElement(e2) && found[e2] == nil {
        counter.Inc()
        
        recipeList, exists := utilities.Recipes[e2]
        if !exists || len(recipeList) == 0 {
            fmt.Printf("  Can't find recipes for '%s'\n", e2)
            return false
        }
        
        // Try the first recipe for this element
        recipe := recipeList[0]
        ing1 := recipe.Element1
        ing2 := recipe.Element2
        
        found[e2] = []string{ing1, ing2}
        
        // Track for live update visualization
        path := []string{e2}
        for k := range found {
            if k != e2 {
                path = append(path, k)
            }
        }
        currentIngredients := make(map[string][]string)
        currentIngredients[e2] = []string{ing1, ing2}
        utilities.TrackLiveUpdate(e2, path, currentIngredients)
        
        // Recursively process the ingredients
        if !processDFSRecipe(ing1, ing2, found, counter) {
            fmt.Printf("  Failed to process recipe for '%s'\n", e2)
            return false
        }
    }
    
    return true
}

func ExploreAllCombinations(e1, e2 string, baseMap map[string][]string, results *[]map[string][]string, counter *SafeCounter) {
	counter.Inc()

	e1Maps := ExploreElementRecipes(e1, utilities.CopyMap(baseMap), counter)

	for _, map1 := range e1Maps {
		e2Maps := ExploreElementRecipes(e2, utilities.CopyMap(map1), counter)

		for _, completeMap := range e2Maps {
			*results = append(*results, completeMap)
		}
	}
}

func ExploreElementRecipes(element string, currentMap map[string][]string, counter *SafeCounter) []map[string][]string {

	if utilities.IsBaseElement(element) {
		return []map[string][]string{currentMap}
	}

	if _, ok := currentMap[element]; ok {
		return []map[string][]string{currentMap}
	}

	counter.Inc()

	recipeList, exists := utilities.Recipes[element]
	if !exists {
		return nil
	}

	var results []map[string][]string

	for _, recipe := range recipeList {
		e1 := recipe.Element1
		e2 := recipe.Element2

		if utilities.Tiers[e1] >= utilities.Tiers[element] || utilities.Tiers[e2] >= utilities.Tiers[element] {
			continue
		}

		newMap := utilities.CopyMap(currentMap)
		newMap[element] = []string{e1, e2}

		e1Maps := ExploreElementRecipes(e1, utilities.CopyMap(newMap), counter)

		for _, map1 := range e1Maps {
			e2Maps := ExploreElementRecipes(e2, utilities.CopyMap(map1), counter)

			results = append(results, e2Maps...)
		}
	}

	return results
}
