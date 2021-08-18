package safer

import "github.com/gocolly/colly/v2"

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

// NewClient build's a new Client interface
func NewClient() Client {
	c := &client{
		scraper: &scraper{
			baseCollector:      colly.NewCollector(),
			companySnapshotURL: companySnapshotURL,
			searchURL:          searchURL,
		},
	}
	return c
}

type client struct {
	scraper *scraper
}

func (c *client) GetCompanyByDOTNumber(dotNumber string) (*CompanySnapshot, error) {
	return c.scraper.scrapeCompanySnapshot(paramUSDOT, dotNumber)
}

func (c *client) GetCompanyByMCMX(mcmx string) (*CompanySnapshot, error) {
	return c.scraper.scrapeCompanySnapshot(paramMCMX, mcmx)
}

func (c *client) SearchCompaniesByName(name string) ([]CompanyResult, error) {
	return c.scraper.scrapeCompanyNameSearch(name)
}
