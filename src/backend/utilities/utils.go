package utilities

import "fmt"

func IsBaseElement(element string) bool {
	for _, base := range BaseElements {
		if element == base {
			return true
		}
	}
	return false
}

func BuildRecipeTree(element string, ingredients map[string][]string) RecipeTree {
	tree := RecipeTree{Element: element}
	
	if ingList, exists := ingredients[element]; exists {
		for _, ing := range ingList {
			tree.Ingredients = append(tree.Ingredients, BuildRecipeTree(ing, ingredients))
		}
	}
	
	return tree
}

func CalculateTreeDepth(tree RecipeTree) int {
	if len(tree.Ingredients) == 0 {
		return 1
	}
	
	maxDepth := 0
	for _, ing := range tree.Ingredients {
		depth := CalculateTreeDepth(ing)
		if depth > maxDepth {
			maxDepth = depth
		}
	}
	
	return maxDepth + 1
}

func PrintRecipeTree(tree RecipeTree, indent string) {
	fmt.Printf("%s%s\n", indent, tree.Element)
	if len(tree.Ingredients) > 0 {
		fmt.Printf("%s└─ combines:\n", indent)
		for i, ing := range tree.Ingredients {
			if i == len(tree.Ingredients)-1 {
				PrintRecipeTree(ing, indent+"   ")
			} else {
				PrintRecipeTree(ing, indent+"│  ")
			}
		}
	}
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func IsSameRecipeTree(tree1, tree2 RecipeTree) bool {
    if tree1.Element != tree2.Element {
        return false
    }
    
    if len(tree1.Ingredients) != len(tree2.Ingredients) {
        return false
    }
    
    if len(tree1.Ingredients) == 2 && len(tree2.Ingredients) == 2 {
        normalOrder := IsSameRecipeTree(tree1.Ingredients[0], tree2.Ingredients[0]) &&
                       IsSameRecipeTree(tree1.Ingredients[1], tree2.Ingredients[1])
        
        reversedOrder := IsSameRecipeTree(tree1.Ingredients[0], tree2.Ingredients[1]) &&
                         IsSameRecipeTree(tree1.Ingredients[1], tree2.Ingredients[0])
        
        return normalOrder || reversedOrder
    }
    
    for i := range tree1.Ingredients {
        if !IsSameRecipeTree(tree1.Ingredients[i], tree2.Ingredients[i]) {
            return false
        }
    }
    
    return true
}

func CopyMap(original map[string][]string) map[string][]string {
    newMap := make(map[string][]string)
    for k, v := range original {
        newSlice := make([]string, len(v))
        copy(newSlice, v)
        newMap[k] = newSlice
    }
    return newMap
}