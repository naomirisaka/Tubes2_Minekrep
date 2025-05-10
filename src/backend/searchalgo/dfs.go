package searchalgo

import (
    "fmt"
    "strings"
    "math/rand"
    "time"
)

type Recipes map[string][][]string

type Node struct {
    element   string
    path      []string
    recipes   []string
    available map[string]bool
    step      int
}
func DFSSingle(startElements []string, target string, recipes Recipes, tiers map[string]int) ([]string, int, bool) {
    var stack []Node
    visited := make(map[string]bool)
    totalVisited := 0

    initialAvailable := make(map[string]bool)
    for _, el := range startElements {
        initialAvailable[el] = true
    }

    for _, el := range startElements {
        stack = append(stack, Node{
            element:   el,
            path:      []string{el},
            recipes:   []string{},
            available: copyMaps(initialAvailable),
            step:      0,
        })
    }

    startTime := time.Now()

    for len(stack) > 0 {
        current := stack[len(stack)-1]
        stack = stack[:len(stack)-1]

        if visited[current.element] {
            continue
        }
        visited[current.element] = true

        totalVisited++
        fmt.Printf("Visiting: %s, Path: %v\n", current.element, current.path)

        // Cek apakah sudah mencapai target
        if current.element == target {
            duration := time.Since(startTime)
            fmt.Printf("\nPath found: %v\n", current.path)
            fmt.Printf("Recipes used: %v\n", current.recipes)
            fmt.Printf("Total nodes visited: %d\n", totalVisited)
            fmt.Printf("Processing time: %v\n", duration)
            return current.path, current.step, true
        }

        var elementsToAdd []Node
        for product, combos := range recipes {
            if visited[product] {
                continue
            }

            if tiers[product] >= tiers[target] {
                continue
            }

            for _, combo := range combos {
                if current.available[combo[0]] && current.available[combo[1]] {
                    recipeStep := fmt.Sprintf("%s + %s => %s", combo[0], combo[1], product)
                    fmt.Printf("%s→ Combine %s + %s => %s\n", strings.Repeat("  ", current.step), combo[0], combo[1], product)

                    newPath := append([]string{}, current.path...)
                    newPath = append(newPath, product)

                    newRecipes := append([]string{}, current.recipes...)
                    newRecipes = append(newRecipes, recipeStep)

                    newAvailable := copyMaps(current.available)
                    newAvailable[product] = true

                    elementsToAdd = append(elementsToAdd, Node{
                        element:   product,
                        path:      newPath,
                        recipes:   newRecipes,
                        available: newAvailable,
                        step:      current.step + 1,
                    })
                }
            }
        }

        for _, node := range elementsToAdd {
            fmt.Printf("Adding to stack: %s, Path: %v\n", node.element, node.path)
        }

        shuffleNodes(elementsToAdd)

        stack = append(stack, elementsToAdd...)
    }

    duration := time.Since(startTime)
    fmt.Printf("\nTarget not found. Total nodes visited: %d\n", totalVisited)
    fmt.Printf("Processing time: %v\n", duration)
    return nil, 0, false
}

func CalculateTiers(recipes Recipes, startElements []string) map[string]int {
    tiers := make(map[string]int)
    queue := startElements

    for _, el := range startElements {
        tiers[el] = 0
    }

    for len(queue) > 0 {
        current := queue[0]
        queue = queue[1:]

        for product, combos := range recipes {
            for _, combo := range combos {
                if (combo[0] == current || combo[1] == current) && tiers[product] == 0 {
                    tiers[product] = tiers[current] + 1
                    queue = append(queue, product)
                }
            }
        }
    }

    return tiers
}

func copyMaps(original map[string]bool) map[string]bool {
	newMap := make(map[string]bool)
	for k, v := range original {
		newMap[k] = v
	}
	return newMap
}

func shuffleNodes(nodes []Node) {
    rand.Seed(time.Now().UnixNano())
    rand.Shuffle(len(nodes), func(i, j int) {
        nodes[i], nodes[j] = nodes[j], nodes[i]
    })
}