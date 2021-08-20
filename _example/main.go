package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/brandenc40/safer"
)

func main() {
	client := safer.NewClient()

	//// by mc/mx
	//s := time.Now()
	//snapshot, err := client.GetCompanyByMCMX("133655")
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//log.Println(time.Since(s).String())
	//log.Printf("%+v", snapshot)

	// by dot
	s := time.Now()
	snapshot, err := client.GetCompanyByDOTNumber("1003306")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(time.Since(s).String())
	v, _ := json.Marshal(snapshot)
	log.Printf(string(v))

	//// search results and grab snapshot from result
	//res, err := client.SearchCompaniesByName("Schneider")
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//log.Printf("%#v", res[0])
	//snapshot, err = client.GetCompanyByDOTNumber(res[0].DOTNumber)
	//if err != nil {
	//	log.Fatalln(err)
	//}
}
