package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/beevik/etree"
	"github.com/topmedia/timevacuum/entities"
)

var (
	destURL     = "https://webservices7.autotask.net/ATServices/1.5/atws.asmx"
	soapAction  = "http://autotask.net/ATWS/v1_5/query"
	contentType = "text/xml"
)

type Client struct {
	HTTPClient *http.Client
	User       string
	Pass       string
	Response   *http.Response
}

func (c *Client) Request(q *entities.QueryXML) []byte {
	req, err := http.NewRequest("POST", destURL, q.ToReader())

	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.SetBasicAuth(c.User, c.Pass)
	req.Header.Add("SOAPAction", soapAction)
	req.Header.Add("Content-Type", contentType)

	c.Response, err = c.HTTPClient.Do(req)

	if code := c.Response.StatusCode; code != 200 {
		log.Fatalf("HTTP Request failed: %s", c.Response.Status)
	}

	if err != nil {
		log.Fatalf("Error reading response XML: %v", err)
	}

	return c.ExtractResults()

}

func (c *Client) ExtractResults() []byte {
	src := etree.NewDocument()
	src.ReadFrom(c.Response.Body)
	res := src.FindElement("//EntityResults")

	dst := etree.CreateDocument(res)
	b, err := dst.WriteToBytes()

	if err != nil {
		log.Fatalf("Error creating target document: %v", err)

	}
	return b
}

func (c *Client) Body() []byte {
	body, err := ioutil.ReadAll(c.Response.Body)

	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}
	return body
}

func (c *Client) FetchResources(qe *entities.QueryExpression) map[int]entities.Resource {
	q := qe.ToQueryXML()
	q.Entity("Resource")
	res := c.Request(q)
	var rr entities.ResourceResults
	xml.Unmarshal(res, &rr)

	m := make(map[int]entities.Resource, len(rr.Resources))
	for _, r := range rr.Resources {
		m[r.ID] = r
	}
	return m
}

func (c *Client) FetchTickets(qe *entities.QueryExpression) []entities.Ticket {
	q := qe.ToQueryXML()
	q.Entity("Ticket")
	res := c.Request(q)
	var tr entities.TicketResults
	xml.Unmarshal(res, &tr)
	return tr.Tickets
}

func (c *Client) FetchTicketByID(id int) *entities.Ticket {
	tr := c.FetchTickets(&entities.QueryExpression{Field: "id", Op: "equals", Value: fmt.Sprintf("%d", id)})

	if len(tr) == 0 {
		return nil
	}
	return &tr[0]
}

func (c *Client) FetchTimeEntries(qe *entities.QueryExpression) []entities.TimeEntry {
	q := qe.ToQueryXML()
	q.Entity("Timeentry")
	res := c.Request(q)
	var ter entities.TimeEntryResults
	xml.Unmarshal(res, &ter)

	rr := c.FetchResources(&entities.QueryExpression{
		Field: "Active", Op: "Equals", Value: "true"})

	for k, te := range ter.TimeEntries {
		r := rr[te.ResourceID]
		ter.TimeEntries[k].ResourceName = fmt.Sprintf("%s %s", r.FirstName, r.LastName)
	}

	return ter.TimeEntries
}

func NewClient(user, pass string) *Client {
	c := &Client{}
	c.User = user
	c.Pass = pass
	c.HTTPClient = &http.Client{}
	return c
}
