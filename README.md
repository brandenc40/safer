# SAFER Scraper API [GoLang]

[![Go Report Card](https://goreportcard.com/badge/github.com/brandenc40/go-safer)](https://goreportcard.com/report/github.com/brandenc40/go-safer)
[![Go Reference](https://pkg.go.dev/badge/github.com/brandenc40/safer.svg)](https://pkg.go.dev/github.com/brandenc40/safer)
[![codecov](https://codecov.io/gh/brandenc40/safer/branch/master/graph/badge.svg?token=4BSF2R1OGP)](https://codecov.io/gh/brandenc40/safer)
[![Tests](https://github.com/brandenc40/safer/actions/workflows/go.yml/badge.svg)](https://github.com/brandenc40/safer/actions/workflows/go.yml)

An API to scrape data from the Department of Transportation's Safety and Fitness Electronic Records 
([SAFER](https://safer.fmcsa.dot.gov/CompanySnapshot.aspx)) System.

Scraping is performed using [htmlquery](https://github.com/antchfx/htmlquery), this project's only non std lib dependency.


## Installation

```shell
go get -u github.com/brandenc40/safer
```

## Client Methods

```go
// GetCompanyByDOTNumber - Get a company snapshot by the companies DOT number. Returns ErrCompanyNotFound if
// no company is found
func (c *Client) GetCompanyByDOTNumber(dotNumber string) (*CompanySnapshot, error)

// GetCompanyByMCMX - Get a company snapshot by the companies MC/MX number. Returns ErrCompanyNotFound if no
// company is found.
//
// Note: do not include the prefix. (e.g. use "133655" not "MC-133655")
func (c *Client) GetCompanyByMCMX(mcmx string) (*CompanySnapshot, error)

// SearchCompaniesByName - Search for all carriers with a given name. Name queries will return the best matched results
// in a slice of CompanyResult structs.
func (c *Client) SearchCompaniesByName(name string) ([]CompanyResult, error)
```

### Build a new Client

```go
package main

import (
	"github.com/brandenc40/safer"
)

func main() {
	client := safer.NewClient()
	//... use the client
}
```

### Scraping Benchmark

Benchmarks only test the time taken to parse the html and map it back to the output. Server time is ignored here.

```shell 
goos: darwin
goarch: arm64
pkg: github.com/brandenc40/safer
BenchmarkClient_GetCompanyByDOTNumber-8             9145            132111 ns/op           91557 B/op       2672 allocs/op
BenchmarkClient_Search_4Results-8                  94274             12821 ns/op            9559 B/op        285 allocs/op
BenchmarkClient_Search_484Results-8                 1048           1137919 ns/op          826245 B/op      24774 allocs/op
PASS
ok      github.com/brandenc40/safer     4.942s
```
