package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/OwnLocal/goes"
	"github.com/topmedia/go-autotask/client"
	"github.com/topmedia/go-autotask/entities"
)

var (
	user = flag.String("user", os.Getenv("AUTOTASK_USER"),
		"Username for API access")
	pass = flag.String("pass", os.Getenv("AUTOTASK_PASS"),
		"Password for API access")
	since = flag.String("since", time.Now().Format("2006-01-02"),
		"Fetch entries since this date (e.g. 2015-01-16)")
	es_host = flag.String("es_host", "localhost",
		"Elasticsearch Host")
	es_port = flag.String("es_port", "9200",
		"Elasticsearch Port")
	es_index = flag.String("es_index", "autotask",
		"Elasticsearch Index Name")
)

func CreateIndex(c *goes.Connection) {
	if exists, _ := c.IndicesExist([]string{*es_index}); exists {
		return
	}

	mapping := map[string]interface{}{
		"mappings": map[string]interface{}{
			"_default_": map[string]interface{}{
				"properties": map[string]interface{}{
					"role_name": map[string]interface{}{
						"type": "string",
						"fields": map[string]interface{}{
							"raw": map[string]interface{}{
								"index": "not_analyzed",
								"type":  "string",
							},
						},
					},
					"resource_name": map[string]interface{}{
						"type": "string",
						"fields": map[string]interface{}{
							"raw": map[string]interface{}{
								"index": "not_analyzed",
								"type":  "string",
							},
						},
					},
					"account_name": map[string]interface{}{
						"type": "string",
						"fields": map[string]interface{}{
							"raw": map[string]interface{}{
								"index": "not_analyzed",
								"type":  "string",
							},
						},
					},
					"title": map[string]interface{}{
						"type": "string",
						"fields": map[string]interface{}{
							"raw": map[string]interface{}{
								"index": "not_analyzed",
								"type":  "string",
							},
						},
					},
					"hours_worked": map[string]interface{}{
						"type": "float",
					},
				},
			},
		},
	}

	_, err := c.CreateIndex(*es_index, mapping)

	if err != nil {
		log.Fatal("Failed to create index with mapping: %s", err)
	}
}

func TimeEntryDocument(te *entities.TimeEntry, t *entities.Ticket, a *entities.Account) goes.Document {
	return goes.Document{
		Id:    fmt.Sprintf("%s-%d", t.TicketNumber, te.ID),
		Index: "autotask",
		Type:  "timeentry",
		Fields: map[string]interface{}{
			"account_name":  a.AccountName,
			"ticket_number": t.TicketNumber,
			"title":         t.Title,
			"resource_name": te.ResourceName,
			"role_name":     te.RoleName,
			"hours_worked":  te.HoursWorked,
			"date_time":     te.StartDateTime,
		},
	}
}

func main() {
	flag.Parse()

	if user == nil || pass == nil {
		log.Fatal("Please specify your Autotask username and password via-user and -pass")
	}

	api := client.NewClient(*user, *pass)
	es := goes.NewConnection(*es_host, *es_port)
	CreateIndex(es)

	tes := api.FetchTimeEntries(&entities.QueryExpression{Field: "CreateDateTime", Op: "GreaterThan", Value: *since})

	for k, te := range tes {
		if te.TicketID == 0 {
			continue
		}
		t := api.FetchTicketByID(te.TicketID)
		if t == nil {
			continue
		}
		a := api.FetchAccountByID(t.AccountID)
		tes[k].Ticket = t
		tes[k].Account = a

		_, err := es.Index(TimeEntryDocument(&te, t, a), url.Values{})

		if err != nil {
			log.Fatal("Indexing failed: %s", err)
		}

		fmt.Printf("%s #%s %s (%s: %v - %v)\n", a.AccountName, t.TicketNumber, t.Title,
			te.ResourceName, te.HoursWorked, te.RoleName)
	}

}
