package searchalgo

import(
	"fmt"
    "tubes2/utilities"
)

var visited int

func DFSSearch(target string, maxRecipes int) ([]utilities.RecipeTree, int) {
    visited = 0
    if utilities.IsBaseElement(target) {
        tree := utilities.RecipeTree{Element: target}
        return []utilities.RecipeTree{tree}, visited
    }

    if _, exists := utilities.Recipes[target]; !exists {
        fmt.Printf("Target element '%s' doesn't exist or can't be created\n", target)
        return nil, visited
    }

    var allResults []utilities.RecipeTree
    
    recipeList, _ := utilities.Recipes[target]
    
    for _, recipe := range recipeList {
        if maxRecipes > 0 && len(allResults) >= maxRecipes {
            break
        }
        
        e1 := recipe.Element1
        e2 := recipe.Element2
        
        // Skip jika lower tier
        e1Tier, e1Exists := utilities.Tiers[e1]
        e2Tier, e2Exists := utilities.Tiers[e2]
        targetTier, targetExists := utilities.Tiers[target]
        
        if e1Exists && e2Exists && targetExists && 
           (e1Tier >= targetTier || e2Tier >= targetTier) {
            continue
        }
        
        found := make(map[string][]string)
        found[target] = []string{e1, e2} 
        
        visitCount := 0
        FindRecipeAll(e1, found, &visitCount)
        FindRecipeAll(e2, found, &visitCount)
        visited += visitCount
        
        valid := true
        for elem, ingredients := range found {
            if utilities.IsBaseElement(elem) {
                continue
            }
            for _, ing := range ingredients {
                if !utilities.IsBaseElement(ing) && found[ing] == nil {
                    valid = false
                    break
                }
            }
            if !valid {
                break
            }
        }
        
        if valid {
            recipeTree := utilities.BuildRecipeTree(target, found)
            allResults = append(allResults, recipeTree)
        }
    }

    return allResults, visited
}

func FindRecipeAll(element string, found map[string][]string, visitCount *int) {
    *visitCount++
    
    if utilities.IsBaseElement(element) || found[element] != nil {
        return
    }
    
    recipeList, exists := utilities.Recipes[element]
    if !exists {
        return
    }
    
    for _, recipe := range recipeList {
        e1 := recipe.Element1
        e2 := recipe.Element2
        
        if utilities.Tiers[e1] >= utilities.Tiers[element] || utilities.Tiers[e2] >= utilities.Tiers[element] {
            continue
        }
        
        found[element] = []string{e1, e2}
        
        // Cari resep untuk komponen-komponen
        FindRecipeAll(e1, found, visitCount)
        FindRecipeAll(e2, found, visitCount)
        
        if (utilities.IsBaseElement(e1) || found[e1] != nil) && 
           (utilities.IsBaseElement(e2) || found[e2] != nil) {
            return
        }
        
        delete(found, element)
    }
}