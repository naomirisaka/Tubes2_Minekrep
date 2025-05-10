package searchalgo

import(
	""
)

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