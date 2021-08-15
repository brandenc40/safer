package safer

import (
	"github.com/gocolly/colly/v2"
)

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

type mapperFunc func(element *colly.XMLElement, snapshot *CompanySnapshot)

var snapshotTableXMLMapping = map[int]mapperFunc{
	tableIdxGeneralInfo: func(element *colly.XMLElement, snapshot *CompanySnapshot) {
		snapshot.EntityType = element.ChildText("//tr[2]/td/text()")
		snapshot.LegalName = element.ChildText("//tr[4]/td/text()")
		snapshot.DBAName = element.ChildText("//tr[5]/td/text()")
		snapshot.PhysicalAddress = parseAddress(element.ChildTexts("//tr[6]/td/text()")...)
		snapshot.Phone = element.ChildText("//tr[7]/td/text()")
		snapshot.MailingAddress = parseAddress(element.ChildTexts("//tr[8]/td/text()")...)
		snapshot.DOTNumber = element.ChildText("//tr[9]/td[1]/text()")
		snapshot.StateCarrierID = element.ChildText("//tr[9]/td[2]/text()")
		snapshot.MCMXFFNumbers = element.ChildText("//tr[10]/td[1]/a/text()")
		snapshot.DUNSNumber = element.ChildText("//tr[10]/td[2]/text()")
		snapshot.PowerUnits = parseInt(element.ChildText("//tr[11]/td[1]/text()"))
		snapshot.Drivers = parseInt(element.ChildText("//tr[11]/td[2]/font/b/text()"))
		snapshot.MCS150FormDate = parseDate(element.ChildText("//tr[12]/td[1]/text()"))
		snapshot.MCS150Mileage, snapshot.MCS150Year = parseMCS150MileageYear(element.ChildText("//tr[12]/td[2]/font/b/text()"))
		snapshot.OutOfServiceDate = parseDate(element.ChildText("//tr[3]/td[2]/text()"))
		snapshot.OperatingStatus = element.ChildText("//tr[3]/td[1]/text()")
		if snapshot.OperatingStatus == "" {
			snapshot.OperatingStatus = element.ChildText("//tr[3]/td[1]/font/b/text()")
		}
	},
	tableIdxOperationClass: func(element *colly.XMLElement, snapshot *CompanySnapshot) {
		classifications := element.ChildTexts("//tr[2]/td/table/tbody/tr[.//td[@class='queryfield']/text() = 'X']/td/font/text()")
		for _, classification := range classifications {
			snapshot.OperationClassification = append(snapshot.OperationClassification, classification)
		}
		if lastVal := element.ChildText("//tr[2]/td[3]/table/tr[5]/td[2]/text()"); lastVal != "" {
			snapshot.OperationClassification = append(snapshot.OperationClassification, lastVal)
		}
	},
	tableIdxCarrierOp: func(element *colly.XMLElement, snapshot *CompanySnapshot) {
		operations := element.ChildTexts("//tr[2]/td/table/tbody/tr[.//td[@class='queryfield']/text() = 'X']/td/font/text()")
		for _, op := range operations {
			snapshot.CarrierOperation = append(snapshot.CarrierOperation, op)
		}
	},
	tableIdxCargoCarried: func(element *colly.XMLElement, snapshot *CompanySnapshot) {
		cargos := element.ChildTexts("//tr[2]/td/table/tbody/tr[.//td[@class='queryfield']/text() = 'X']/td/font/text()")
		for _, cargo := range cargos {
			snapshot.CargoCarried = append(snapshot.CargoCarried, cargo)
		}
	},
	tableIdxUsInspections: func(element *colly.XMLElement, snapshot *CompanySnapshot) {
		snapshot.USVehicleInspections.Inspections = parseInt(element.ChildText("//tr[2]/td[1]/text()"))
		snapshot.USVehicleInspections.OutOfService = parseInt(element.ChildText("//tr[3]/td[1]/text()"))
		snapshot.USVehicleInspections.OutOfServicePct = parsePctToFloat32(element.ChildText("//tr[4]/td[1]/text()"))
		snapshot.USVehicleInspections.NationalAverage = parsePctToFloat32(element.ChildText("//tr[5]/td[1]/font/text()"))
		snapshot.USDriverInspections.Inspections = parseInt(element.ChildText("//tr[2]/td[2]/text()"))
		snapshot.USDriverInspections.OutOfService = parseInt(element.ChildText("//tr[3]/td[2]/text()"))
		snapshot.USDriverInspections.OutOfServicePct = parsePctToFloat32(element.ChildText("//tr[4]/td[2]/text()"))
		snapshot.USDriverInspections.NationalAverage = parsePctToFloat32(element.ChildText("//tr[5]/td[2]/font/text()"))
		snapshot.USHazmatInspections.Inspections = parseInt(element.ChildText("//tr[2]/td[3]/text()"))
		snapshot.USHazmatInspections.OutOfService = parseInt(element.ChildText("//tr[3]/td[3]/text()"))
		snapshot.USHazmatInspections.OutOfServicePct = parsePctToFloat32(element.ChildText("//tr[4]/td[3]/text()"))
		snapshot.USHazmatInspections.NationalAverage = parsePctToFloat32(element.ChildText("//tr[5]/td[3]/font/text()"))
		snapshot.USIEPInspections.Inspections = parseInt(element.ChildText("//tr[2]/td[4]/text()"))
		snapshot.USIEPInspections.OutOfService = parseInt(element.ChildText("//tr[3]/td[4]/text()"))
		snapshot.USIEPInspections.OutOfServicePct = parsePctToFloat32(element.ChildText("//tr[4]/td[4]/text()"))
		snapshot.USIEPInspections.NationalAverage = parsePctToFloat32(element.ChildText("//tr[5]/td[4]/font/text()"))
	},
	tableIdxUsCrashes: func(element *colly.XMLElement, snapshot *CompanySnapshot) {
		snapshot.USCrashes.Fatal = parseInt(element.ChildText("//tr[2]/td[1]/text()"))
		snapshot.USCrashes.Injury = parseInt(element.ChildText("//tr[2]/td[2]/text()"))
		snapshot.USCrashes.Tow = parseInt(element.ChildText("//tr[2]/td[3]/text()"))
		snapshot.USCrashes.Total = parseInt(element.ChildText("//tr[2]/td[4]/text()"))
	},
	tableIdxCanadaInspections: func(element *colly.XMLElement, snapshot *CompanySnapshot) {
		snapshot.CanadaVehicleInspections.Inspections = parseInt(element.ChildText("//tr[2]/td[1]/text()"))
		snapshot.CanadaVehicleInspections.OutOfService = parseInt(element.ChildText("//tr[3]/td[1]/text()"))
		snapshot.CanadaVehicleInspections.OutOfServicePct = parsePctToFloat32(element.ChildText("//tr[4]/td[1]/text()"))
		snapshot.CanadaDriverInspections.Inspections = parseInt(element.ChildText("//tr[2]/td[2]/text()"))
		snapshot.CanadaDriverInspections.OutOfService = parseInt(element.ChildText("//tr[3]/td[2]/text()"))
		snapshot.CanadaDriverInspections.OutOfServicePct = parsePctToFloat32(element.ChildText("//tr[4]/td[2]/text()"))
	},
	tableIdxCanadaCrashes: func(element *colly.XMLElement, snapshot *CompanySnapshot) {
		snapshot.CanadaCrashes.Fatal = parseInt(element.ChildText("//tr[2]/td[1]/text()"))
		snapshot.CanadaCrashes.Injury = parseInt(element.ChildText("//tr[2]/td[2]/text()"))
		snapshot.CanadaCrashes.Tow = parseInt(element.ChildText("//tr[2]/td[3]/text()"))
		snapshot.CanadaCrashes.Total = parseInt(element.ChildText("//tr[2]/td[4]/text()"))
	},
	tableIdxSafetyRating: func(element *colly.XMLElement, snapshot *CompanySnapshot) {
		snapshot.Safety.RatingDate = parseDate(element.ChildText("//tr[2]/td[1]/text()"))
		snapshot.Safety.ReviewDate = parseDate(element.ChildText("//tr[2]/td[2]/text()"))
		snapshot.Safety.Rating = element.ChildText("//tr[3]/td[1]/text()")
		snapshot.Safety.Type = element.ChildText("//tr[3]/td[2]/text()")
	},
}
