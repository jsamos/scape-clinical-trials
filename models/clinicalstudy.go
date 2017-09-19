package models

import (
	"encoding/xml"
)

type ClinicalStudy struct {
	XMLName xml.Name `xml:"clinical_study"`
	Title    string   `xml:"official_title"`
	Status    string   `xml:"overall_status"`
	Url    string   `xml:"required_header>url"`
	LeadSponsor string   `xml:"sponsors>lead_sponsor>agency"`
	Collaborators []string   `xml:"sponsors>collaborator>agency"`
	DateUpdated string   `xml:"lastchanged_date"`
}

func BuildClinicalStudyFromXml(content string) ClinicalStudy {
	v := ClinicalStudy{}
	xml.Unmarshal([]byte(content), &v)
	return v
}