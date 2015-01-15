package main

import (
	"encoding/xml"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/beevik/etree"
)

var (
	destURL     = "https://webservices7.autotask.net/ATServices/1.5/atws.asmx"
	soapAction  = "http://autotask.net/ATWS/v1_5/query"
	contentType = "text/xml"
)

type Envelope struct {
	XMLName xml.Name
	Body    Body
}

type Body struct {
	XMLName       xml.Name
	QueryResponse QueryResponse `xml:"queryResponse"`
}

type QueryResponse struct {
	XMLName     xml.Name
	QueryResult QueryResult `xml:"queryResult"`
}

type QueryResult struct {
	XMLName       xml.Name
	EntityResults EntityResults
}

type EntityResults struct {
	XMLName  xml.Name
	Entities []Entity `xml:"Entity"`
}

type Entity struct {
	XMLName      xml.Name
	ID           int `xml:"id"`
	HoursWorked  float32
	ResourceID   int
	ResourceName string
	Resources    map[int]string
}

func (e *Entity) CacheResources() {
}

type QueryXML struct {
	Doc  *etree.Document
	Qxml *etree.Element
}

func (q *QueryXML) ToReader() io.Reader {
	return strings.NewReader(q.String())
}

func (q *QueryXML) String() string {
	tmpl, err := ioutil.ReadFile(path.Join("templates", "query.xml"))

	if err != nil {
		return ""
	}

	out, err := q.Doc.WriteToString()

	if err != nil {
		return ""
	}

	return strings.Replace(string(tmpl), "{queryxml}", out, 1)
}

// func (q *QueryXML) String() string {
// 	q.Doc.Indent(2)
// 	s, _ := q.Doc.WriteToString()
// 	return s

// 	// if err != nil {
// 	// 	log.Fatalf("Failed to open request XML: %v", err)
// 	// }
// 	// reqBody, err := os.Open(path.Join("templates", "query.xml"))

// 	// defer reqBody.Close()
// 	// req

// }
func (q *QueryXML) Entity(name string) {
	e := q.Qxml.CreateElement("entity")
	e.SetText(name)
}
func (q *QueryXML) FieldExpression(name, op, value string) {
	qry := q.Qxml.CreateElement("query")
	f := qry.CreateElement("field")
	f.SetText(name)
	e := f.CreateElement("expression")
	e.SetText(value)
	e.CreateAttr("op", op)
}

func NewQueryXML() *QueryXML {
	q := &QueryXML{}
	q.Doc = etree.NewDocument()
	q.Qxml = q.Doc.CreateElement("queryxml")

	return q
}

type Client struct {
	HTTPClient *http.Client
	User       string
	Pass       string
	Response   *http.Response
}

func (c *Client) Request(q *QueryXML) *http.Response {
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

	return c.Response
}

func (c *Client) Body() []byte {
	body, err := ioutil.ReadAll(c.Response.Body)

	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}
	return body
}

func NewClient(user, pass string) *Client {
	c := &Client{}
	c.User = user
	c.Pass = pass
	c.HTTPClient = &http.Client{}
	return c
}
