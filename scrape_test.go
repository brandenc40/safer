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
	mux.HandleFunc("/snapshot", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write(readTestData("./testdata/snapshot-basic.html"))
	})
	mux.HandleFunc("/snapshot-extras", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write(readTestData("./testdata/snapshot-extras.html"))
	})
	mux.HandleFunc("/snapshot-oos", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write(readTestData("./testdata/snapshot-oos.html"))
	})
	mux.HandleFunc("/snapshot-not-found", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write(readTestData("./testdata/not-found.html"))
	})
	mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})
	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write(readTestData("./testdata/search-result.html"))
	})
	return httptest.NewServer(mux)
}

func readTestData(path string) []byte {
	snapshotHTML, err := os.Open(path)
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

func TestScrapeSnapshot_Basic(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	s := &scraper{
		companySnapshotURL: ts.URL + "/snapshot",
	}
	snapshot, err := s.scrapeCompanySnapshot("", "")
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
		Safety:                   SafetyRating{RatingDate: &ratingDate, ReviewDate: &reviewDate, Rating: "Satisfactory", Type: "Non-Ratable"},
		LatestUpdateDate:         &updateDate,
		OutOfServiceDate:         (*time.Time)(nil),
		MCS150FormDate:           &mcsDate,
		OperationClassification:  []string{"Auth. For Hire"},
		CarrierOperation:         []string{"Interstate"},
		CargoCarried:             []string{"General Freight", "Logs, Poles, Beams, Lumber", "Building Materials", "Fresh Produce", "Intermodal Cont.", "Meat", "Chemicals", "Commodities Dry Bulk", "Refrigerated Food", "Beverages", "Paper Products"},
		LegalName:                "SCHNEIDER NATIONAL CARRIERS INC",
		DBAName:                  "",
		EntityType:               "CARRIER/CARGO TANK/BROKER",
		PhysicalAddress:          "3101 S PACKERLAND DR GREEN BAY, WI 54313",
		Phone:                    "(800) 558-6767",
		MailingAddress:           "PO BOX 2545 GREEN BAY, WI 54306-2545",
		DOTNumber:                "264184",
		StateCarrierID:           "",
		MCMXFFNumbers:            []string{"MC-133655"},
		DUNSNumber:               "15-730-4676",
		MCS150Mileage:            1100158928,
		MCS150Year:               "2020",
		OperatingStatus:          "AUTHORIZED",
		PowerUnits:               10884,
		Drivers:                  12239,
	}
	if !reflect.DeepEqual(expected, snapshot) {
		t.Errorf("scrapeCompanySnapshot() = \n %v, want \n %v", snapshot, expected)
	}
}

func TestScrapeSnapshot_Extras(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	s := &scraper{
		companySnapshotURL: ts.URL + "/snapshot-extras",
	}
	snapshot, err := s.scrapeCompanySnapshot("", "")
	if err != nil {
		t.Errorf("scrapeCompanySnapshot should return no error, but got %v", err)
	}
	if snapshot == nil {
		t.Errorf("snapshot should not return nil")
	}
	updateDate := time.Unix(1629244800, 0).UTC()
	expected := &CompanySnapshot{
		USVehicleInspections:     InspectionSummary{Inspections: 0, OutOfService: 0, OutOfServicePct: 0, NationalAverage: 0.2084},
		USDriverInspections:      InspectionSummary{Inspections: 0, OutOfService: 0, OutOfServicePct: 0, NationalAverage: 0.0545},
		USHazmatInspections:      InspectionSummary{Inspections: 0, OutOfService: 0, OutOfServicePct: 0, NationalAverage: 0.0441},
		USIEPInspections:         InspectionSummary{Inspections: 0, OutOfService: 0, OutOfServicePct: 0, NationalAverage: 0},
		CanadaVehicleInspections: InspectionSummary{Inspections: 0, OutOfService: 0, OutOfServicePct: 0, NationalAverage: 0},
		CanadaDriverInspections:  InspectionSummary{Inspections: 0, OutOfService: 0, OutOfServicePct: 0, NationalAverage: 0},
		USCrashes:                CrashSummary{Fatal: 0, Injury: 0, Tow: 0, Total: 0},
		CanadaCrashes:            CrashSummary{Fatal: 0, Injury: 0, Tow: 0, Total: 0},
		Safety:                   SafetyRating{RatingDate: (*time.Time)(nil), ReviewDate: (*time.Time)(nil), Rating: "None", Type: "None"},
		LatestUpdateDate:         &updateDate,
		OutOfServiceDate:         (*time.Time)(nil),
		MCS150FormDate:           (*time.Time)(nil),
		OperationClassification:  []string{"Private(Property)", "APPLYING F"},
		CarrierOperation:         []string{"Intrastate Only (Non-HM)"},
		CargoCarried:             []string{"Grain, Feed, Hay", "Agricultural/Farm Supplies", "Construction", "ROCK SAND DIRT"},
		LegalName:                "DONALD R SCHNEIDER",
		DBAName:                  "",
		EntityType:               "CARRIER",
		PhysicalAddress:          "230 W BLACKHAWK OLD MONROE, MO 63369",
		Phone:                    "(636) 665-5500",
		MailingAddress:           "230 W BLACKHAWK OLD MONROE, MO 63369",
		DOTNumber:                "884762",
		StateCarrierID:           "",
		MCMXFFNumbers:            []string{},
		DUNSNumber:               "",
		MCS150Mileage:            10000,
		MCS150Year:               "1999",
		OperatingStatus:          "ACTIVE",
		PowerUnits:               1,
		Drivers:                  1,
	}
	if !reflect.DeepEqual(expected, snapshot) {
		t.Errorf("scrapeCompanySnapshot() = \n %#v, want \n %#v", snapshot, expected)
	}
}

func TestScrapeSnapshot_OOS(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	s := &scraper{
		companySnapshotURL: ts.URL + "/snapshot-oos",
	}
	snapshot, err := s.scrapeCompanySnapshot("", "")
	if err != nil {
		t.Errorf("scrapeCompanySnapshot should return no error, but got %v", err)
	}
	if snapshot == nil {
		t.Errorf("snapshot should not return nil")
	}
	oosDate := time.Unix(1019606400, 0).UTC()
	updateDate := time.Unix(1629244800, 0).UTC()
	mcsDate := time.Unix(1088467200, 0).UTC()
	expected := &CompanySnapshot{
		USVehicleInspections:     InspectionSummary{Inspections: 0, OutOfService: 0, OutOfServicePct: 0, NationalAverage: 0.2084},
		USDriverInspections:      InspectionSummary{Inspections: 0, OutOfService: 0, OutOfServicePct: 0, NationalAverage: 0.0545},
		USHazmatInspections:      InspectionSummary{Inspections: 0, OutOfService: 0, OutOfServicePct: 0, NationalAverage: 0.0441},
		USIEPInspections:         InspectionSummary{Inspections: 0, OutOfService: 0, OutOfServicePct: 0, NationalAverage: 0},
		CanadaVehicleInspections: InspectionSummary{Inspections: 0, OutOfService: 0, OutOfServicePct: 0, NationalAverage: 0},
		CanadaDriverInspections:  InspectionSummary{Inspections: 0, OutOfService: 0, OutOfServicePct: 0, NationalAverage: 0},
		USCrashes:                CrashSummary{Fatal: 0, Injury: 0, Tow: 0, Total: 0},
		CanadaCrashes:            CrashSummary{Fatal: 0, Injury: 0, Tow: 0, Total: 0},
		Safety:                   SafetyRating{RatingDate: (*time.Time)(nil), ReviewDate: (*time.Time)(nil), Rating: "", Type: ""},
		LatestUpdateDate:         &updateDate,
		OutOfServiceDate:         &oosDate,
		MCS150FormDate:           &mcsDate,
		OperationClassification:  []string{"Private(Property)"},
		CarrierOperation:         []string{"Intrastate Only (Non-HM)"},
		CargoCarried:             []string{"Machinery, Large Objects", "Construction", "AUGER RIG"},
		LegalName:                "LARRY R RIGGS",
		DBAName:                  "R&J EARTHBORING",
		EntityType:               "CARRIER",
		PhysicalAddress:          "",
		Phone:                    "",
		MailingAddress:           "",
		DOTNumber:                "1003306",
		StateCarrierID:           "",
		MCMXFFNumbers:            []string{},
		DUNSNumber:               "",
		MCS150Mileage:            16000,
		MCS150Year:               "2001",
		OperatingStatus:          "OUT-OF-SERVICE",
		PowerUnits:               2,
		Drivers:                  1,
	}
	if !reflect.DeepEqual(expected, snapshot) {
		t.Errorf("scrapeCompanySnapshot() = \n %#v, want \n %#v", snapshot, expected)
	}
}

func TestScrapeSnapshot_NotFound(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	s := &scraper{
		companySnapshotURL: ts.URL + "/snapshot-not-found",
	}
	snapshot, err := s.scrapeCompanySnapshot("", "")
	if err != ErrCompanyNotFound {
		t.Errorf("scrapeCompanySnapshot should return ErrCompanyNotFound but got %v", err)
	}
	if snapshot != nil {
		t.Errorf("snapshot should return nil")
	}
}

func TestScrapeSnapshot_Error(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	s := &scraper{
		companySnapshotURL: ts.URL + "/error",
	}
	snapshot, err := s.scrapeCompanySnapshot("a", "a")
	if err == nil {
		t.Errorf("scrapeCompanySnapshot should return an error but got %v", err)
	}
	if snapshot != nil {
		t.Errorf("snapshot should return nil but got %v", snapshot)
	}
}

func TestScrapeCompanyNameSearch(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	s := &scraper{
		searchURL: ts.URL + "/search",
	}
	result, err := s.scrapeCompanyNameSearch("")
	if err != nil {
		t.Errorf("scrapeCompanyNameSearch should return no error, but got %v", err)
	}
	if result == nil {
		t.Errorf("result should not return nil")
	}
	expected := []CompanyResult{{Name: "A- SCHNEIDER CONSTRUCTION LLC", DOTNumber: "1261876", Location: "CADOTT, WI"}, {Name: "AARON DUANE SCHNEIDER", DOTNumber: "2563009", Location: "BRODHEAD, WI"}, {Name: "AARON SCHNEIDER", DOTNumber: "3456265", Location: "STANLEY, WI"}, {Name: "ABRAHAM SCHNEIDER", DOTNumber: "2560384", Location: "ROCKFORD, MI"}, {Name: "AL J SCHNEIDER COMPANY", DOTNumber: "123907", Location: "LOUISVILLE, KY"}, {Name: "ALAN SCHNEIDER", DOTNumber: "2153816", Location: "BALDWIN, NY"}, {Name: "ALAN SCHNEIDER TRUCKING COMPANY", DOTNumber: "974012", Location: "SHERIDAN, TX"}, {Name: "ALLAN A SCHNEIDER", DOTNumber: "3649886", Location: "LITTLE FALLS, MN"}, {Name: "ANDREW SCHNEIDER", DOTNumber: "2900237", Location: "GILMAN, WI"}, {Name: "ANTHONY Q SCHNEIDER", DOTNumber: "2589751", Location: "ASHLEY, ND"}, {Name: "ARLO SCHNEIDER", DOTNumber: "2189204", Location: "STAPLETON, GA"}, {Name: "BARRY L SCHNEIDER JR", DOTNumber: "646237", Location: "HENDERSON, KY"}, {Name: "BART SCHNEIDERMAN", DOTNumber: "1463475", Location: "WEST BURLINGTON, IA"}, {Name: "BERGSCHNEIDER FARMS", DOTNumber: "2055571", Location: "WAVERLY, IL"}, {Name: "BILL SCHNEIDER LANDSCAPING", DOTNumber: "2901731", Location: "E NORTHPORT, NY"}, {Name: "BILL SCHNEIDER TRUCKING", DOTNumber: "719325", Location: "NORTH VERNON, IN"}, {Name: "BRADLEY BRETTSCHNEIDER", DOTNumber: "3011702", Location: "IDA, MI"}, {Name: "BRADLEY SCHNEIDER", DOTNumber: "3295014", Location: "FL RIVER MLS, CA"}, {Name: "BRENT SCHNEIDER", DOTNumber: "2200357", Location: "GRIFTON, NC"}, {Name: "BRENT SCHNEIDER", DOTNumber: "2094550", Location: "EL CAMPO, TX"}, {Name: "BRETSCHNEIDER CO", DOTNumber: "1130627", Location: "HASTINGS, MN"}, {Name: "BRIAN CURTIS SCHNEIDER", DOTNumber: "2411806", Location: "GANADO, TX"}, {Name: "BRIAN M SCHNEIDER", DOTNumber: "2827619", Location: "GREENBAY, WI"}, {Name: "BRIAN SCHNEIDER", DOTNumber: "1584782", Location: "CHILTON, WI"}, {Name: "BRUCE A SCHNEIDER", DOTNumber: "1406168", Location: "JANESVILLE, WI"}, {Name: "BRUCE E SCHNEIDER", DOTNumber: "2686007", Location: "PETALUMA, CA"}, {Name: "BRUCE SCHNEIDER", DOTNumber: "2503593", Location: "MANCHESTER, IA"}, {Name: "C J SCHNEIDER TRUCKING INC", DOTNumber: "993997", Location: "CHILTON, WI"}, {Name: "CARL SCHNEIDER", DOTNumber: "638744", Location: "MADISON, CT"}, {Name: "CARL SCHNEIDER LOGGING", DOTNumber: "2029598", Location: "WORLAND, WY"}, {Name: "CASEY M SCHNEIDER", DOTNumber: "1359358", Location: "CHILTON, WI"}, {Name: "CASEY SCHNEIDER", DOTNumber: "3229487", Location: "RICHMOND, KY"}, {Name: "CHARLES C RIEMENSCHNEIDER", DOTNumber: "769523", Location: "FAIRLESS HILLS, PA"}, {Name: "CHARLES H SCHNEIDER", DOTNumber: "2768702", Location: "CHINO, CA"}, {Name: "CHARLES SCHNEIDER", DOTNumber: "2079873", Location: "AUGUSTA, KS"}, {Name: "CHARLES SCHNEIDER CONSTRUCTION CORP", DOTNumber: "999699", Location: "HAZELHURST, WI"}, {Name: "CHARLES SCHNEIDER'S SERVICES", DOTNumber: "1128542", Location: "NORFOLK, NE"}, {Name: "CHRIS SCHNEIDER", DOTNumber: "1949319", Location: "FLANAGAN, IL"}, {Name: "CHRIS SCHNEIDER", DOTNumber: "2051868", Location: "TAYLOR, TX"}, {Name: "CHUCK SCHNEIDER TRUCKING LLC", DOTNumber: "3212697", Location: "GENOA CITY, WI"}, {Name: "CLARK SCHNEIDER", DOTNumber: "1554563", Location: "PERRINTON, MI"}, {Name: "CLIFTON SCHNEIDER", DOTNumber: "2201375", Location: "EVART, MI"}, {Name: "CRAIG SCHNEIDER", DOTNumber: "2468479", Location: "HOLLAND, TX"}, {Name: "CURTIS SCHNEIDER", DOTNumber: "2454459", Location: "ALAMOSA, CO"}, {Name: "D SCHNEIDER TRUCKING LLC", DOTNumber: "3654223", Location: "WAPELLO, IA"}, {Name: "DALE E SCHNEIDER", DOTNumber: "2004211", Location: "CHARTER OAK, IA"}, {Name: "DANIEL L SCHNEIDER", DOTNumber: "700698", Location: "KENNARD, NE"}, {Name: "DANIEL L SCHNEIDER", DOTNumber: "1649526", Location: "ONTARIO, NY"}, {Name: "DANIEL SCHNEIDER", DOTNumber: "3030698", Location: "SAINT LOUIS, MO"}, {Name: "DANIEL SCHNEIDER", DOTNumber: "1402912", Location: "CASSVILLE, WI"}, {Name: "DARVID A SCHNEIDER & ALBERT A COBLENTZ", DOTNumber: "701326", Location: "BRECKENRIDGE, MI"}, {Name: "DAVID LLOYD REIFSCHNEIDER", DOTNumber: "2104236", Location: "HUBBARD, IA"}, {Name: "DAVID N SCHNEIDER", DOTNumber: "1994662", Location: "WHEATLAND, IA"}, {Name: "DAVID SCHNEIDER", DOTNumber: "3388590", Location: "BEL AIR, MD"}, {Name: "DAVID SCHNEIDER", DOTNumber: "1205234", Location: "SARASOTA, FL"}, {Name: "DAVID SCHNEIDER", DOTNumber: "1539096", Location: "FOND DU LAC, WI"}, {Name: "DAVID SCHNEIDER FARMS", DOTNumber: "1467712", Location: "MELVILLE, NY"}, {Name: "DAVID SCHNEIDERHAN", DOTNumber: "1556534", Location: "MALONE, WI"}, {Name: "DAVID SCHNEIDERMANN", DOTNumber: "1211666", Location: "ULEN, MN"}, {Name: "DAVID W SCHNEIDER", DOTNumber: "2375131", Location: "UPPERCO, MD"}, {Name: "DAVID W SCHNEIDER", DOTNumber: "2368265", Location: "RHINELANDER, WI"}, {Name: "DENNIS RIEMENSCHNEIDER", DOTNumber: "1369384", Location: "RIVER FALLS, WI"}, {Name: "DEVIN BERGSCHNEIDER TRUCKING INC", DOTNumber: "2888036", Location: "FRANKLIN, IL"}, {Name: "DONALD R SCHNEIDER", DOTNumber: "884762", Location: "OLD MONROE, MO"}, {Name: "DONALD W SCHATTSCHNEIDER", DOTNumber: "524150", Location: "PHILLIPS, ME"}, {Name: "DOUG BERGSCHNEIDER", DOTNumber: "1842500", Location: "GREENFIELD, IL"}, {Name: "DOUG BERGSCHNEIDER TRUCKING", DOTNumber: "1842500", Location: "GREENFIELD, IL"}, {Name: "DR SCHNEIDER AUTOMOTIVE SYSTEMS INC", DOTNumber: "2986460", Location: "RUSSELL SPGS, KY"}, {Name: "DR SCHNEIDER AUTOMOTIVE SYSTEMS LLC", DOTNumber: "1855404", Location: "BRIGHTON, MI"}, {Name: "DW SCHNEIDER ENGINEERING LLC", DOTNumber: "2918924", Location: "WALDO, WI"}, {Name: "E SCHNEIDER DISTRIBUTION", DOTNumber: "3045414", Location: "JEFFERSONVLLE, IN"}, {Name: "E SCHNEIDER ENTERPRISES", DOTNumber: "1236560", Location: "TAYLORS, SC"}, {Name: "E SCHNEIDER SONS INC", DOTNumber: "10478", Location: "ALLENTOWN, PA"}, {Name: "EDWARD A SCHNEIDER", DOTNumber: "1062163", Location: "KANOPOLIS, KS"}, {Name: "EDWARD SCHNEIDER", DOTNumber: "1290469", Location: "VINTON, TX"}, {Name: "ELIZABETH SCHNEIDER", DOTNumber: "1511796", Location: "RICE LAKE, WI"}, {Name: "ELROY SCHNEIDER", DOTNumber: "1429478", Location: "CHILTON, WI"}, {Name: "EMJ SCHNEIDER FARMS", DOTNumber: "1555374", Location: "OAKLEY, MI"}, {Name: "ERIC D SCHNEIDER", DOTNumber: "3240096", Location: "LODI, CA"}, {Name: "ERIC SCHNEIDER", DOTNumber: "2139871", Location: "SAINT JOHNS, MI"}, {Name: "ERIC SCHNEIDER", DOTNumber: "2519968", Location: "PARADISE, MT"}, {Name: "ERIC SCHNEIDER LANDSCAPES LLC", DOTNumber: "3319159", Location: "SPRING CITY, PA"}, {Name: "EUGENE MARTY JOE SCHNEIDER", DOTNumber: "1555374", Location: "OAKLEY, MI"}, {Name: "EVAN SCHNEIDER", DOTNumber: "3445423", Location: "ROCKLAND, MA"}, {Name: "EWALD I SCHNEIDER TRUCKING LLC", DOTNumber: "1942803", Location: "OTTO, TX"}, {Name: "F SCHNEIDER CONTRACTING CORP", DOTNumber: "3306200", Location: "HAMPTON BAYS, NY"}, {Name: "FRANCIS D BRESCHNEIDER", DOTNumber: "1130627", Location: "HASTINGS, MN"}, {Name: "FRED SCHNEIDER", DOTNumber: "1667853", Location: "MIDDLETON, MI"}, {Name: "FRED SCHNEIDER FARM", DOTNumber: "1667853", Location: "MIDDLETON, MI"}, {Name: "FREDERICK D SCHNEIDER", DOTNumber: "2340303", Location: "MCDONALD, PA"}, {Name: "FREDERICK JOHN SCHNEIDER", DOTNumber: "1136590", Location: "COLUMBIA HEIGHTS, MN"}, {Name: "G SCHNEIDER FARMS LLC", DOTNumber: "1400363", Location: "KIEL, WI"}, {Name: "GARY E SCHNEIDER", DOTNumber: "1524672", Location: "CASSADAGA, NY"}, {Name: "GARY P SCHNEIDER", DOTNumber: "1846913", Location: "MACOM, MI"}, {Name: "GENE SCHNEIDER EXCAVATING", DOTNumber: "1047523", Location: "MARSHFIELD, WI"}, {Name: "GENE STUCKENSCHNEIDER TRUCKING LLC", DOTNumber: "554633", Location: "MARTINSBURG, MO"}, {Name: "GENE W SCHNEIDER", DOTNumber: "1047523", Location: "MARSHFIELD, WI"}, {Name: "GEOFF SCHNEIDER", DOTNumber: "1188943", Location: "WALTON, NE"}, {Name: "GEOFFREY G SCHNEIDER", DOTNumber: "2048115", Location: "CLERMONT, IA"}, {Name: "GEORGE S SCHNEIDER JR", DOTNumber: "2082571", Location: "FRANKLINVILLE, NY"}, {Name: "GEORGE SCHNEIDER CONTRACTING", DOTNumber: "1732092", Location: "BEAVERTON, MI"}, {Name: "GEORGE SCHNEIDER JR", DOTNumber: "2517872", Location: "SAINT JAMES, NY"}, {Name: "GERALD SCHNEIDER", DOTNumber: "1375403", Location: "CECIL, WI"}, {Name: "GERD K SCHNEIDER", DOTNumber: "2605159", Location: "GILROY, CA"}, {Name: "GERD SCHNEIDER NURSERY", DOTNumber: "2605159", Location: "GILROY, CA"}, {Name: "GERD SCHNEIDER NURSERY LLC", DOTNumber: "2886077", Location: "GILROY, CA"}, {Name: "GERRIT J SCHNEIDER", DOTNumber: "872418", Location: "ALAMOSA, CO"}, {Name: "GLENN SCHNEIDER", DOTNumber: "2063239", Location: "BELTON, TX"}, {Name: "GORDON SCHNEIDER", DOTNumber: "1567319", Location: "CHILTON, WI"}, {Name: "GRANT S SCHNEIDER", DOTNumber: "3425516", Location: "ALBANY, NY"}, {Name: "GREG SCHNEIDER ELECTRIC INC", DOTNumber: "1535120", Location: "REMINGTON, IN"}, {Name: "GUSTAV A SCHNEIDER IV", DOTNumber: "2594105", Location: "NEKOOSA, WI"}, {Name: "HANS J SCHNEIDER", DOTNumber: "2698396", Location: "HAYWARD, CA"}, {Name: "HAROLD SCHNEIDER", DOTNumber: "2004853", Location: "SAINT CLAIR, MI"}, {Name: "HARRY SCHNEIDER LEASING LLC", DOTNumber: "2225041", Location: "LAKE ST. LOUIS, MO"}, {Name: "HARVEY M SCHNEIDER", DOTNumber: "2181827", Location: "CHILTON, WI"}, {Name: "HARVEY N SCHNEIDER", DOTNumber: "1414510", Location: "REEDSBURG, WI"}, {Name: "HOFFSCHNEIDER CONSTRUCTION LLC", DOTNumber: "2795992", Location: "BLAIR, NE"}, {Name: "J & K POHLSCHNEIDER INC", DOTNumber: "3563804", Location: "SAINT PAUL, OR"}, {Name: "J ALLAN SCHNEIDER CORP", DOTNumber: "2651675", Location: "CALISTOGA, CA"}, {Name: "J R SCHNEIDER CONSTRUCTION INC", DOTNumber: "1932441", Location: "AUSTIN, TX"}, {Name: "J SCHNEIDER GROUP INC", DOTNumber: "2643904", Location: "CEDARHURST, NY"}, {Name: "J SCHNEIDER TRANSPORT LLC", DOTNumber: "2310000", Location: "NAPPANEE, IN"}, {Name: "J.A. SCHNEIDER INC.", DOTNumber: "2606154", Location: "MONTEBELLO, CA"}, {Name: "JACOB J SCHNEIDER", DOTNumber: "2060648", Location: "RUTHVEN, IA"}, {Name: "JACOB SCHNEIDER", DOTNumber: "3230603", Location: "NYA, MN"}, {Name: "JAKE SCHNEIDER", DOTNumber: "2060648", Location: "RUTHVEN, IA"}, {Name: "JAMES A BERGSCHNEIDER", DOTNumber: "1845959", Location: "SOUTH DAYTONA, FL"}, {Name: "JAMES DIRKSCHNEIDER", DOTNumber: "1896618", Location: "MEDFORD, NY"}, {Name: "JAMES E SCHNEIDER", DOTNumber: "1006221", Location: "KIEL, WI"}, {Name: "JAMES P SCHNEIDER AND ASSOCIATES", DOTNumber: "2157015", Location: "AUBURN HILLS, MI"}, {Name: "JAMES S SCHNEIDER", DOTNumber: "2728757", Location: "MOUNTAIN HOUSE, CA"}, {Name: "JAMES SCHNEIDER", DOTNumber: "1126759", Location: "SCANDIA, MN"}, {Name: "JAMES SCHNEIDER", DOTNumber: "1571972", Location: "CENTRAL CITY, NE"}, {Name: "JAMES WILLIAM SCHNEIDER", DOTNumber: "2343284", Location: "FOREST CITY, PA"}, {Name: "JAMIE P SCHNEIDER", DOTNumber: "1678114", Location: "TERRE HAUTE, IN"}, {Name: "JARROD SCHNEIDER", DOTNumber: "2435721", Location: "WALNUT HILL, FL"}, {Name: "JASON SCHNEIDER", DOTNumber: "2276267", Location: "FORESTVILLE, NY"}, {Name: "JEFFREY & SCOT BERGSCHNEIDER", DOTNumber: "2055571", Location: "WAVERLY, IL"}, {Name: "JEFFREY L SCHNEIDER", DOTNumber: "1412101", Location: "ST PAUL, IN"}, {Name: "JEFFREY THOMAS RIEMENSCHNEIDER CATERING", DOTNumber: "1122392", Location: "VADNAIS HEIGHTS, MN"}, {Name: "JEROME SCHNEIDER", DOTNumber: "1589212", Location: "BISMARCK, ND"}, {Name: "JERRY JOHN SCHNEIDER JR", DOTNumber: "1408536", Location: "CRAWFORDVILLE, FL"}, {Name: "JERRY R SCHNEIDER", DOTNumber: "2691170", Location: "CAMARILLO, CA"}, {Name: "JERRY SCHNEIDER", DOTNumber: "1365230", Location: "SPARTA, WI"}, {Name: "JOAN M SCHNEIDER LLC", DOTNumber: "1540788", Location: "EVANSVILLE, IN"}, {Name: "JODY LEE SCHNEIDER", DOTNumber: "2297496", Location: "SEYMOUR, WI"}, {Name: "JOEL A SCHNEIDER", DOTNumber: "1001708", Location: "FREDONIA, WI"}, {Name: "JOHN C SCHNEIDER", DOTNumber: "1472585", Location: "CECIL, WI"}, {Name: "JOHN D SCHNEIDER", DOTNumber: "411381", Location: "MOUNT VERNON, IN"}, {Name: "JOHN E SCHNEIDER", DOTNumber: "1389845", Location: "GLADBROOK, IA"}, {Name: "JOHN E SCHNEIDER CONSTRUCTION LLC", DOTNumber: "3040564", Location: "VANCOUVER, WA"}, {Name: "JOHN M SCHNEIDER", DOTNumber: "1921832", Location: "CHATTAROY, WA"}, {Name: "JOHN P REIFSCHNEIDER", DOTNumber: "1983585", Location: "LINCOLN, NE"}, {Name: "JOHN SCHNEIDER", DOTNumber: "2311339", Location: "WATERLOO, IL"}, {Name: "JOHN SCHNEIDER", DOTNumber: "2356015", Location: "BURNSVILLE, MN"}, {Name: "JOHN SCHNEIDER", DOTNumber: "1950691", Location: "CHESANING, MI"}, {Name: "JON SCHNEIDER", DOTNumber: "1706697", Location: "EAGLE LAKE, MN"}, {Name: "JOSEPH A SCHNEIDER", DOTNumber: "3372049", Location: "S EGREMONT, MA"}, {Name: "JOSEPH MICHAEL SCHNEIDER", DOTNumber: "1844651", Location: "CAPE GIRARDEAU, MO"}, {Name: "JOSEPH SCHNEIDER", DOTNumber: "2076552", Location: "GEORGETOWN, KY"}, {Name: "JOSEPH SCHNEIDER", DOTNumber: "3021608", Location: "HIGHLAND PARK, NJ"}, {Name: "JT SCHNEIDER TRUCKING LLC", DOTNumber: "2461918", Location: "SILEX, MO"}}
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("scrapeCompanyNameSearch() = \n %#v \n , want \n %v", result, expected)
	}
}

func TestScrapeCompanyNameSearch_Error(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	s := &scraper{
		searchURL: ts.URL + "/error",
	}
	result, err := s.scrapeCompanyNameSearch("")
	if err == nil {
		t.Errorf("scrapeCompanyNameSearch should return an error but got %v", err)
	}
	if result != nil {
		t.Errorf("result should return nil")
	}
}
