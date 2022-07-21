package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/pdok/atom-generator/feeds"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

const FILE string = `file`
const OUTPUT string = `output`

func main() {
	app := cli.NewApp()
	app.Name = "Atom Generator"
	app.Usage = "A Golang Atom generation application"

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     FILE,
			Aliases:  []string{"f"},
			Usage:    "Config file",
			Required: true,
			EnvVars:  []string{"FILE"},
		},
		&cli.StringFlag{
			Name:     OUTPUT,
			Aliases:  []string{"o"},
			Usage:    "Output directory",
			Required: true,
			EnvVars:  []string{"OUTPUT"},
		},
	}

	app.Action = func(c *cli.Context) error {

		// Read and unmarshal config file
		doc, err := ioutil.ReadFile(c.String(FILE))
		if err != nil {
			log.Fatalf("error: %v, with file: %v", err, `no file`)
		}
		var config feeds.Feeds
		if err := yaml.Unmarshal(doc, &config); err != nil {
			log.Fatalf("error: %v", err)
		}

		processedFeeds := feeds.ProcessFeeds(config)

		// write both service and dataset feeds
		for _, feed := range processedFeeds {
			if err := feed.Valid(); err != nil {
				log.Fatalf(`ATOM Feeds with the id: %s is not valid. With the error: %s`, feed.ID, err.Error())
			}

			filename, err := feed.GetFileName()
			if err == nil {
				feed.WriteATOM(c.String(OUTPUT) + `/` + filename)
			} else {
				log.Fatalf(`ATOM Feed NOT generated the id: %s`, feed.ID)
			}
		}

		log.Println(`ATOM Feeds generated`)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
