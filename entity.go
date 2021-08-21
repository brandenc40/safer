package safer

import "time"

// CompanyResult is the search result returned from a company query by name
type CompanyResult struct {
	Name      string `json:"name"`
	DOTNumber string `json:"dot_number"`
	Location  string `json:"location"`
}

// CompanySnapshot data parsed from the https://safer.fmcsa.dot.gov/CompanySnapshot.aspx website
type CompanySnapshot struct {
	USVehicleInspections     InspectionSummary `json:"us_vehicle_inspections"`
	USDriverInspections      InspectionSummary `json:"us_driver_inspections"`
	USHazmatInspections      InspectionSummary `json:"us_hazmat_inspections"`
	USIEPInspections         InspectionSummary `json:"us_iep_inspections"`
	CanadaVehicleInspections InspectionSummary `json:"canada_vehicle_inspections"`
	CanadaDriverInspections  InspectionSummary `json:"canada_driver_inspections"`
	USCrashes                CrashSummary      `json:"us_crashes"`
	CanadaCrashes            CrashSummary      `json:"canada_crashes"`
	Safety                   SafetyRating      `json:"safety"`
	LatestUpdateDate         *time.Time        `json:"latest_update_date"`
	OutOfServiceDate         *time.Time        `json:"out_of_service_date"`
	MCS150FormDate           *time.Time        `json:"mcs_150_form_date"`
	OperationClassification  []string          `json:"operation_classification"`
	CarrierOperation         []string          `json:"carrier_operation"`
	CargoCarried             []string          `json:"cargo_carried"`
	LegalName                string            `json:"legal_name"`
	DBAName                  string            `json:"dba_name"`
	EntityType               string            `json:"entity_type"`
	PhysicalAddress          string            `json:"physical_address"`
	Phone                    string            `json:"phone"`
	MailingAddress           string            `json:"mailing_address"`
	DOTNumber                string            `json:"dot_number"`
	StateCarrierID           string            `json:"state_carrier_id"`
	MCMXFFNumbers            []string          `json:"mc_mx_ff_numbers"`
	DUNSNumber               string            `json:"duns_number"`
	MCS150Mileage            int               `json:"mcs_150_mileage"`
	MCS150Year               string            `json:"mcs_150_year"`
	OperatingStatus          string            `json:"operating_status"`
	PowerUnits               int               `json:"power_units"`
	Drivers                  int               `json:"drivers"`
}

// InspectionSummary for 24 months prior to LatestUpdateDate.
//
// Note: NationalAverage not available for Canadian summaries
type InspectionSummary struct {
	Inspections     int     `json:"inspections"`
	OutOfService    int     `json:"out_of_service"`
	OutOfServicePct float32 `json:"out_of_service_pct"`
	NationalAverage float32 `json:"national_average"`
}

// CrashSummary for 24 months prior to LatestUpdateDate
type CrashSummary struct {
	Fatal  int `json:"fatal"`
	Injury int `json:"injury"`
	Tow    int `json:"tow"`
	Total  int `json:"total"`
}

// SafetyRating current as of LatestUpdateDate
type SafetyRating struct {
	RatingDate *time.Time `json:"rating_date"`
	ReviewDate *time.Time `json:"review_date"`
	Rating     string     `json:"rating"`
	Type       string     `json:"type"`
}
