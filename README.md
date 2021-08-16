# SAFER Scraper API [GoLang]

[![Go Report Card](https://goreportcard.com/badge/github.com/brandenc40/go-safer)](https://goreportcard.com/report/github.com/brandenc40/go-safer)
[![Go Reference](https://pkg.go.dev/badge/github.com/brandenc40/safer.svg)](https://pkg.go.dev/github.com/brandenc40/safer)
[![codecov](https://codecov.io/gh/brandenc40/safer/branch/master/graph/badge.svg?token=4BSF2R1OGP)](https://codecov.io/gh/brandenc40/safer)

A web scraping API to fetch data from the Department of Transportation's Safety and Fitness Electronic Records 
([SAFER](https://safer.fmcsa.dot.gov/CompanySnapshot.aspx)) System.

Scaping is performed using [Colly](https://github.com/gocolly/colly), this project's only non standard library dependency.


## Installation

```shell
go get github.com/brandenc40/safer
```

### Available Functions

```go
package safer

// GetCompanyByDOTNumber - Get a company snapshot by the companies DOT number
func GetCompanyByDOTNumber(dotNumber string) (*CompanySnapshot, error) 

// GetCompanyByMCMX - Get a company snapshot by the companies MC/MX number
//
// Note: do not include the prefix. (e.g. use "133655" not "MC-133655")
func GetCompanyByMCMX(mcmx string) (*CompanySnapshot, error) 

// SearchCompaniesByName - Search for all carriers with a given name. Name queries will return the best matched results
// in a slice of CompanyResult structs.
func SearchCompaniesByName(name string) ([]CompanyResult, error) 
```

### Example Usage

```go
package main

import (
	"log"

	"github.com/brandenc40/safer"
)

func main() {
	// by mc/mx
	snapshot, err := safer.GetCompanyByMCMX("133655")
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("%+v", snapshot)

	// by dot
	snapshot, err = safer.GetCompanyByDOTNumber("264184")
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("%+v", snapshot)

	// search by name and grab snapshot from result
	companies, err := safer.SearchCompaniesByName("Schneider")
	if err != nil {
		log.Fatalln(err)
	}
	topResult := companies[0]
	log.Printf("%#v", topResult)
	snapshot, err = topResult.GetSnapshot()
	if err != nil {
		log.Fatalln(err)
	}
}
```
