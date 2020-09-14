package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/pdok/atom-generator/feeds"
	"gopkg.in/yaml.v2"
)

var file, output *string

func init() {
	file = flag.String("f", "", "yaml file containing the atom feed configuration")
	output = flag.String("o", ".", "the output directory where the files are written")
	flag.Parse()

	if len(*file) == 0 {
		log.Fatal("No configuration file found")
		return
	}

	if len(*output) == 0 {
		log.Fatal("No output directory found")
		return
	}
}

func main() {
	var doc []byte
	var config feeds.Feeds

	// Read and unmarshal config file
	doc, err := ioutil.ReadFile(*file)
	if err != nil {
		log.Fatalf("error: %v, with file: %v", err, `no file`)
	}
	if err := yaml.Unmarshal(doc, &config); err != nil {
		log.Fatalf("error: %v", err)
	}

	var processedFeeds []feeds.Feed

	// process Feeds
	for _, feed := range config.Feeds {
		processedFeeds = append(processedFeeds, feeds.ProcessFeed(feed))
	}

	// write both service and dataset feeds
	for _, feed := range processedFeeds {
		if !feed.Valid() {
			log.Fatalf(`ATOM Feeds with the id: %s is not valid`, feed.ID)
		}

		filename, err := feed.GetFileName()
		if err == nil {
			feed.WriteATOM(*output + `/` + filename)
		} else {
			log.Fatalf(`ATOM Feed NOT generated the id: %s`, feed.ID)
			return
		}
	}

	log.Println(`ATOM Feeds generated`)
}
