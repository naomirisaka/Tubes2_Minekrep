package searchalgo

import (
	"fmt"
	"strings"
	"time"
)

type Recipe map[string][][]string

// masih ada json yg blm properly read (max visited nodes baru sampe 80)??? masih ada yg not found, ex: Wine, Vinegar
func BFSSingle(startElements []string, target string, recipes Recipe) ([]string, int, bool) {
	type Node struct {
		element   string
		path      []string
		recipes   []string
		available map[string]bool
		step      int
	}

	queue := []Node{}
	visited := make(map[string]bool)
	totalVisited := 0

	initialAvailable := make(map[string]bool)
	for _, el := range startElements {
		initialAvailable[el] = true
	}

	for _, el := range startElements {
		// initialize queue
		queue = append(queue, Node{
			element:   el,
			path:      []string{el},
			recipes:   []string{},
			available: copyMap(initialAvailable),
			step:      0,
		})
		visited[el] = true
	}

	startTime := time.Now()

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		totalVisited++

		if current.element == target {
			duration := time.Since(startTime)
			fmt.Printf("\nShortest path found: %v\n", current.path)
			fmt.Printf("Recipes used: %v\n", current.recipes)
			fmt.Printf("Total nodes visited: %d\n", totalVisited)
			fmt.Printf("Processing time: %v\n", duration)
			return current.path, current.step, true
		}

		for product, combos := range recipes {
			for _, combo := range combos {
				if current.available[combo[0]] && current.available[combo[1]] && !visited[product] {
					recipeStep := fmt.Sprintf("%s + %s => %s", combo[0], combo[1], product)
					fmt.Printf("%s→ Combine %s + %s => %s\n", strings.Repeat("  ", current.step), combo[0], combo[1], product)

					// add recipe to the path
					newPath := append([]string{}, current.path...)
					newPath = append(newPath, product)

					// record recipe used
					newRecipes := append([]string{}, current.recipes...)
					newRecipes = append(newRecipes, recipeStep)

					newAvailable := copyMap(current.available)
					newAvailable[product] = true

					queue = append(queue, Node{
						element:   product,
						path:      newPath,
						recipes:   newRecipes,
						available: newAvailable,
						step:      current.step + 1,
					})
					visited[product] = true
				}
			}
		}
	}

	duration := time.Since(startTime)
	fmt.Printf("\nTarget not found. Total nodes visited: %d\n", totalVisited)
	fmt.Printf("Processing time: %v\n", duration)
	return nil, 0, false
}

// masih baca 1 json itu 1 recipe, pdhl kalo salah satu elemennya bisa dibentuk dr banyak cara = banyak recipe
func BFSMultiple(startElements []string, target string, recipes Recipe, maxRecipes int) ([][]string, int) {
	type Node struct {
		element   string
		path      []string
		recipes   []string
		available map[string]bool
		step      int
	}
	visited := make(map[string]bool)
	queue := []Node{}
	var foundRecipes [][]string
	var foundRecipePaths [][]string
	totalVisited := 0

	initialAvailable := make(map[string]bool)
	for _, el := range startElements {
		initialAvailable[el] = true
	}

	for _, el := range startElements {
		queue = append(queue, Node{
			element:   el,
			path:      []string{el},
			recipes:   []string{},
			available: copyMap(initialAvailable),
			step:      0,
		})
		visited[el] = true
	}

	startTime := time.Now()

	steps := 0
	for len(queue) > 0 && len(foundRecipes) < maxRecipes {
		current := queue[0]
		queue = queue[1:]
		steps++

		totalVisited++

		if current.element == target {
			foundRecipes = append(foundRecipes, current.path)
			foundRecipePaths = append(foundRecipePaths, current.recipes)
		}

		for product, combos := range recipes {
			for _, combo := range combos {
				if current.available[combo[0]] && current.available[combo[1]] && !visited[product] {
					recipeStep := fmt.Sprintf("%s + %s => %s", combo[0], combo[1], product)
					fmt.Printf("%s→ Combine %s + %s => %s\n", strings.Repeat("  ", current.step), combo[0], combo[1], product)

					newPath := append([]string{}, current.path...)
					newPath = append(newPath, product)

					newRecipes := append([]string{}, current.recipes...)
					newRecipes = append(newRecipes, recipeStep)

					newAvailable := copyMap(current.available)
					newAvailable[product] = true

					queue = append(queue, Node{
						element:   product,
						path:      newPath,
						recipes:   newRecipes,
						available: newAvailable,
						step:      current.step + 1,
					})
					visited[product] = true
				}
			}
		}
	}

	duration := time.Since(startTime)
	fmt.Printf("\nTotal recipes found: %d\n", len(foundRecipes))
	fmt.Printf("Total nodes visited: %d\n", totalVisited)
	fmt.Printf("Processing time: %v\n", duration)

	if len(foundRecipes) < maxRecipes {
		fmt.Println("Not enough recipes found. Returning the available ones.")
	}

	return foundRecipePaths, totalVisited
}

func copyMap(original map[string]bool) map[string]bool {
	newMap := make(map[string]bool)
	for k, v := range original {
		newMap[k] = v
	}
	return newMap
}
