package safer

import (
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"
)

func newTestServer() *httptest.Server {
	mux := http.NewServeMux()
	data := readSnapshotTestData()
	mux.HandleFunc("/snapshot", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write(data)
	})
	return httptest.NewServer(mux)
}

func readSnapshotTestData() []byte {
	snapshotHTML, err := os.Open("testdata/snapshot.html")
	defer snapshotHTML.Close()
	if err != nil {
		panic(err)
	}
	data := make([]byte, 1024*60)
	_, err = snapshotHTML.Read(data)
	if err != nil {
		panic(err)
	}
	return data
}

func TestScrapeSnapshot(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	snapshot, err := scrapeCompanySnapshot("", "", ts.URL+"/snapshot")
	if err != nil {
		t.Errorf("scrapeCompanySnapshot should return no error, but got %v", err)
	}
	if snapshot == nil {
		t.Errorf("snapshot should not return nil")
	}
	ratingDate := time.Unix(1045699200, 0).UTC()
	reviewDate := time.Unix(1602633600, 0).UTC()
	updateDate := time.Unix(1628899200, 0).UTC()
	mcsDate := time.Unix(1618790400, 0).UTC()
	expected := &CompanySnapshot{
		USVehicleInspections:     InspectionSummary{Inspections: 7276, OutOfService: 991, OutOfServicePct: 0.136, NationalAverage: 0.2084},
		USDriverInspections:      InspectionSummary{Inspections: 13728, OutOfService: 71, OutOfServicePct: 0.005, NationalAverage: 0.0545},
		USHazmatInspections:      InspectionSummary{Inspections: 426, OutOfService: 6, OutOfServicePct: 0.014, NationalAverage: 0.0441},
		USIEPInspections:         InspectionSummary{Inspections: 2, OutOfService: 0, OutOfServicePct: 0, NationalAverage: 0},
		CanadaVehicleInspections: InspectionSummary{Inspections: 24, OutOfService: 8, OutOfServicePct: 0.333, NationalAverage: 0},
		CanadaDriverInspections:  InspectionSummary{Inspections: 30, OutOfService: 8, OutOfServicePct: 0.267, NationalAverage: 0},
		USCrashes:                CrashSummary{Fatal: 15, Injury: 248, Tow: 574, Total: 837},
		CanadaCrashes:            CrashSummary{Fatal: 0, Injury: 0, Tow: 1, Total: 1},
		Safety: SafetyRating{
			RatingDate: &ratingDate,
			ReviewDate: &reviewDate,
			Rating:     "Satisfactory",
			Type:       "Non-Ratable",
		},
		LatestUpdateDate:        &updateDate,
		OutOfServiceDate:        (*time.Time)(nil),
		MCS150FormDate:          &mcsDate,
		OperationClassification: []string{"Auth. For Hire"},
		CarrierOperation:        []string{"Interstate"},
		CargoCarried:            []string{"General Freight", "Logs, Poles, Beams, Lumber", "Building Materials", "Fresh Produce", "Intermodal Cont.", "Meat", "Chemicals", "Commodities Dry Bulk", "Refrigerated Food", "Beverages", "Paper Products"},
		LegalName:               "SCHNEIDER NATIONAL CARRIERS INC",
		DBAName:                 "",
		EntityType:              "CARRIER/CARGO TANK/BROKER",
		PhysicalAddress:         "3101 S PACKERLAND DR GREEN BAY, WI 54313",
		Phone:                   "(800) 558-6767",
		MailingAddress:          "PO BOX 2545 GREEN BAY, WI 54306-2545",
		DOTNumber:               "264184",
		StateCarrierID:          "",
		MCMXFFNumbers:           "MC-133655",
		DUNSNumber:              "15-730-4676",
		MCS150Mileage:           1100158928,
		MCS150Year:              "2020",
		OperatingStatus:         "AUTHORIZED",
		PowerUnits:              10884,
		Drivers:                 12239,
	}
	if !reflect.DeepEqual(expected, snapshot) {
		t.Errorf("scrapeCompanySnapshot() = \n %v, want \n %v", snapshot, expected)
	}
}
