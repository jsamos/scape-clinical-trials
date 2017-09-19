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

type FeedJob struct {
	Counter int
	FeedItem *gofeed.Item
}

func fetchStudy(study *gofeed.Item) (ClinicalStudy, error) {
			itemGUID := study.GUID
			url := fmt.Sprintf("https://clinicaltrials.gov/ct2/show/%s?displayxml=true", itemGUID)
			resp, err := grequests.Get(url, nil)
			v := ClinicalStudy{}
			
			if err != nil {
				return v, err
			}

			content := resp.String()

		
			xml.Unmarshal([]byte(content), &v)

			return v, nil
}

func main() {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("https://clinicaltrials.gov/ct2/results/rss.xml?rcv_d=&lup_d=14&sel_rss=mod14&recrs=dghi&count=10000")

	jobs := make(chan FeedJob, len(feed.Items))
	done := make(chan bool)

  go func() {
      for {
          job, more := <-jobs
          
          if more {
						study, err := fetchStudy(job.FeedItem)

						if err != nil {
							log.Fatalln("Unable to make request: ", err)
						}

						fmt.Printf("Title: %q\n", study.Title)
						fmt.Printf("Url: %q\n", study.Url)
						fmt.Printf("Status: %q\n", study.Status)
						fmt.Printf("LeadSponsor: %q\n", study.LeadSponsor)
						fmt.Println("Request Complete", job.Counter)
          } else {
              fmt.Println("received all jobs")
              done <- true
              return
          }
      }
  }()

	//for i := 0; i < len(feed.Items); i++ {
	for i := 0; i < 20; i++ {
			jobs <- FeedJob{i, feed.Items[i]}
			
	}

	close(jobs)
	<-done
}