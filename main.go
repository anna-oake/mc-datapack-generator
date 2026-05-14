package main

import (
	"flag"
	"fmt"

	"github.com/anna-oake/mc-datapack-generator/internal/clientjar"
	"github.com/anna-oake/mc-datapack-generator/internal/datapack"
	"github.com/anna-oake/mc-datapack-generator/internal/manifest"
)

func main() {
	var previousVersion = flag.String("previous-version", "1.20.6", "specify the latest version to compare against")
	var datapackPath = flag.String("out", "./output/", "specify the path to save the datapack")
	flag.Parse()

	versions, err := manifest.GetVersionList()
	if err != nil {
		panic(err)
	}

	versions = manifest.FilterVersionList(versions, true, *previousVersion)
	if len(versions) == 0 {
		fmt.Print("version=none")
		return
	}
	v := versions[len(versions)-1]

	jarFile, err := v.FetchClientJar()
	if err != nil {
		panic(err)
	}
	jar, err := clientjar.Open(jarFile)
	if err != nil {
		panic(err)
	}
	recipes, err := jar.SmeltingRecipes()
	if err != nil {
		panic(err)
	}

	fullMap := makeIngredientsMap(recipes)
	dp := datapack.Datapack{
		Namespace: "smeltable",
		ItemTags:  fullMap,
	}
	err = dp.Save(*datapackPath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("version=%s", v.ID)
}
