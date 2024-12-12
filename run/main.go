package main

import (
	"fmt"

	"github.com/custodia-cenv/hldc/src/hldcfs"
)

func main() {
	img, err := hldcfs.OpenHldcVfsImage("test.img")
	if err != nil {
		panic(err)
	}

	listResult, err := img.List()
	if err != nil {
		panic(err)
	}
	for _, item := range listResult {
		fmt.Println(item)
	}
	fmt.Println("Blöcke:", img.TotalBlocks())
	data, err := img.ReadData("tst")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}
