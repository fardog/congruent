package main

import (
	"flag"
	"io/ioutil"

	"fmt"
	"github.com/fardog/congruent"
)

var files []string

func main() {
	flag.Parse()
	files := flag.Args()

	if len(files) < 1 {
		panic("Must provide at least one config file")
	}

	configs := make([]*congruent.Config, len(files))

	for i, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}

		config, err := congruent.NewConfigFromJSON(data)
		if err != nil {
			panic(err)
		}

		configs[i] = config
	}

	for _, c := range configs {
		servers := c.ResolvedServerConfigs()

		for _, r := range c.Requests {
			for _, s := range servers {
				req := congruent.NewRequest(s, r)
				body, err := req.Do()
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Printf("%s: %s\n", req.URI, body)
			}
		}
	}
}
