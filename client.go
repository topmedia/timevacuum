package main

import (
	"encoding/xml"
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

func (c *Client) FetchResources(qe *entities.QueryExpression) []entities.Resource {
	q := qe.ToQueryXML()
	q.Entity("resource")
	res := c.Request(q)
	var rr entities.ResourceResults
	xml.Unmarshal(res, &rr)
	return rr.Resources
}

func (c *Client) FetchTimeEntries(qe *entities.QueryExpression) []entities.TimeEntry {
	q := qe.ToQueryXML()
	q.Entity("timeentry")
	res := c.Request(q)
	var te entities.TimeEntryResults
	xml.Unmarshal(res, &te)
	return te.TimeEntries
}

func NewClient(user, pass string) *Client {
	c := &Client{}
	c.User = user
	c.Pass = pass
	c.HTTPClient = &http.Client{}
	return c
}
