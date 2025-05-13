package utilities

<<<<<<< HEAD
// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"time"
// )

=======
>>>>>>> dd6ca3248ae7b2d1452d4e3847b539e901bdfce9
type Recipe struct {
	Element1     string `json:"element1"`
	Element2     string `json:"element2"`
	Result       string `json:"result"`
	IconFilename string `json:"icon_filename"`
}

type Element struct {
	Name string
	Tier int
}

type Node struct {
<<<<<<< HEAD
	Element    string
	Path       []string
	Visited    map[string]bool
	Depth      int
	Ingredients map[string][]string 
}

type RecipeTree struct {
	Element    string      `json:"element"`
	Ingredients []RecipeTree `json:"ingredients,omitempty"`
}
=======
	Element     string
	Path        []string
	Visited     map[string]bool
	Depth       int
	Ingredients map[string][]string
}

type RecipeTree struct {
	Element     string       `json:"element"`
	Ingredients []RecipeTree `json:"ingredients,omitempty"`
}

>>>>>>> dd6ca3248ae7b2d1452d4e3847b539e901bdfce9
type Step struct {
	Current   string   `json:"current"`
	Queue     []string `json:"queue"`
	Element1  string   `json:"element1"`
	Element2  string   `json:"element2"`
	Result    string   `json:"result"`
}

type ResultStep struct {
	Element1     string `json:"element1"`
	Element2     string `json:"element2"`
	Result       string `json:"result"`
	IconFilename string `json:"icon_filename"`
}

type RecipeResult struct {
	TargetElement string       `json:"targetElement"`
	Steps         []ResultStep `json:"steps"`
}

type LiveUpdateStep struct {
	Step           int           `json:"step"`
	Message        string        `json:"message"`
	PartialTree    *RecipeResult `json:"partial_tree,omitempty"`
	HighlightNodes []string      `json:"highlight_nodes"`
}

var (
<<<<<<< HEAD
	Elements   = make(map[string]Element)
	Recipes    = make(map[string][]Recipe)
	BaseElements = []string{"Water", "Fire", "Earth", "Air"}
	Tiers      = make(map[string]int)
)
=======
	Elements     = make(map[string]Element)
	Recipes      = make(map[string][]Recipe)
	BaseElements = []string{"Water", "Fire", "Earth", "Air"}
	Tiers        = make(map[string]int)
)
>>>>>>> dd6ca3248ae7b2d1452d4e3847b539e901bdfce9
