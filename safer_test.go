package safer

import (
	"testing"

	"github.com/antchfx/htmlquery"
)

func TestNewClient(t *testing.T) {
	c := NewClient()
	if c == nil {
		t.Error("expected client to not be nil but it was")
	}
}

func TestClient_GetCompanyByDOTNumber(t *testing.T) {
	s := newTestServer()
	defer s.Close()

	c := &Client{
		scraper: scraper{
			companySnapshotURL: s.URL + "/snapshot",
			searchURL:          s.URL + "/search",
		},
	}
	snapshot, err := c.GetCompanyByDOTNumber("")
	if snapshot == nil {
		t.Error("snapshot returned nil")
	}
	if err != nil {
		t.Errorf("error expected nil but got %v", err)
	}

	snapshot, err = c.GetCompanyByMCMX("")
	if snapshot == nil {
		t.Error("snapshot returned nil")
	}
	if err != nil {
		t.Errorf("error expected nil but got %v", err)
	}

	results, err := c.SearchCompaniesByName("")
	if results == nil {
		t.Error("results returned nil")
	}
	if len(results) == 0 {
		t.Error("results length = 0")
	}
	if err != nil {
		t.Errorf("error expected nil but got %v", err)
	}
}

func BenchmarkClient_GetCompanyByDOTNumber(b *testing.B) {
	node, _ := htmlquery.LoadDoc("./testdata/snapshot-basic.html")
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = htmlNodeToCompanySnapshot(node)
	}
}

func BenchmarkClient_Search_4Results(b *testing.B) {
	node, _ := htmlquery.LoadDoc("./testdata/search-result-short.html")
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = htmlNodeToCompanyResults(node)
	}
}

func BenchmarkClient_Search_484Results(b *testing.B) {
	node, _ := htmlquery.LoadDoc("./testdata/search-result.html")
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = htmlNodeToCompanyResults(node)
	}
}
