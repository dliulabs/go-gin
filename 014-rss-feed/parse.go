package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

type Feed struct {
	/*
		XMLName  xml.Name `xml:"feed"`
		Text     string   `xml:",chardata"`
		Xmlns    string   `xml:"xmlns,attr"`
		Media    string   `xml:"media,attr"`
		Category struct {
			Text  string `xml:",chardata"`
			Term  string `xml:"term,attr"`
			Label string `xml:"label,attr"`
		} `xml:"category"`
		Updated string `xml:"updated"`
		Icon    string `xml:"icon"`
		ID      string `xml:"id"`
		Link    []struct {
			Text string `xml:",chardata"`
			Rel  string `xml:"rel,attr"`
			Href string `xml:"href,attr"`
			Type string `xml:"type,attr"`
		} `xml:"link"`
		Logo     string  `xml:"logo"`
		Subtitle string  `xml:"subtitle"`
		Title    string  `xml:"title"`
	*/
	Entry []Entry `xml:"entry"`
}

type Entry struct {
	/*
		Text   string `xml:",chardata"`
		Author struct {
			Text string `xml:",chardata"`
			Name string `xml:"name"`
			URI  string `xml:"uri"`
		} `xml:"author"`
		Category struct {
			Text  string `xml:",chardata"`
			Term  string `xml:"term,attr"`
			Label string `xml:"label,attr"`
		} `xml:"category"`
		Content struct {
			Text string `xml:",chardata"`
			Type string `xml:"type,attr"`
		} `xml:"content"`
		ID   string `xml:"id"`
	*/
	Link struct {
		Text string `xml:",chardata"`
		Href string `xml:"href,attr"`
	} `xml:"link"`
	//Updated   string `xml:"updated"`
	//Published string `xml:"published"`
	Title     string `xml:"title"`
	Thumbnail struct {
		Text string `xml:",chardata"`
		URL  string `xml:"url,attr"`
	} `xml:"thumbnail"`
}

func main() {
	data, err := ioutil.ReadFile("rss.xml")
	if err != nil {
		fmt.Printf("Failed to read file: %v\n", err.Error())
	}
	var feed Feed
	err = xml.Unmarshal([]byte(data), &feed)
	if err != nil {
		fmt.Printf("failed to load data: %v\n", err.Error())
	}

	for _, entry := range feed.Entry {
		fmt.Printf("Title: %v\n", entry.Title)
	}
}
