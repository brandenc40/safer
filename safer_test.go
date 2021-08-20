package safer

import (
	"testing"

	"github.com/gocolly/colly/v2"
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

	c := (Client)(&client{
		scraper: &scraper{
			baseCollector:      colly.NewCollector(),
			companySnapshotURL: s.URL + "/snapshot",
			searchURL:          s.URL + "/search",
		},
	})
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
	if results == nil || len(results) == 0 {
		t.Error("snapshot returned nil")
	}
	if err != nil {
		t.Errorf("error expected nil but got %v", err)
	}
}