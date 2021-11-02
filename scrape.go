package safer

import (
	"errors"
	"net/http"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

const (
	companySnapshotURL = "https://safer.fmcsa.dot.gov/query.asp"
	searchURL          = "https://safer.fmcsa.dot.gov/keywordx.asp"
	paramUSDOT         = "USDOT"
	paramMCMX          = "MC_MX"
)

var headers = http.Header{
	"Accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
	"Accept-Encoding":           {"gzip, deflate, br"},
	"Accept-Language":           {"en-US,en;q=0.9"},
	"Cache-Control":             {"max-age=0"},
	"Connection":                {"keep-alive"},
	"Host":                      {"safer.fmcsa.dot.gov"},
	"Upgrade-Insecure-Requests": {"1"},
	"User-Agent":                {"Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Mobile Safari/537.36"},
}

type scraper struct {
	companySnapshotURL string
	searchURL          string
}

func (s *scraper) scrapeCompanySnapshot(queryParam, queryString string) (*CompanySnapshot, error) {
	params := "?searchType=ANY&query_type=queryCarrierSnapshot&query_param=" + queryParam + "&query_string=" + queryString
	reqURL := companySnapshotURL
	if s.companySnapshotURL != "" {
		reqURL = s.companySnapshotURL
	}
	node, err := postRequestToHTMLNode(reqURL + params)
	if err != nil {
		return nil, err
	}
	return htmlNodeToCompanySnapshot(node)
}

func (s *scraper) scrapeCompanyNameSearch(queryString string) ([]CompanyResult, error) {
	params := "?SEARCHTYPE=&searchstring=*" + strings.ToUpper(queryString) + "*"
	reqURL := searchURL
	if s.searchURL != "" {
		reqURL = s.searchURL
	}
	node, err := postRequestToHTMLNode(reqURL + params)
	if err != nil {
		return nil, err
	}
	return htmlNodeToCompanyResults(node)
}

func postRequestToHTMLNode(reqURL string) (*html.Node, error) {
	req, err := http.NewRequest(http.MethodPost, reqURL, http.NoBody)
	if err != nil {
		return nil, err
	}
	req.Header = headers
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status + " Response from SAFER")
	}
	return htmlquery.Parse(resp.Body)
}
