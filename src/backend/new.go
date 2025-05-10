func dfsSearch(target string, maxRecipes int) ([]RecipeTree, int) {
    visited = 0
    if isBaseElement(target) {
        tree := RecipeTree{Element: target}
        return []RecipeTree{tree}, visited
    }

    if _, exists := recipes[target]; !exists {
        fmt.Printf("Target element '%s' doesn't exist or can't be created\n", target)
        return nil, visited
    }

    // Array untuk menyimpan semua resep yang ditemukan
    var allResults []RecipeTree
    
    // Dapatkan semua resep langsung untuk target
    recipeList, _ := recipes[target]
    
    for _, recipe := range recipeList {
        if maxRecipes > 0 && len(allResults) >= maxRecipes {
            break
        }
        
        e1 := recipe.Element1
        e2 := recipe.Element2
        
        // Skip jika melanggar aturan tier
        if tiers[e1] >= tiers[target] || tiers[e2] >= tiers[target] {
            continue
        }
        
        // Cari resep untuk kombinasi ini
        found := make(map[string][]string)
        found[target] = []string{e1, e2} // Tambahkan resep target terlebih dahulu
        
        visitCount := 0
        findRecipeAll(e1, found, &visitCount)
        findRecipeAll(e2, found, &visitCount)
        visited += visitCount
        
        // Cek apakah semua elemen memiliki resep atau base elements
        valid := true
        for elem, ingredients := range found {
            if isBaseElement(elem) {
                continue
            }
            for _, ing := range ingredients {
                if !isBaseElement(ing) && found[ing] == nil {
                    valid = false
                    break
                }
            }
            if !valid {
                break
            }
        }
        
        if valid {
            recipeTree := buildRecipeTree(target, found)
            allResults = append(allResults, recipeTree)
        }
    }

    return allResults, visited
}

// Mencari semua resep mungkin untuk sebuah elemen
func findRecipeAll(element string, found map[string][]string, visitCount *int) {
    *visitCount++
    
    // Jika sudah base element atau sudah punya resep, kita selesai
    if isBaseElement(element) || found[element] != nil {
        return
    }
    
    // Dapatkan resep untuk elemen ini
    recipeList, exists := recipes[element]
    if !exists {
        return
    }
    
    // Coba setiap resep
    for _, recipe := range recipeList {
        e1 := recipe.Element1
        e2 := recipe.Element2
        
        // Skip jika melanggar aturan tier
        if tiers[e1] >= tiers[element] || tiers[e2] >= tiers[element] {
            continue
        }
        
        // Tambahkan resep ini ke map
        found[element] = []string{e1, e2}
        
        // Cari resep untuk komponen-komponen
        findRecipeAll(e1, found, visitCount)
        findRecipeAll(e2, found, visitCount)
        
        // Jika kita menemukan resep yang valid, kembali sekarang
        if (isBaseElement(e1) || found[e1] != nil) && 
           (isBaseElement(e2) || found[e2] != nil) {
            return
        }
        
        // Jika tidak valid, hapus dan coba resep berikutnya
        delete(found, element)
    }
}