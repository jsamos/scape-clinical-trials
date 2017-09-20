package main

import (
	"fmt"
	"time"
	"log"
	"flag"
	"github.com/mmcdole/gofeed"
	"github.com/levigross/grequests"
	"runtime"
	"strings"
	"clinicaltrials/trialdate"
	"clinicaltrials/models"
)

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

func worker(jobs <-chan *gofeed.Item, output chan<- models.ClinicalStudy) {
	for job := range jobs {
		study, err := fetchStudy(job)

		if err != nil {
			log.Fatalln("Unable to make request: ", err)
		} else {
			output <- study
		}
	}
}

func main() {
	var fromDaysAgo int
	var processCount int
	var workersCount int
	flag.IntVar(&fromDaysAgo, "days", 5, "how many days back to allow")
	flag.IntVar(&processCount, "processes", 4, "how processes to use")
	flag.IntVar(&workersCount, "workers", 5, "how workers to use")
	flag.Parse()
	fmt.Println("days", workersCount)
	dateLimit := time.Now().AddDate(0, 0, -fromDaysAgo)
	trialDateStructure := "2006-01-02"
	excludeStudiesWith := []string{"Universi", "School", "College", "Hospital"}
	runtime.GOMAXPROCS(processCount)
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("https://clinicaltrials.gov/ct2/results/rss.xml?rcv_d=&lup_d=14&sel_rss=mod14&recrs=eghi&count=10000")
	jobs := make(chan *gofeed.Item, len(feed.Items))
	results := make(chan models.ClinicalStudy, len(feed.Items))

  for w := 1; w <= workersCount; w++ {
      go worker(jobs, results)
  }

  start := time.Now()

	for i := 0; i < len(feed.Items); i++ {
			jobs <- feed.Items[i]
			
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