package searchalgo

import (
	"fmt"
	"strings"
	"time"
	"math/rand"
)

type Node struct {
    element   string
    path      []string
    recipes   []string
    available map[string]bool
    step      int
}
// DFSSingle menggunakan algoritma DFS untuk mencari satu jalur terpendek dari elemen dasar ke target
func DFSSingle(startElements []string, target string, recipes Recipe, tiers map[string]int) ([]string, int, bool) {

    var stack []Node
    visited := make(map[string]bool)
    totalVisited := 0

    initialAvailable := make(map[string]bool)
    for _, el := range startElements {
        initialAvailable[el] = true
    }

    // Inisialisasi stack dengan elemen-elemen awal
    for _, el := range startElements {
        stack = append(stack, Node{
            element:   el,
            path:      []string{el},
            recipes:   []string{},
            available: copyMaps(initialAvailable),
            step:      0,
        })
        visited[el] = true
    }

    startTime := time.Now()

    for len(stack) > 0 {
        // Pop dari stack (ambil elemen terakhir)
        current := stack[len(stack)-1]
        stack = stack[:len(stack)-1]

        totalVisited++

        // Cek apakah sudah mencapai target
        if current.element == target {
            duration := time.Since(startTime)
            fmt.Printf("\nPath found: %v\n", current.path)
            fmt.Printf("Recipes used: %v\n", current.recipes)
            fmt.Printf("Total nodes visited: %d\n", totalVisited)
            fmt.Printf("Processing time: %v\n", duration)
            return current.path, current.step, true
        }

        // Cari kombinasi elemen yang bisa dibuat
        var elementsToAdd []Node
        for product, combos := range recipes {
            if visited[product] {
                continue
            }

            // Pastikan tier elemen lebih kecil dari tier target
            if tiers[product] >= tiers[target] {
                continue
            }

            for _, combo := range combos {
                if current.available[combo[0]] && current.available[combo[1]] {
                    recipeStep := fmt.Sprintf("%s + %s => %s", combo[0], combo[1], product)
                    fmt.Printf("%s→ Combine %s + %s => %s\n", strings.Repeat("  ", current.step), combo[0], combo[1], product)

                    // Tambahkan elemen baru ke path
                    newPath := append([]string{}, current.path...)
                    newPath = append(newPath, product)

                    // Catat resep yang digunakan
                    newRecipes := append([]string{}, current.recipes...)
                    newRecipes = append(newRecipes, recipeStep)

                    // Update elemen yang tersedia
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

        // Acak urutan elemen yang akan ditambahkan ke stack
        shuffleNodes(elementsToAdd)

        // Tambahkan elemen-elemen baru ke stack
        for _, node := range elementsToAdd {
            stack = append(stack, node)
            visited[node.element] = true
        }
    }

    duration := time.Since(startTime)
    fmt.Printf("\nTarget not found. Total nodes visited: %d\n", totalVisited)
    fmt.Printf("Processing time: %v\n", duration)
    return nil, 0, false
}

// Fungsi untuk mengacak urutan node
func shuffleNodes(nodes []Node) {
    rand.Seed(time.Now().UnixNano())
    rand.Shuffle(len(nodes), func(i, j int) {
        nodes[i], nodes[j] = nodes[j], nodes[i]
    })
}

// DFSMultiple menggunakan algoritma DFS untuk menemukan beberapa jalur dari elemen dasar ke target
func DFSMultiple(startElements []string, target string, recipes Recipe, maxRecipes int) ([][]string, int) {
	type Node struct {
		element   string
		path      []string
		recipes   []string
		available map[string]bool
		step      int
	}

	var stack []Node
	visited := make(map[string]bool)
	var foundRecipes [][]string
	var foundRecipePaths [][]string
	totalVisited := 0

	initialAvailable := make(map[string]bool)
	for _, el := range startElements {
		initialAvailable[el] = true
	}

	// Inisialisasi stack dengan elemen-elemen awal
	for _, el := range startElements {
		stack = append(stack, Node{
			element:   el,
			path:      []string{el},
			recipes:   []string{},
			available: copyMaps(initialAvailable),
			step:      0,
		})
		visited[el] = true
	}

	startTime := time.Now()

	for len(stack) > 0 && len(foundRecipes) < maxRecipes {
		// Pop dari stack (ambil elemen terakhir)
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		totalVisited++

		// Cek apakah sudah mencapai target
		if current.element == target {
			foundRecipes = append(foundRecipes, current.path)
			foundRecipePaths = append(foundRecipePaths, current.recipes)
			
			// Jika sudah mencapai jumlah resep yang diinginkan, selesai
			if len(foundRecipes) >= maxRecipes {
				break
			}
			
			// Lanjutkan mencari resep lain
			continue
		}

		// Cari semua kombinasi elemen yang bisa dibuat
		elementsToAdd := []Node{}
		
		for product, combos := range recipes {
			if visited[product] {
				continue
			}

			for _, combo := range combos {
				if current.available[combo[0]] && current.available[combo[1]] {
					recipeStep := fmt.Sprintf("%s + %s => %s", combo[0], combo[1], product)
					fmt.Printf("%s→ Combine %s + %s => %s\n", strings.Repeat("  ", current.step), combo[0], combo[1], product)

					// Tambahkan elemen baru ke path
					newPath := append([]string{}, current.path...)
					newPath = append(newPath, product)

					// Catat resep yang digunakan
					newRecipes := append([]string{}, current.recipes...)
					newRecipes = append(newRecipes, recipeStep)

					// Update elemen yang tersedia
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
		
		// Untuk DFS, kita push elemen-elemen baru ke stack dalam urutan terbalik
		// agar elemen yang ditambahkan terakhir akan diproses terlebih dahulu
		for i := len(elementsToAdd) - 1; i >= 0; i-- {
			stack = append(stack, elementsToAdd[i])
			visited[elementsToAdd[i].element] = true
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

// Fungsi pembantu untuk menyalin map (sama seperti yang ada di implementasi BFS)
func copyMaps(original map[string]bool) map[string]bool {
	newMap := make(map[string]bool)
	for k, v := range original {
		newMap[k] = v
	}
	return newMap
}

func calculateTiers(recipes Recipe, startElements []string) map[string]int {
    tiers := make(map[string]int)
    queue := startElements

    // Elemen awal memiliki tier 0
    for _, el := range startElements {
        tiers[el] = 0
    }

    // Lakukan BFS untuk menghitung tier
    for len(queue) > 0 {
        current := queue[0]
        queue = queue[1:]

        for product, combos := range recipes {
            for _, combo := range combos {
                if combo[0] == current || combo[1] == current {
                    if _, exists := tiers[product]; !exists {
                        tiers[product] = tiers[current] + 1
                        queue = append(queue, product)
                    }
                }
            }
        }
    }

    return tiers
}