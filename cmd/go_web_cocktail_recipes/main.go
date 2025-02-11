package main

import (
	"fmt"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)

}
