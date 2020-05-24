package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/shamansanchez/sc4go/dbpf"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("USAGE: %s <region_path>\n", os.Args[0])
		return
	}

	path := os.Args[1]

	paths, _ := ioutil.ReadDir(path)

	for _, c := range paths {
		n := c.Name()
		if filepath.Ext(n) == ".sc4" {
			city := dbpf.ReadDBPF(path + n)
			info := dbpf.GetRegionData(city)

			fmt.Printf("===== %s =====\nVersion: %d.%d\nMayor: %s\nPopulation: %d\nGUID: 0x%X\n\n",
				info.Name,
				info.MajorVersion,
				info.MinorVersion,
				info.MayorName,
				info.ResidentialPopulation+info.CommercialPopulation+info.IndustrialPopulation,
				info.GUID)
		}
	}

}
