package clientjar

import (
	"io"
	"strings"

	"github.com/pkg/errors"
	"github.com/valyala/fastjson"
)

type SmeltingRecipe struct {
	Type        string   `json:"type"`
	Ingredients []string `json:"ingredient"`
}

func parseIngredients(ing *fastjson.Value) ([]string, error) {
	switch ing.Type() {
	case fastjson.TypeString:
		return []string{string(ing.GetStringBytes())}, nil
	case fastjson.TypeObject:
		item := string(ing.GetStringBytes("item"))
		if item != "" {
			return []string{item}, nil
		}
		tag := string(ing.GetStringBytes("tag"))
		if tag != "" {
			return []string{"#" + tag}, nil
		}
		return nil, errors.New("unexpected ingredient format")
	case fastjson.TypeArray:
		var ingredients []string
		for _, i := range ing.GetArray() {
			el, err := parseIngredients(i)
			if err != nil {
				return nil, err
			}
			ingredients = append(ingredients, el...)
		}
		return ingredients, nil
	}

	return nil, errors.New("unexpected ingredient type " + ing.Type().String())
}

func (j *ClientJar) SmeltingRecipes() ([]SmeltingRecipe, error) {
	var recipes []SmeltingRecipe
	for _, file := range j.rd.File {
		if !strings.HasPrefix(file.Name, "data/minecraft/recipe/") {
			continue
		}
		if !strings.HasSuffix(file.Name, ".json") {
			continue
		}
		f, err := file.Open()
		if err != nil {
			return nil, errors.Wrap(err, "failed to open recipe file "+file.Name)
		}
		b, err := io.ReadAll(f)
		f.Close()
		if err != nil {
			return nil, errors.Wrap(err, "failed to read recipe file "+file.Name)
		}

		var p fastjson.Parser
		v, err := p.ParseBytes(b)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse recipe file "+file.Name)
		}

		t := strings.TrimPrefix(string(v.GetStringBytes("type")), "minecraft:")
		if t != "smelting" && t != "blasting" && t != "smoking" {
			continue
		}

		ing, err := parseIngredients(v.Get("ingredient"))
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse ingredients "+file.Name)
		}

		recipes = append(recipes, SmeltingRecipe{
			Type:        t,
			Ingredients: ing,
		})
		f.Close()
	}
	return recipes, nil
}
