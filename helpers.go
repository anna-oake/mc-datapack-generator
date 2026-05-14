package main

import (
	"sort"

	"github.com/anna-oake/mc-datapack-generator/internal/clientjar"
)

func makeIngredientsMap(recipes []clientjar.SmeltingRecipe) map[string][]string {
	smelting := make(map[string]bool)
	blasting := make(map[string]bool)
	smoking := make(map[string]bool)

	for _, recipe := range recipes {
		for _, ingredient := range recipe.Ingredients {
			switch recipe.Type {
			case "smelting":
				smelting[ingredient] = true
			case "blasting":
				blasting[ingredient] = true
			case "smoking":
				smoking[ingredient] = true
			}
		}
	}

	var smeltingKeys []string
	for i := range smelting {
		smeltingKeys = append(smeltingKeys, i)
	}
	sort.Strings(smeltingKeys)

	var blastingKeys []string
	for i := range blasting {
		blastingKeys = append(blastingKeys, i)
	}
	sort.Strings(blastingKeys)

	var smokingKeys []string
	for i := range smoking {
		smokingKeys = append(smokingKeys, i)
	}
	sort.Strings(smokingKeys)

	return map[string][]string{
		"via_smelting": smeltingKeys,
		"via_blasting": blastingKeys,
		"via_smoking":  smokingKeys,
	}
}
