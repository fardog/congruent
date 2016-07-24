package main

import (
	"flag"
	"io/ioutil"

	"fmt"
	"github.com/fardog/congruent"
	"github.com/imdario/mergo"
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
		var servers congruent.ServerDefs

		for _, s := range c.Servers {
			gr := c.Global
			if err := mergo.MergeWithOverwrite(&gr, s); err != nil {
				panic(err)
			}

			servers = append(servers, gr)
		}

		for _, r := range c.Requests {
			for _, s := range servers {
				req := congruent.NewRequest(
					r.Method, s.BaseURI+r.Path, s.Headers,
				)
				body, err := req.Do()
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Printf("%s: %s\n", req.URI, body)
			}
		}
	}
}
