package main

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/levigross/grequests"
	"log"
	"encoding/xml"
)

type ClinicalStudy struct {
	XMLName xml.Name `xml:"clinical_study"`
	Title    string   `xml:"official_title"`
	Status    string   `xml:"overall_status"`
	Url    string   `xml:"required_header>url"`
	LeadSponsor string   `xml:"sponsors>lead_sponsor>agency"`
	Collaborators []string   `xml:"sponsors>collaborator>agency"`
}

func main() {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("https://clinicaltrials.gov/ct2/results/rss.xml?rcv_d=&lup_d=14&sel_rss=mod14&recrs=dghi&count=10000")

	//for i := 0; i < len(feed.Items); i++ {
	for i := 0; i < 10; i++ {
			item := feed.Items[i]
			itemGUID := item.GUID
			url := fmt.Sprintf("https://clinicaltrials.gov/ct2/show/%s?displayxml=true", itemGUID)
			//fmt.Println(item.Title)
			//fmt.Println(url)
			resp, err := grequests.Get(url, nil)

			if err != nil {
				log.Fatalln("Unable to make request: ", err)
			}
			content := resp.String()
			v := ClinicalStudy{Title: "none"}
			xml.Unmarshal([]byte(content), &v)
			fmt.Printf("Title: %q\n", v.Title)
			fmt.Printf("Url: %q\n", v.Url)
			fmt.Printf("Status: %q\n", v.Status)
			fmt.Printf("LeadSponsor: %q\n", v.LeadSponsor)
			//fmt.Printf("Collaborators: %q\n", v.Collaborators)
	}
}