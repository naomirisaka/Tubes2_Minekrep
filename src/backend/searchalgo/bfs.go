package searchalgo

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type Result struct {
	Path    []string
	Recipes []string
}

type Recipe map[string][][]string

type MultiVisited map[string]map[int]bool

func BFSSingle(startElements []string, target string, recipes Recipe) ([]string, int, bool) {
	type Node struct {
		element   string
		path      []string
		recipes   []string
		available map[string]bool
		step      int
	}

	queue := []Node{}
	visited := MultiVisited{}
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
		markVisited(visited, el, 0)
	}

	startTime := time.Now()

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		totalVisited++

		if current.element == target {
			duration := time.Since(startTime)
			// fmt.Printf("\nShortest path found: %v\n", current.path)
			fmt.Printf("Recipes used: %v\n", current.recipes)
			fmt.Printf("Total nodes visited: %d\n", totalVisited)
			fmt.Printf("Processing time: %v\n", duration)
			return current.path, current.step, true
		}

		for product, combos := range recipes {
			for _, combo := range combos {
				if current.available[combo[0]] && current.available[combo[1]] &&
					!isVisited(visited, product, current.step+1) {

					recipeStep := fmt.Sprintf("%s + %s => %s", combo[0], combo[1], product)
					fmt.Printf("%s→ Combine %s + %s => %s\n", strings.Repeat("  ", current.step), combo[0], combo[1], product)

					newPath := append([]string{}, current.path...)
					newPath = append(newPath, product)

					newRecipes := append([]string{}, current.recipes...)
					newRecipes = append(newRecipes, recipeStep)

					newAvailable := copyMap(current.available)
					newAvailable[product] = true

					if product == target {
						duration := time.Since(startTime)
						// fmt.Printf("\nShortest path found: %v\n", newPath)
						fmt.Printf("Recipes used: %v\n", newRecipes)
						fmt.Printf("Total nodes visited: %d\n", totalVisited+1)
						fmt.Printf("Processing time: %v\n", duration)
						return newPath, current.step + 1, true
					}

					queue = append(queue, Node{
						element:   product,
						path:      newPath,
						recipes:   newRecipes,
						available: newAvailable,
						step:      current.step + 1,
					})
					markVisited(visited, product, current.step+1)
				}
			}
		}
	}

	duration := time.Since(startTime)
	fmt.Printf("\nTarget not found. Total nodes visited: %d\n", totalVisited)
	fmt.Printf("Processing time: %v\n", duration)
	return nil, 0, false
}


func BFSMultipleParallel(startElements []string, target string, recipes Recipe, maxRecipes int, maxWorkers int) ([][]string, int) {
	type Node struct {
		element   string
		path      []string
		recipes   []string
		available map[string]bool
		step      int
	}

	var mutex sync.Mutex
	var foundResults []Result
	totalVisited := 0
	var visitedMutex sync.Mutex
	visited := MultiVisited{}

	workQueue := make(chan Node, 1000)
	resultChan := make(chan Result, maxRecipes)
	done := make(chan struct{})
	var wg sync.WaitGroup // WaitGroup for worker goroutines
	
	initialAvailable := make(map[string]bool)
	for _, el := range startElements {
		initialAvailable[el] = true
	}

	for _, el := range startElements {
		workQueue <- Node{
			element:   el,
			path:      []string{el},
			recipes:   []string{},
			available: copyMap(initialAvailable),
			step:      0,
		}
		
		mutex.Lock()
		markVisited(visited, el, 0)
		mutex.Unlock()
	}

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			
			for {
				select {
				case <-done:
					return
				case current, ok := <-workQueue:
					if !ok {
						return
					}

					visitedMutex.Lock()
					totalVisited++
					visitedMutex.Unlock()

					if current.element == target {
						resultChan <- Result{
							Path:    current.path,
							Recipes: current.recipes,
						}
					}

					for product, combos := range recipes {
						for _, combo := range combos {
							if current.available[combo[0]] && current.available[combo[1]] {
								visitedMutex.Lock()
								isVisitedAlready := isVisited(visited, product, current.step+1)
								visitedMutex.Unlock()
								
								if !isVisitedAlready {
									recipeStep := fmt.Sprintf("%s + %s => %s", combo[0], combo[1], product)
									
									newPath := append([]string{}, current.path...)
									newPath = append(newPath, product)
									
									newRecipes := append([]string{}, current.recipes...)
									newRecipes = append(newRecipes, recipeStep)
									
									newAvailable := copyMap(current.available)
									newAvailable[product] = true
									
									visitedMutex.Lock()
									markVisited(visited, product, current.step+1)
									visitedMutex.Unlock()
									
									fmt.Printf("[Worker %d] %s→ Combine %s + %s => %s\n", 
										workerID,
										strings.Repeat("  ", current.step), 
										combo[0], combo[1], product)
									
									select {
									case workQueue <- Node{
										element:   product,
										path:      newPath,
										recipes:   newRecipes,
										available: newAvailable,
										step:      current.step + 1,
									}:
									case <-done:
										return
									}
								}
							}
						}
					}
				}
			}
		}(i)
	}
	
	// collect results in a separate goroutine
	go func() {
		for result := range resultChan {
			mutex.Lock()
			foundResults = append(foundResults, result)
			
			if len(foundResults) >= maxRecipes {
				close(done) 
			}
			mutex.Unlock()
		}
	}()

	wg.Wait()
	close(workQueue)
	close(resultChan)

	var foundRecipePaths [][]string
	for _, result := range foundResults {
		foundRecipePaths = append(foundRecipePaths, result.Recipes)
	}

	return foundRecipePaths, totalVisited
}

func BFSMultiple(startElements []string, target string, recipes Recipe, maxRecipes int, maxWorkers int, timeoutSeconds int) ([][]string, int, bool) {
	startTime := time.Now()
	resultChan := make(chan [][]string, 1)
	visitedChan := make(chan int, 1)
	
	go func() {
		results, visited := BFSMultipleParallel(startElements, target, recipes, maxRecipes, maxWorkers)
		resultChan <- results
		visitedChan <- visited
	}()
	
	select {
	case results := <-resultChan:
		visited := <-visitedChan
		duration := time.Since(startTime)
		fmt.Printf("\nFound %d recipes in %v\n", len(results), duration)
		return results, visited, true
	case <-time.After(time.Duration(timeoutSeconds) * time.Second):
		duration := time.Since(startTime)
		fmt.Printf("\nSearch timed out after %v\n", duration)
		return nil, 0, false
	}
}

func copyMap(original map[string]bool) map[string]bool {
	newMap := make(map[string]bool)
	for k, v := range original {
		newMap[k] = v
	}
	return newMap
}
func isVisited(v MultiVisited, el string, step int) bool {
	if v[el] == nil {
		return false
	}
	return v[el][step]
}

func markVisited(v MultiVisited, el string, step int) {
	if v[el] == nil {
		v[el] = make(map[int]bool)
	}
	v[el][step] = true
}
