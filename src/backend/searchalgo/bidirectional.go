package searchalgo

import (
	"fmt"
	"sync"
	"tubes2/utilities"
)

const MaxDepth = 40

func BiDirectionalSearch(target string, maxRecipes int) ([]utilities.RecipeTree, int) {
	counter := &SafeCounter{v: 0}

	if utilities.IsBaseElement(target) {
		tree := utilities.RecipeTree{Element: target}
		return []utilities.RecipeTree{tree}, 0
	}

	recipeList, exists := utilities.Recipes[target]
	if !exists || len(recipeList) == 0 {
		fmt.Printf("Target element '%s' doesn't exist or can't be created\n", target)
		return nil, 0
	}

	var allResults []utilities.RecipeTree
	var mutex sync.Mutex

	forwardQueue := []*utilities.Node{}
	backwardQueue := []*utilities.Node{}

	forwardVisitedMap := map[string]map[string][]string{}
	backwardVisitedMap := map[string]map[string][]string{}

	globalForwardVisited := map[string]bool{}
	globalBackwardVisited := map[string]bool{}

	startNode := &utilities.Node{
		Element:     target,
		Path:        []string{target},
		Visited:     map[string]bool{target: true},
		Depth:       0,
		Ingredients: map[string][]string{},
	}
	forwardQueue = append(forwardQueue, startNode)
	globalForwardVisited[target] = true

	for _, base := range utilities.BaseElements {
		node := &utilities.Node{
			Element:     base,
			Path:        []string{base},
			Visited:     map[string]bool{base: true},
			Depth:       0,
			Ingredients: map[string][]string{},
		}
		backwardQueue = append(backwardQueue, node)
		globalBackwardVisited[base] = true
		backwardVisitedMap[base] = map[string][]string{}
	}

	recipesFound := 0

	for len(forwardQueue) > 0 && len(backwardQueue) > 0 && (maxRecipes <= 0 || recipesFound < maxRecipes) {
		fq := forwardQueue
		forwardQueue = []*utilities.Node{}
		for _, node := range fq {
			counter.Inc()
			if node.Depth >= MaxDepth {
				continue
			}
			if recipes, ok := utilities.Recipes[node.Element]; ok {
				for _, recipe := range recipes {
					e1, e2 := recipe.Element1, recipe.Element2
					if forwardVisitedMap[node.Element] == nil {
						forwardVisitedMap[node.Element] = map[string][]string{}
					}
					forwardVisitedMap[node.Element][fmt.Sprintf("%s+%s", e1, e2)] = []string{e1, e2}
					utilities.TrackLiveUpdate(node.Element, node.Path, map[string][]string{node.Element: {e1, e2}})

					if backwardVisitedMap[e1] != nil || backwardVisitedMap[e2] != nil {
						complete := buildCompleteRecipe(target, []string{e1, e2}, e1, e2, forwardVisitedMap, backwardVisitedMap)
						tree := utilities.BuildRecipeTree(target, complete)

						mutex.Lock()
						allResults = append(allResults, tree)
						recipesFound++
						mutex.Unlock()

						fmt.Printf("Found recipe #%d for %s\n", recipesFound, target)
						if maxRecipes > 0 && recipesFound >= maxRecipes {
							return allResults, counter.Value()
						}
					}

					for _, elem := range []string{e1, e2} {
						if !globalForwardVisited[elem] && !utilities.IsBaseElement(elem) {
							globalForwardVisited[elem] = true
							forwardQueue = append(forwardQueue, &utilities.Node{
								Element:     elem,
								Path:        append(copySlice(node.Path), elem),
								Visited:     nil,
								Depth:       node.Depth + 1,
								Ingredients: copyIngredients(node.Ingredients),
							})
							fmt.Printf("Added to forwardQueue: %s, Depth: %d, Path: %v\n", elem, node.Depth+1, append(copySlice(node.Path), elem))

						}
					}
				}
			}
		}

		bq := backwardQueue
		backwardQueue = []*utilities.Node{}
		for _, node := range bq {
			counter.Inc()
			if node.Depth >= MaxDepth {
				continue
			}
			for result, recipes := range utilities.Recipes {
				for _, recipe := range recipes {
					if recipe.Element1 != node.Element && recipe.Element2 != node.Element {
						continue
					}
					other := recipe.Element2
					if recipe.Element2 == node.Element {
						other = recipe.Element1
					}

					if globalBackwardVisited[result] {
						continue
					}
					globalBackwardVisited[result] = true

					newIng := copyIngredients(node.Ingredients)
					newIng[result] = []string{node.Element, other}
					if backwardVisitedMap[result] == nil {
						backwardVisitedMap[result] = map[string][]string{}
					}
					backwardVisitedMap[result][fmt.Sprintf("%s+%s", node.Element, other)] = []string{node.Element, other}
					utilities.TrackLiveUpdate(result, append(copySlice(node.Path), result), map[string][]string{result: {node.Element, other}})

					fmt.Printf("globalForwardVisited: %v\n", globalForwardVisited)
					fmt.Printf("globalBackwardVisited: %v\n", globalBackwardVisited)
					if forwardVisitedMap[result] != nil {
						complete := buildCompleteRecipe(target, []string{recipe.Element1, recipe.Element2}, node.Element, other, forwardVisitedMap, backwardVisitedMap)
						tree := utilities.BuildRecipeTree(target, complete)

						mutex.Lock()
						allResults = append(allResults, tree)
						recipesFound++
						mutex.Unlock()

						fmt.Printf("Found recipe #%d for %s\n", recipesFound, target)
						if maxRecipes > 0 && recipesFound >= maxRecipes {
							return allResults, counter.Value()
						}
					}

					fmt.Printf("globalForwardVisited: %v\n", globalForwardVisited)
					fmt.Printf("globalBackwardVisited: %v\n", globalBackwardVisited)
					if !globalBackwardVisited[other] && !utilities.IsBaseElement(other) {
						backwardQueue = append(backwardQueue, &utilities.Node{
							Element:     other,
							Path:        append(copySlice(node.Path), result),
							Visited:     nil,
							Depth:       node.Depth + 1,
							Ingredients: newIng,
						})
						fmt.Printf("Added to backwardQueue: %s, Depth: %d, Path: %v\n", other, node.Depth+1, append(copySlice(node.Path), result))
					}
				}
			}
		}
	}

	fmt.Printf("Bidirectional search complete. Found %d recipes for %s\n", len(allResults), target)
	return allResults, counter.Value()
}

func buildCompleteRecipe(target string, forwardIngredients []string, backwardE1, backwardE2 string,
	forwardVisited, backwardVisited map[string]map[string][]string) map[string][]string {

	complete := map[string][]string{}
	complete[target] = forwardIngredients

	var expand func(string)
	expand = func(elem string) {
		if _, exists := complete[elem]; exists {
			return
		}

		if utilities.IsBaseElement(elem) {
			return
		}

		var ingredients []string
		if m, ok := backwardVisited[elem]; ok {
			for _, v := range m {
				ingredients = v
				break
			}
		} 
		if len(ingredients) == 0 {
			if m, ok := forwardVisited[elem]; ok {
				for _, v := range m {
					ingredients = v
					break
				}
			}
		}

		if len(ingredients) == 2 {
			complete[elem] = ingredients
			expand(ingredients[0])
			expand(ingredients[1])
		} else {
			fmt.Printf("Warning: Could not find ingredients for element %s\n", elem)
		}
	}

	for _, ing := range forwardIngredients {
		expand(ing)
	}

	return complete
}

func copySlice(src []string) []string {
	dst := make([]string, len(src))
	copy(dst, src)
	return dst
}

func copyIngredients(src map[string][]string) map[string][]string {
	dst := make(map[string][]string)
	for k, v := range src {
		dst[k] = append([]string{}, v...)
	}
	return dst
}
