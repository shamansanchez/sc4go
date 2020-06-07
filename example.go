package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/shamansanchez/sc4go/dbpf"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("USAGE: %s <region_path>\n", os.Args[0])
		return
	}

	path := os.Args[1]

	paths, err := ioutil.ReadDir(path)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(paths)

	for _, c := range paths {
		n := c.Name()
		if filepath.Ext(n) == ".sc4" {
			city, err := dbpf.ReadDBPF(fmt.Sprintf("%s%c%s", path, filepath.Separator, n))

			if err != nil {
				log.Fatal(err)
			}

			// for k, v := range city.Index {

			// 	if dir, ok := city.Directory[k]; ok {
			// 		log.Printf("%s: %d bytes, %d bytes uncompressed", k, v.FileSize, dir.FileSize)
			// 	} else {
			// 		log.Printf("%s: %d bytes", k, v.FileSize)

			// 	}
			// }

			info := dbpf.GetRegionData(city)

			log.Println(time.Unix(int64(city.Header.Created), 0))
			log.Println(time.Unix(int64(city.Header.Modified), 0))

			fmt.Printf("===== %s =====\nVersion: %d.%d\nMayor: %s\nPopulation: %d\nGUID: 0x%X\n\n",
				info.Name,
				info.MajorVersion,
				info.MinorVersion,
				info.MayorName,
				info.ResidentialPopulation+info.CommercialPopulation+info.IndustrialPopulation,
				info.GUID)

			// image := city.FileContents["8A2482B9.4A2482BB.0"]
			// ioutil.WriteFile(fmt.Sprintf("%s-%d.%d.png", info.Name, info.TileX, info.TileY), image, 755)
		}

	}

}
