package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/brandenc40/safer"
)

func main() {
	// by mc/mx
	s := time.Now()
	snapshot, err := safer.GetCompanyByMCMX("133655")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(time.Since(s).String())
	log.Printf("%+v", snapshot)

	// by dot
	s = time.Now()
	snapshot, err = safer.GetCompanyByDOTNumber("3653803")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(time.Since(s).String())
	v, _ := json.Marshal(snapshot)
	log.Printf(string(v))

	// search results and grab snapshot from result
	res, err := safer.SearchCompaniesByName("Schneider")
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("%#v", res[0])
	snapshot, err = res[0].GetSnapshot()
	if err != nil {
		log.Fatalln(err)
	}
}
