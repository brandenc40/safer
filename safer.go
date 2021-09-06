package safer

// NewClient build's a new Client interface
func NewClient() *Client {
	return &Client{
		scraper: scraper{},
	}
}

// Client for scraping company details from SAFER
type Client struct {
	scraper
}

// GetCompanyByDOTNumber - Get a company snapshot by the companies DOT number. Returns ErrCompanyNotFound if
// no company is found
func (c *Client) GetCompanyByDOTNumber(dotNumber string) (*CompanySnapshot, error) {
	return c.scraper.scrapeCompanySnapshot(paramUSDOT, dotNumber)
}

// GetCompanyByMCMX - Get a company snapshot by the companies MC/MX number. Returns ErrCompanyNotFound if no
// company is found.
//
// Note: do not include the prefix. (e.g. use "133655" not "MC-133655")
func (c *Client) GetCompanyByMCMX(mcmx string) (*CompanySnapshot, error) {
	return c.scraper.scrapeCompanySnapshot(paramMCMX, mcmx)
}

// SearchCompaniesByName - Search for all carriers with a given name. Name queries will return the best matched results
// in a slice of CompanyResult structs.
func (c *Client) SearchCompaniesByName(name string) ([]CompanyResult, error) {
	return c.scraper.scrapeCompanyNameSearch(name)
}
