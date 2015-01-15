package main

import (
	"encoding/xml"
	"flag"
	"fmt"
)

var (
	user = flag.String("user", "", "Username for API access")
	pass = flag.String("pass", "", "Password for API access")
)

func main() {
	flag.Parse()

	q := NewQueryXML()
	q.Entity("timeentry")
	q.FieldExpression("dateworked", "greaterthan", "2015-01-14")

	c := NewClient(*user, *pass)
	c.Request(q)

	var env Envelope

	xml.Unmarshal(c.Body(), &env)
	fmt.Printf("%v", env)
}
