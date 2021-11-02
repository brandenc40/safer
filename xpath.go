package safer

import (
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

// company snapshot xpath constants
const (
	snapshotNotFoundXpath      = "/html/head/title[text()='SAFER Web - Company Snapshot RECORD NOT FOUND' or text()='SAFER Web - Company Snapshot RECORD INACTIVE']"
	srcTableXpath              = "/html/body/p/table/tbody/tr[2]/td/table/tbody/tr[2]/td"
	latestUpdateDateXpath      = "/table/tbody/tr[3]/td/font/b[3]/font/text()"
	tableGeneralInfoXpath      = "/center[1]/table/tbody"
	tableOperationClassXpath   = "/tr[14]/td/table/tbody/tr[2]/td/table/tbody/tr[.//td[@class='queryfield']/text() = 'X']"
	tableCarrierOpXpath        = "/tr[16]/td/table/tbody/tr[2]/td/table/tbody/tr[.//td[@class='queryfield']/text() = 'X']/td/font/text()"
	tableCargoCarriedXpath     = "/tr[19]/td/table/tbody/tr[2]/td/table/tbody/tr[.//td[@class='queryfield']/text() = 'X']"
	tableUSInspectionXpath     = "/center[3]/table/tbody"
	tableUSCrashXpath          = "/center[4]/table/tbody/tr[2]"
	tableCanadaInspectionXpath = "/center[6]/table/tbody"
	tableCanadaCrashXpath      = "/center[7]/table/tbody/tr[2]/td/text()"
	tableSafetyRatingXpath     = "/center[9]/table/tbody"
)

// company search xpath constants
const (
	companyResultXpath = "/html/body/table[3]/tbody/tr[.//*[@scope='rpw']]"
)

func htmlNodeToCompanySnapshot(root *html.Node) (*CompanySnapshot, error) {
	if found := htmlquery.Find(root, snapshotNotFoundXpath); found != nil && len(found) > 0 {
		return nil, ErrCompanyNotFound
	}
	snapshot := new(CompanySnapshot)
	if srcNode := htmlquery.FindOne(root, srcTableXpath); srcNode != nil {
		snapshot.LatestUpdateDate = parseDate(getNodeText(srcNode, latestUpdateDateXpath))
		// general info
		if node := htmlquery.FindOne(srcNode, tableGeneralInfoXpath); node != nil {
			snapshot.EntityType = getNodeText(node, "/tr[2]/td/text()")
			if tr3 := htmlquery.FindOne(node, "/tr[3]"); tr3 != nil {
				snapshot.OutOfServiceDate = parseDate(getNodeText(tr3, "/td[2]/text()"))
				snapshot.OperatingStatus = getNodeText(tr3, "/td[1]/text()")
				if snapshot.OperatingStatus == "" {
					// out-of-service is bolded and not caught by the previous xpath
					snapshot.OperatingStatus = getNodeText(tr3, "/td[1]/font/b/text()")
				}
			}
			snapshot.LegalName = getNodeText(node, "/tr[4]/td/text()")
			snapshot.DBAName = getNodeText(node, "/tr[5]/td/text()")
			snapshot.PhysicalAddress = parseAddress(getNodeTexts(node, "/tr[6]/td/text()")...)
			snapshot.Phone = getNodeText(node, "/tr[7]/td/text()")
			snapshot.MailingAddress = parseAddress(getNodeTexts(node, "/tr[8]/td/text()")...)
			if tr9 := htmlquery.FindOne(node, "/tr[9]"); tr9 != nil {
				snapshot.DOTNumber = getNodeText(tr9, "/td[1]/text()")
				snapshot.StateCarrierID = getNodeText(tr9, "/td[2]/text()")
			}
			if tr10 := htmlquery.FindOne(node, "/tr[10]"); tr10 != nil {
				snapshot.MCMXFFNumbers = getNodeTexts(tr10, "/td[1]/a/text()")
				snapshot.DUNSNumber = getNodeText(tr10, "/td[2]/text()")
				if snapshot.DUNSNumber == "--" {
					snapshot.DUNSNumber = ""
				}
			}
			if tr11 := htmlquery.FindOne(node, "/tr[11]"); tr11 != nil {
				snapshot.PowerUnits = parseInt(getNodeText(tr11, "/td[1]/text()"))
				snapshot.Drivers = parseInt(getNodeText(tr11, "/td[2]/font/b/text()"))
			}
			if tr12 := htmlquery.FindOne(node, "/tr[12]"); tr12 != nil {
				snapshot.MCS150FormDate = parseDate(getNodeText(tr12, "/td[1]/text()"))
				snapshot.MCS150Mileage, snapshot.MCS150Year = parseMCS150MileageYear(getNodeText(tr12, "/td[2]/font/b/text()"))
			}
			// carrier classification
			for _, classNode := range htmlquery.Find(node, tableOperationClassXpath) {
				classification := getNodeText(classNode, "/td/font/text()")
				if classification == "" {
					// optional extra classifications (not all will have this)
					classification = getNodeText(classNode, "/td[2]/text()")
				}
				if classification != "" {
					snapshot.OperationClassification = append(snapshot.OperationClassification, classification)
				}
			}
			// carrier operation
			operations := getNodeTexts(node, tableCarrierOpXpath)
			for _, op := range operations {
				snapshot.CarrierOperation = append(snapshot.CarrierOperation, op)
			}
			// cargo carried
			for _, cargoNode := range htmlquery.Find(node, tableCargoCarriedXpath) {
				cargo := getNodeText(cargoNode, "/td/font/text()")
				if cargo == "" {
					// optional extra classifications (not all will have this)
					cargo = getNodeText(cargoNode, "/td[2]/text()")
				}
				if cargo != "" {
					snapshot.CargoCarried = append(snapshot.CargoCarried, cargo)
				}
			}
		}
		// us inspections
		if node := htmlquery.FindOne(srcNode, tableUSInspectionXpath); node != nil {
			if tr2 := htmlquery.FindOne(node, "/tr[2]"); tr2 != nil {
				snapshot.USVehicleInspections.Inspections = parseInt(getNodeText(tr2, "/td[1]/text()"))
				snapshot.USDriverInspections.Inspections = parseInt(getNodeText(tr2, "/td[2]/text()"))
				snapshot.USHazmatInspections.Inspections = parseInt(getNodeText(tr2, "/td[3]/text()"))
				snapshot.USIEPInspections.Inspections = parseInt(getNodeText(tr2, "/td[4]/text()"))
			}
			if tr3 := htmlquery.FindOne(node, "/tr[3]"); tr3 != nil {
				snapshot.USVehicleInspections.OutOfService = parseInt(getNodeText(tr3, "/td[1]/text()"))
				snapshot.USDriverInspections.OutOfService = parseInt(getNodeText(tr3, "/td[2]/text()"))
				snapshot.USHazmatInspections.OutOfService = parseInt(getNodeText(tr3, "/td[3]/text()"))
				snapshot.USIEPInspections.OutOfService = parseInt(getNodeText(tr3, "/td[4]/text()"))
			}
			if tr4 := htmlquery.FindOne(node, "/tr[4]"); tr4 != nil {
				snapshot.USVehicleInspections.OutOfServicePct = parsePctToFloat32(getNodeText(tr4, "/td[1]/text()"))
				snapshot.USDriverInspections.OutOfServicePct = parsePctToFloat32(getNodeText(tr4, "/td[2]/text()"))
				snapshot.USHazmatInspections.OutOfServicePct = parsePctToFloat32(getNodeText(tr4, "/td[3]/text()"))
				snapshot.USIEPInspections.OutOfServicePct = parsePctToFloat32(getNodeText(tr4, "/td[4]/text()"))
			}
			if tr5 := htmlquery.FindOne(node, "/tr[5]"); tr5 != nil {
				snapshot.USVehicleInspections.NationalAverage = parsePctToFloat32(getNodeText(tr5, "/td[1]/font/text()"))
				snapshot.USDriverInspections.NationalAverage = parsePctToFloat32(getNodeText(tr5, "/td[2]/font/text()"))
				snapshot.USHazmatInspections.NationalAverage = parsePctToFloat32(getNodeText(tr5, "/td[3]/font/text()"))
				snapshot.USIEPInspections.NationalAverage = parsePctToFloat32(getNodeText(tr5, "/td[4]/font/text()"))
			}
		}
		// us crash
		if node := htmlquery.FindOne(srcNode, tableUSCrashXpath); node != nil {
			snapshot.USCrashes.Fatal = parseInt(getNodeText(node, "/td[1]/text()"))
			snapshot.USCrashes.Injury = parseInt(getNodeText(node, "/td[2]/text()"))
			snapshot.USCrashes.Tow = parseInt(getNodeText(node, "/td[3]/text()"))
			snapshot.USCrashes.Total = parseInt(getNodeText(node, "/td[4]/text()"))
		}
		// canada inspection
		if node := htmlquery.FindOne(srcNode, tableCanadaInspectionXpath); node != nil {
			if tr2 := htmlquery.FindOne(node, "/tr[2]"); tr2 != nil {
				snapshot.CanadaVehicleInspections.Inspections = parseInt(getNodeText(tr2, "/td[1]/text()"))
				snapshot.CanadaDriverInspections.Inspections = parseInt(getNodeText(tr2, "/td[2]/text()"))
			}
			if tr3 := htmlquery.FindOne(node, "/tr[3]"); tr3 != nil {
				snapshot.CanadaVehicleInspections.OutOfService = parseInt(getNodeText(tr3, "/td[1]/text()"))
				snapshot.CanadaDriverInspections.OutOfService = parseInt(getNodeText(tr3, "/td[2]/text()"))
			}
			if tr4 := htmlquery.FindOne(node, "/tr[4]"); tr4 != nil {
				snapshot.CanadaVehicleInspections.OutOfServicePct = parsePctToFloat32(getNodeText(tr4, "/td[1]/text()"))
				snapshot.CanadaDriverInspections.OutOfServicePct = parsePctToFloat32(getNodeText(tr4, "/td[2]/text()"))
			}
		}
		// canada crash
		if nodes := htmlquery.Find(srcNode, tableCanadaCrashXpath); nodes != nil {
			snapshot.CanadaCrashes.Fatal = parseInt(strings.TrimSpace(nodes[0].Data))
			snapshot.CanadaCrashes.Injury = parseInt(strings.TrimSpace(nodes[1].Data))
			snapshot.CanadaCrashes.Tow = parseInt(strings.TrimSpace(nodes[2].Data))
			snapshot.CanadaCrashes.Total = parseInt(strings.TrimSpace(nodes[3].Data))
		}
		// canada crash
		if node := htmlquery.FindOne(srcNode, tableSafetyRatingXpath); node != nil {
			if tr2 := htmlquery.Find(node, "/tr[2]/td/text()"); tr2 != nil {
				snapshot.Safety.RatingDate = parseDate(strings.TrimSpace(tr2[0].Data))
				snapshot.Safety.ReviewDate = parseDate(strings.TrimSpace(tr2[1].Data))
			}
			if tr3 := htmlquery.Find(node, "/tr[3]/td/text()"); tr3 != nil {
				snapshot.Safety.Rating = strings.TrimSpace(tr3[0].Data)
				snapshot.Safety.Type = strings.TrimSpace(tr3[1].Data)
			}
		}
	}
	return snapshot, nil
}

func htmlNodeToCompanyResults(node *html.Node) ([]CompanyResult, error) {
	resultNodes := htmlquery.Find(node, companyResultXpath)
	if resultNodes == nil || len(resultNodes) == 0 {
		return []CompanyResult{}, nil
	}
	companyResults := make([]CompanyResult, len(resultNodes), len(resultNodes))
	for i, n := range resultNodes {
		companyResults[i] = companyResultFromNode(n)
	}
	return companyResults, nil
}

func companyResultFromNode(n *html.Node) CompanyResult {
	res := CompanyResult{
		Location: getNodeText(n, "/td/b/text()"),
	}
	if node := htmlquery.FindOne(n, "/th/b/a"); node != nil {
		res.Name = getNodeText(node, "/text()")
		res.DOTNumber = parseDotFromSearchParams(htmlquery.SelectAttr(node, "href"))
	}
	return res
}

func getNodeText(node *html.Node, path string) string {
	child := htmlquery.FindOne(node, path)
	if child == nil {
		return ""
	}
	return strings.TrimSpace(child.Data)
}

func getNodeTexts(node *html.Node, path string) []string {
	children := htmlquery.Find(node, path)
	if children == nil || len(children) == 0 {
		return []string{}
	}
	out := make([]string, len(children), len(children))
	for i, child := range children {
		out[i] = strings.TrimSpace(child.Data)
	}
	return out
}
