package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/topmedia/timevacuum/entities"
)

var (
	user = flag.String("user", os.Getenv("AUTOTASK_USER"),
		"Username for API access")
	pass = flag.String("pass", os.Getenv("AUTOTASK_PASS"),
		"Password for API access")
	since = flag.String("since", time.Now().Format("2006-01-02"),
		"Fetch entries since this date (e.g. 2015-01-16)")
)

func main() {
	flag.Parse()

	if user == nil || pass == nil {
		log.Fatal("Please specify your Autotask username and password via-user and -pass")
	}

	api := NewClient(*user, *pass)

	tes := api.FetchTimeEntries(&entities.QueryExpression{Field: "CreateDateTime", Op: "GreaterThan", Value: *since})

	for k, te := range tes {
		if te.TicketID == 0 {
			continue
		}
		t := api.FetchTicketByID(te.TicketID)
		tes[k].Ticket = t
		fmt.Printf("#%s %s (%s: %v)\n", t.TicketNumber, t.Title,
			te.ResourceName, te.HoursWorked)
	}

}
