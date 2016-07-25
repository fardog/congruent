package main

import (
	"flag"
	"io/ioutil"

	"fmt"
	"github.com/fardog/congruent"
	"os"
	"text/tabwriter"
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

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, '\t', 0)

	for _, c := range configs {
		servers := c.ResolvedServerConfigs()

		for _, r := range c.Requests {

			for _, s := range servers {
				req := congruent.NewRequest(s, r)
				resp, err := req.Do()
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Fprintf(w, "URL: %s\tStatus: %d\n", req.URI, resp.StatusCode)
				fmt.Fprintf(w, "%s\n\n", resp.Body)
			}
		}
	}

	if err := w.Flush(); err != nil {
		panic(err)
	}
}
