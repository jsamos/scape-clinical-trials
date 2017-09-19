package main

import (
	"fmt"
	"time"
	"log"
	"github.com/mmcdole/gofeed"
	"github.com/levigross/grequests"
	"runtime"
	"strings"
	"clinicaltrials/trialdate"
	"clinicaltrials/models"
)

type FeedJob struct {
	Counter int
	FeedItem *gofeed.Item
}

func fetchStudy(study *gofeed.Item) (models.ClinicalStudy, error) {
			itemGUID := study.GUID
			url := fmt.Sprintf("https://clinicaltrials.gov/ct2/show/%s?displayxml=true", itemGUID)
			resp, err := grequests.Get(url, nil)
			
			if err != nil {
				return models.ClinicalStudy{}, err
			}

			model := models.BuildClinicalStudyFromXml(resp.String())
			return model, nil
}

func worker(jobs <-chan FeedJob, output chan<- models.ClinicalStudy) {
	for job := range jobs {
		study, err := fetchStudy(job.FeedItem)

		if err != nil {
			log.Fatalln("Unable to make request: ", err)
		} else {
			output <- study
		}
	}
}

func main() {
	dateLimit := time.Now().AddDate(0, 0, -5)
	trialDateStructure := "2006-01-02"
	excludeStudiesWith := []string{"Universi", "School", "College", "Hospital"}
	runtime.GOMAXPROCS(4)
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("https://clinicaltrials.gov/ct2/results/rss.xml?rcv_d=&lup_d=14&sel_rss=mod14&recrs=eghi&count=10000")

	jobs := make(chan FeedJob, len(feed.Items))
	results := make(chan models.ClinicalStudy, len(feed.Items))

  for w := 1; w <= 5; w++ {
      go worker(jobs, results)
  }

  start := time.Now()

	//for i := 0; i < len(feed.Items); i++ {
	for i := 0; i < len(feed.Items); i++ {
			jobs <- FeedJob{i, feed.Items[i]}
			
	}

	close(jobs)

  for a := 1; a <= len(feed.Items); a++ {
    study := <-results

  	isPublicCompany := true

    for word := range excludeStudiesWith {
  		if strings.Contains(study.LeadSponsor, excludeStudiesWith[word]) {
  			isPublicCompany = false
  			break
  		}
  	}

  	dateFormat := trialdate.Formatter(study.DateUpdated)
  	trialFormattedDate, err := dateFormat()

  	if err != nil {
  		continue
  	}

  	updatedAt, err := time.Parse(trialDateStructure, trialFormattedDate)

  	if err != nil {
  		fmt.Println(err)
  	}

    if isPublicCompany == true && updatedAt.After(dateLimit) {
			fmt.Printf("Title: %q\n", study.Title)
			fmt.Printf("Url: %q\n", study.Url)
			fmt.Printf("Status: %q\n", study.Status)
			fmt.Printf("DateUpdated: %q\n", study.DateUpdated)
			fmt.Printf("LeadSponsor: %q\n", study.LeadSponsor)
			fmt.Println("")
			fmt.Println("")
    }
  }

  end := time.Now()
  elapsed := end.Sub(start)

  fmt.Println("Time Elapsed:", elapsed)
}