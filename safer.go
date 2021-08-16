package safer

// GetCompanyByDOTNumber - Get a company snapshot by the companies DOT number
func GetCompanyByDOTNumber(dotNumber string) (*CompanySnapshot, error) {
	return scrapeCompanySnapshot(paramUSDOT, dotNumber, companySnapshotURL)
}

// GetCompanyByMCMX - Get a company snapshot by the companies MC/MX number
//
// Note: do not include the prefix. (e.g. use "133655" not "MC-133655")
func GetCompanyByMCMX(mcmx string) (*CompanySnapshot, error) {
	return scrapeCompanySnapshot(paramMCMX, mcmx, companySnapshotURL)
}

// SearchCompaniesByName - Search for all carriers with a given name. Name queries will return the best matched results
// in a slice of CompanyResult structs.
func SearchCompaniesByName(name string) ([]CompanyResult, error) {
	return scrapeCompanyNameSearch(name, searchURL)
}
