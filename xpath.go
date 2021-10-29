package safer

import (
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

// function to map the XMLElement to the CompanySnapshot
type mapperFunc func(node *html.Node, snapshot *CompanySnapshot)

// company snapshot xpath constants
const (
	snapshotNotFoundXpath = "/html/head/title[text()='SAFER Web - Company Snapshot RECORD NOT FOUND' or text()='SAFER Web - Company Snapshot RECORD INACTIVE']"
	latestUpdateDateXpath = "//b/font[@color='#0000C0']/text()"
	tableXpath            = "//table"
)

// company search xpath constants
const (
	companyResultXpath = "//tr[.//*[@scope='rpw']]"
)

// company snapshot tableXpath indexes
const (
	tableIdxGeneralInfo       = 6
	tableIdxOperationClass    = 7
	tableIdxCarrierOp         = 11
	tableIdxCargoCarried      = 15
	tableIdxUsInspections     = 19
	tableIdxUsCrashes         = 20
	tableIdxCanadaInspections = 21
	tableIdxCanadaCrashes     = 22
	tableIdxSafetyRating      = 23
)

// key=index of tableXpath element to find the data, value=function to map the XMLElement to the CompanySnapshot
var snapshotTableXpathMapping = map[int]mapperFunc{
	tableIdxGeneralInfo: func(node *html.Node, snapshot *CompanySnapshot) {
		snapshot.EntityType = getNodeText(node, "//tr[2]/td/text()")
		snapshot.LegalName = getNodeText(node, "//tr[4]/td/text()")
		snapshot.DBAName = getNodeText(node, "//tr[5]/td/text()")
		snapshot.PhysicalAddress = parseAddress(getNodeTexts(node, "//tr[6]/td/text()")...)
		snapshot.Phone = getNodeText(node, "//tr[7]/td/text()")
		snapshot.MailingAddress = parseAddress(getNodeTexts(node, "//tr[8]/td/text()")...)
		snapshot.DOTNumber = getNodeText(node, "//tr[9]/td[1]/text()")
		snapshot.StateCarrierID = getNodeText(node, "//tr[9]/td[2]/text()")
		snapshot.MCMXFFNumbers = getNodeTexts(node, "//tr[10]/td[1]/a/text()")
		snapshot.DUNSNumber = getNodeText(node, "//tr[10]/td[2]/text()")
		if snapshot.DUNSNumber == "--" {
			snapshot.DUNSNumber = ""
		}
		snapshot.PowerUnits = parseInt(getNodeText(node, "//tr[11]/td[1]/text()"))
		snapshot.Drivers = parseInt(getNodeText(node, "//tr[11]/td[2]/font/b/text()"))
		snapshot.MCS150FormDate = parseDate(getNodeText(node, "//tr[12]/td[1]/text()"))
		snapshot.MCS150Mileage, snapshot.MCS150Year = parseMCS150MileageYear(getNodeText(node, "//tr[12]/td[2]/font/b/text()"))
		snapshot.OutOfServiceDate = parseDate(getNodeText(node, "//tr[3]/td[2]/text()"))
		snapshot.OperatingStatus = getNodeText(node, "//tr[3]/td[1]/text()")
		if snapshot.OperatingStatus == "" {
			// out-of-service is bolded and not caught by the previous xpath
			snapshot.OperatingStatus = getNodeText(node, "//tr[3]/td[1]/font/b/text()")
		}
	},
	tableIdxOperationClass: func(node *html.Node, snapshot *CompanySnapshot) {
		classifications := getNodeTexts(node, "//tr[2]/td/table/tbody/tr[.//td[@class='queryfield']/text() = 'X']/td/font/text()")
		for _, classification := range classifications {
			snapshot.OperationClassification = append(snapshot.OperationClassification, classification)
		}
		// optional extra classifications (not all will have this)
		classifications = getNodeTexts(node, "//tr[2]/td/table/tbody/tr[.//td[@class='queryfield']/text() = 'X']/td[2]/text()")
		for _, classification := range classifications {
			snapshot.OperationClassification = append(snapshot.OperationClassification, classification)
		}
	},
	tableIdxCarrierOp: func(node *html.Node, snapshot *CompanySnapshot) {
		operations := getNodeTexts(node, "//tr[2]/td/table/tbody/tr[.//td[@class='queryfield']/text() = 'X']/td/font/text()")
		for _, op := range operations {
			snapshot.CarrierOperation = append(snapshot.CarrierOperation, op)
		}
	},
	tableIdxCargoCarried: func(node *html.Node, snapshot *CompanySnapshot) {
		cargos := getNodeTexts(node, "//tr[2]/td/table/tbody/tr[.//td[@class='queryfield']/text() = 'X']/td/font/text()")
		for _, cargo := range cargos {
			snapshot.CargoCarried = append(snapshot.CargoCarried, cargo)
		}
		// optional extra cargos (not all will have this)
		cargos = getNodeTexts(node, "//tr[2]/td/table/tbody/tr[.//td[@class='queryfield']/text() = 'X']/td[2]/text()")
		for _, cargo := range cargos {
			snapshot.CargoCarried = append(snapshot.CargoCarried, cargo)
		}
	},
	tableIdxUsInspections: func(node *html.Node, snapshot *CompanySnapshot) {
		snapshot.USVehicleInspections.Inspections = parseInt(getNodeText(node, "//tr[2]/td[1]/text()"))
		snapshot.USVehicleInspections.OutOfService = parseInt(getNodeText(node, "//tr[3]/td[1]/text()"))
		snapshot.USVehicleInspections.OutOfServicePct = parsePctToFloat32(getNodeText(node, "//tr[4]/td[1]/text()"))
		snapshot.USVehicleInspections.NationalAverage = parsePctToFloat32(getNodeText(node, "//tr[5]/td[1]/font/text()"))
		snapshot.USDriverInspections.Inspections = parseInt(getNodeText(node, "//tr[2]/td[2]/text()"))
		snapshot.USDriverInspections.OutOfService = parseInt(getNodeText(node, "//tr[3]/td[2]/text()"))
		snapshot.USDriverInspections.OutOfServicePct = parsePctToFloat32(getNodeText(node, "//tr[4]/td[2]/text()"))
		snapshot.USDriverInspections.NationalAverage = parsePctToFloat32(getNodeText(node, "//tr[5]/td[2]/font/text()"))
		snapshot.USHazmatInspections.Inspections = parseInt(getNodeText(node, "//tr[2]/td[3]/text()"))
		snapshot.USHazmatInspections.OutOfService = parseInt(getNodeText(node, "//tr[3]/td[3]/text()"))
		snapshot.USHazmatInspections.OutOfServicePct = parsePctToFloat32(getNodeText(node, "//tr[4]/td[3]/text()"))
		snapshot.USHazmatInspections.NationalAverage = parsePctToFloat32(getNodeText(node, "//tr[5]/td[3]/font/text()"))
		snapshot.USIEPInspections.Inspections = parseInt(getNodeText(node, "//tr[2]/td[4]/text()"))
		snapshot.USIEPInspections.OutOfService = parseInt(getNodeText(node, "//tr[3]/td[4]/text()"))
		snapshot.USIEPInspections.OutOfServicePct = parsePctToFloat32(getNodeText(node, "//tr[4]/td[4]/text()"))
		snapshot.USIEPInspections.NationalAverage = parsePctToFloat32(getNodeText(node, "//tr[5]/td[4]/font/text()"))
	},
	tableIdxUsCrashes: func(node *html.Node, snapshot *CompanySnapshot) {
		snapshot.USCrashes.Fatal = parseInt(getNodeText(node, "//tr[2]/td[1]/text()"))
		snapshot.USCrashes.Injury = parseInt(getNodeText(node, "//tr[2]/td[2]/text()"))
		snapshot.USCrashes.Tow = parseInt(getNodeText(node, "//tr[2]/td[3]/text()"))
		snapshot.USCrashes.Total = parseInt(getNodeText(node, "//tr[2]/td[4]/text()"))
	},
	tableIdxCanadaInspections: func(node *html.Node, snapshot *CompanySnapshot) {
		snapshot.CanadaVehicleInspections.Inspections = parseInt(getNodeText(node, "//tr[2]/td[1]/text()"))
		snapshot.CanadaVehicleInspections.OutOfService = parseInt(getNodeText(node, "//tr[3]/td[1]/text()"))
		snapshot.CanadaVehicleInspections.OutOfServicePct = parsePctToFloat32(getNodeText(node, "//tr[4]/td[1]/text()"))
		snapshot.CanadaDriverInspections.Inspections = parseInt(getNodeText(node, "//tr[2]/td[2]/text()"))
		snapshot.CanadaDriverInspections.OutOfService = parseInt(getNodeText(node, "//tr[3]/td[2]/text()"))
		snapshot.CanadaDriverInspections.OutOfServicePct = parsePctToFloat32(getNodeText(node, "//tr[4]/td[2]/text()"))
	},
	tableIdxCanadaCrashes: func(node *html.Node, snapshot *CompanySnapshot) {
		snapshot.CanadaCrashes.Fatal = parseInt(getNodeText(node, "//tr[2]/td[1]/text()"))
		snapshot.CanadaCrashes.Injury = parseInt(getNodeText(node, "//tr[2]/td[2]/text()"))
		snapshot.CanadaCrashes.Tow = parseInt(getNodeText(node, "//tr[2]/td[3]/text()"))
		snapshot.CanadaCrashes.Total = parseInt(getNodeText(node, "//tr[2]/td[4]/text()"))
	},
	tableIdxSafetyRating: func(node *html.Node, snapshot *CompanySnapshot) {
		snapshot.Safety.RatingDate = parseDate(getNodeText(node, "//tr[2]/td[1]/text()"))
		snapshot.Safety.ReviewDate = parseDate(getNodeText(node, "//tr[2]/td[2]/text()"))
		snapshot.Safety.Rating = getNodeText(node, "//tr[3]/td[1]/text()")
		snapshot.Safety.Type = getNodeText(node, "//tr[3]/td[2]/text()")
	},
}

func companyResultFromNode(n *html.Node) CompanyResult {
	return CompanyResult{
		Name:      getNodeText(n, "/th/b/a/text()"),
		DOTNumber: parseDotFromSearchParams(getNodeAttrText(n, "/th/b/a/@href", "href")),
		Location:  getNodeText(n, "/td/b/text()"),
	}
}

func getNodeText(node *html.Node, path string) string {
	child := htmlquery.FindOne(node, path)
	if child == nil {
		return ""
	}
	return strings.TrimSpace(child.Data)
}

func getNodeAttrText(node *html.Node, path, attr string) string {
	child := htmlquery.FindOne(node, path)
	if child == nil {
		return ""
	}
	return strings.TrimSpace(htmlquery.SelectAttr(child, attr))
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
