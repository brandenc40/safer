# SAFER Scraper API [GoLang]

[![Go Report Card](https://goreportcard.com/badge/github.com/brandenc40/go-safer)](https://goreportcard.com/report/github.com/brandenc40/go-safer)
[![Go Reference](https://pkg.go.dev/badge/github.com/brandenc40/safer.svg)](https://pkg.go.dev/github.com/brandenc40/safer)
[![codecov](https://codecov.io/gh/brandenc40/safer/branch/master/graph/badge.svg?token=4BSF2R1OGP)](https://codecov.io/gh/brandenc40/safer)
[![Tests](https://github.com/brandenc40/safer/actions/workflows/go.yml/badge.svg)](https://github.com/brandenc40/safer/actions/workflows/go.yml)

An API to scrape data from the Department of Transportation's Safety and Fitness Electronic Records 
([SAFER](https://safer.fmcsa.dot.gov/CompanySnapshot.aspx)) System.

Scaping is performed using [Colly](https://github.com/gocolly/colly), this project's only non standard library dependency.


## Installation

```shell
go get github.com/brandenc40/safer
```

## Client Interface

```go
// Client for scraping company details from SAFER
type Client interface {
    // GetCompanyByDOTNumber - Get a company snapshot by the companies DOT number. Returns ErrCompanyNotFound if
    // no company is found
    GetCompanyByDOTNumber(dotNumber string) (*CompanySnapshot, error)
    
    // GetCompanyByMCMX - Get a company snapshot by the companies MC/MX number. Returns ErrCompanyNotFound if no
    // company is found.
    //
    // Note: do not include the prefix. (e.g. use "133655" not "MC-133655")
    GetCompanyByMCMX(mcmx string) (*CompanySnapshot, error)
    
    // SearchCompaniesByName - Search for all carriers with a given name. Name queries will return the best matched results
    // in a slice of CompanyResult structs.
    SearchCompaniesByName(name string) ([]CompanyResult, error)
}
```

## Example Usage

```go
package main

import (
	"log"

	"github.com/brandenc40/safer"
)

func main() {
	client := safer.NewClient()
	
	// by mc/mx
	snapshot, err := client.GetCompanyByMCMX("133655")
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("%+v", snapshot)

	// by dot
	snapshot, err = client.GetCompanyByDOTNumber("264184")
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("%+v", snapshot)

	// search by name and grab snapshot from result
	companies, err := client.SearchCompaniesByName("Schneider")
	if err != nil {
		log.Fatalln(err)
	}
	topResult := companies[0]
	log.Printf("%#v", topResult)
	snapshot, err = client.GetCompanyByDOTNumber(topResult.DOTNumber)
	if err != nil {
		log.Fatalln(err)
	}
}
```
