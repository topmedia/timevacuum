package entities

import "encoding/xml"

type TimeEntryResults struct {
	XMLName     xml.Name    `xml:"EntityResults"`
	TimeEntries []TimeEntry `xml:"Entity"`
}

type TimeEntry struct {
	XMLName      xml.Name `xml:"Entity"`
	ID           int      `xml:"id"`
	HoursWorked  float32
	ResourceID   int
	ResourceName string
}

type ResourceResults struct {
	XMLName   xml.Name   `xml:"EntityResults"`
	Resources []Resource `xml:"Entity"`
}

type Resource struct {
	XMLName    xml.Name `xml:"Entity"`
	ID         int      `xml:"id"`
	ResourceID int
	FirstName  string
	LastName   string
}
