package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/topmedia/timevacuum/entities"
)

var (
	user = flag.String("user", os.Getenv("AUTOTASK_USER"), "Username for API access")
	pass = flag.String("pass", os.Getenv("AUTOTASK_PASS"), "Password for API access")
)

func main() {
	flag.Parse()

	if user == nil || pass == nil {
		log.Fatal("Please specify your Autotask username and password via-user and -pass")
	}

	api := NewClient(*user, *pass)

	rr := api.FetchResources(&entities.QueryExpression{Field: "active", Op: "equals", Value: "true"})

	for _, r := range rr {
		fmt.Printf("%v\n", r.FirstName)
	}

	te := api.FetchTimeEntries(&entities.QueryExpression{Field: "dateworked", Op: "greaterthan", Value: "2015-01-14"})

	for _, r := range te {
		fmt.Printf("%v\n", r.HoursWorked)
	}

}
