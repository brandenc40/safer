package safer_test

import (
	"fmt"
	"log"

	"github.com/brandenc40/safer"
)

func ExampleClient_GetCompanyByDOTNumber() {
	client := safer.NewClient()

	snapshot, err := client.GetCompanyByDOTNumber("1003306")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%#v", snapshot)
}

func ExampleClient_GetCompanyByMCMX() {
	client := safer.NewClient()

	snapshot, err := client.GetCompanyByMCMX("133655")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%#v", snapshot)
}

func ExampleClient_SearchCompaniesByName() {
	client := safer.NewClient()

	res, err := client.SearchCompaniesByName("Schneider")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%#v", res[0])
}
