package entities

import "encoding/xml"

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

type TicketResults struct {
	XMLName xml.Name `xml:"EntityResults"`
	Tickets []Ticket `xml:"Entity"`
}

type Ticket struct {
	XMLName            xml.Name `xml:"Entity"`
	ID                 int      `xml:"id"`
	AccountID          int
	AssignedResourceID int
	TicketNumber       string
	Title              string
	Description        string
	Status             int
}

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
	TicketID     int
	Ticket       *Ticket
}
