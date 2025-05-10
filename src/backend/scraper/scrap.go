package scraper

import (
	"encoding/json"
	//"fmt"
	//"log"
	//"net/http"
	//"os"
	//"strings"
	"io/ioutil"
	//"sync"
	"tubes2_minekrep/src/backend/utilities"
	// "github.com/PuerkitoBio/goquery"
)

type Recipe struct {
	Element1     string `json:"element1"`
	Element2     string `json:"element2"`
	Result       string `json:"result"`
	IconFilename string `json:"icon_filename"`
}

// func main() {
// 	baseURL := "https://little-alchemy.fandom.com"
// 	elementPageURL := baseURL + "/wiki/Elements_(Little_Alchemy_2)"

// 	os.MkdirAll("data", os.ModePerm)

// 	elementLinks := getElementLinks(elementPageURL, baseURL)

// 	var allRecipes []Recipe
// 	var mutex sync.Mutex
// 	var wg sync.WaitGroup

// 	semaphore := make(chan struct{}, 10) // 10 goroutine bersamaan

// 	for _, link := range elementLinks {
// 		wg.Add(1)
// 		semaphore <- struct{}{} // Acquire semaphore

// 		go func(url string) {
// 			defer wg.Done()
// 			defer func() { <-semaphore }() // Release semaphore

// 			recipes := scrapeElementPage(url)

// 			mutex.Lock()
// 			allRecipes = append(allRecipes, recipes...)
// 			mutex.Unlock()
// 		}(link)
// 	}

// 	wg.Wait() // Tunggu semua goroutine selesai

// 	err1 := saveToJSON(allRecipes, "data/recipes.json")
// 	// err2 := saveToCSV(allRecipes, "data/recipes.csv")

// 	if err1 != nil { // || err2 != nil {
// 		log.Println("Terjadi kesalahan saat menyimpan data.")
// 	} else {
// 		fmt.Println("Selesai! Data disimpan ke folder 'data/'.")
// 	}
// }

// func getElementLinks(url, baseURL string) []string {
// 	client := &http.Client{}
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		log.Fatal("Error creating request:", err)
// 	}

// 	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

// 	res, err := client.Do(req)
// 	if err != nil {
// 		log.Fatal("Gagal akses halaman utama:", err)
// 	}
// 	defer res.Body.Close()

// 	if res.StatusCode != 200 {
// 		log.Fatalf("Status code error: %d %s", res.StatusCode, res.Status)
// 	}

// 	doc, err := goquery.NewDocumentFromReader(res.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	links := make(map[string]bool)
// 	doc.Find("a").Each(func(_ int, s *goquery.Selection) {
// 		href, exists := s.Attr("href")
// 		if exists && strings.HasPrefix(href, "/wiki/") &&
// 			!strings.Contains(href, ":") && !strings.Contains(href, "#") {
// 			fullURL := baseURL + href
// 			links[fullURL] = true
// 		}
// 	})

// 	var uniqueLinks []string
// 	for link := range links {
// 		uniqueLinks = append(uniqueLinks, link)
// 	}

// 	return uniqueLinks
// }

// func scrapeElementPage(url string) []Recipe {
// 	var recipes []Recipe

// 	client := &http.Client{}
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		return recipes
// 	}

// 	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

// 	res, err := client.Do(req)
// 	if err != nil {
// 		return recipes
// 	}
// 	defer res.Body.Close()

// 	if res.StatusCode != 200 {
// 		return recipes
// 	}

// 	doc, err := goquery.NewDocumentFromReader(res.Body)
// 	if err != nil {
// 		return recipes
// 	}

// 	pageTitle := doc.Find("h1.page-header__title").Text()
// 	result := strings.TrimSpace(pageTitle)
// 	iconFilename := strings.ToLower(strings.ReplaceAll(result, " ", "_")) + ".png"

// 	found := false
// 	doc.Find(".mw-parser-output").Children().Each(func(i int, s *goquery.Selection) {
// 		if goquery.NodeName(s) == "h2" && strings.Contains(s.Text(), "Recipes") {
// 			found = true
// 			return
// 		}

// 		if found && goquery.NodeName(s) == "ul" {
// 			s.Find("li").Each(func(_ int, li *goquery.Selection) {
// 				text := li.Text()
// 				if strings.Contains(text, "+") {
// 					parts := strings.Split(text, "+")
// 					if len(parts) == 2 {
// 						element1 := strings.TrimSpace(parts[0])
// 						element2 := strings.TrimSpace(parts[1])
// 						recipes = append(recipes, Recipe{
// 							Element1:     element1,
// 							Element2:     element2,
// 							Result:       result,
// 							IconFilename: iconFilename,
// 						})
// 					}
// 				}
// 			})
// 			found = false
// 		}
// 	})

// 	return recipes
// }

// func saveToJSON(data []Recipe, filename string) error {
// 	file, err := os.Create(filename)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	encoder := json.NewEncoder(file)
// 	encoder.SetIndent("", "  ")
// 	return encoder.Encode(data)
// }
func initializeTiers() {
	// Set base tier 1
	for _, element := range utilities.BaseElements {
		utilities.Tiers[element] = 1
	}
	queue := make([]string, 0)
    queue = append(queue, utilities.BaseElements...)
    processed := make(map[string]bool)
    
    for _, elem := range utilities.BaseElements {
        processed[elem] = true
    }
	for len(queue) > 0 {
        current := queue[0]
        queue = queue[1:]
        
        for result, recipeList := range utilities.Recipes {
            if processed[result] {
                continue 
            }
            
            for _, recipe := range recipeList {
                if (recipe.Element1 == current || recipe.Element2 == current) {
                    // Only compute tier if both ingredients have tiers
                    if tier1, ok1 := utilities.Tiers[recipe.Element1]; ok1 {
                        if tier2, ok2 := utilities.Tiers[recipe.Element2]; ok2 {
                            resultTier := utilities.Max(tier1, tier2) + 1
                            existingTier, exists := utilities.Tiers[result]
                            
                            // Update tier kalau ada yagn lebih pendek
                            if !exists || resultTier < existingTier {
                                utilities.Tiers[result] = resultTier
                                if !processed[result] {
                                    queue = append(queue, result)
                                }
                            }
                            
                            processed[result] = true
                        }
                    }
                }
            }
        }
    }
}

func LoadRecipes(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	var recipeList []utilities.Recipe
	err = json.Unmarshal(data, &recipeList)
	if err != nil {
		return err
	}

	for _, recipe := range recipeList {
		utilities.Recipes[recipe.Result] = append(utilities.Recipes[recipe.Result], recipe)
	}

	initializeTiers()

	return nil
}
// func saveToCSV(data []Recipe, filename string) error {
// 	file, err := os.Create(filename)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	writer := csv.NewWriter(file)
// 	defer writer.Flush()

// 	writer.Write([]string{"Element1", "Element2", "Result", "IconFilename"})

// 	for _, r := range data {
// 		writer.Write([]string{r.Element1, r.Element2, r.Result, r.IconFilename})
// 	}

// 	return nil
// }
