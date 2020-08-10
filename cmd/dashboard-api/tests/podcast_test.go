package tests

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	// NOTE: Models should not be imported, we want to test the exact JSON. We
	// make the comparison process easier using the go-cmp library.
	"github.com/dapperAuteur/dashboard-go-api/cmd/dashboard-api/internal/handlers"
	"github.com/dapperAuteur/dashboard-go-api/internal/schema"
	"github.com/dapperAuteur/dashboard-go-api/internal/tests"
	"github.com/google/go-cmp/cmp"
)

// TestPodcasts runs a series of tests to exercise Podcast behavior from the
// API level. The subtests all share the same database and application for
// speed and convenience. The downside is the order the tests are ran matters
// and one test may break if other tests are not ran before it. If a particular
// subtest needs a fresh instance of the application it can make it or it
// should be its own Test* function.
func TestPodcasts(t *testing.T) {
	db, teardown := tests.NewUnit(t)
	defer teardown()

	if err := schema.Seed(db); err != nil {
		t.Fatal(err)
	}

	log := log.New(os.Stderr, "TEST : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	tests := PodcastTests{app: handlers.API(db, log)}

	t.Run("List", tests.List)
	t.Run("PodcastCRUD", tests.PodcastCRUD)
}

// PodcastTests holds methods for each podcast subtest. This type allows
// passing dependencies for tests while still providing a convenient syntax
// when subtests are registered.
type PodcastTests struct {
	app http.Handler
}

func (p *PodcastTests) List(t *testing.T) {
	req := httptest.NewRequest("GET", "/v1/podcasts", nil)
	resp := httptest.NewRecorder()

	p.app.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("getting: expected status code %v, got %v", http.StatusOK, resp.Code)
	}

	var list []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		t.Fatalf("decoding: %s", err)
	}

	want := []map[string]interface{}{
		{
			"id":        "a2b0639f-2cc6-44b8-b97b-15d69dbb511e",
			"name":      "Comic Books",
			"cost":      float64(50),
			"quantity":  float64(42),
			"createdAt": "2019-01-01T00:00:01.000001Z",
			"updatedAt": "2019-01-01T00:00:01.000001Z",
		},
		{
			"id":        "72f8b983-3eb4-48db-9ed0-e45cc6bd716b",
			"name":      "McDonalds Toys",
			"cost":      float64(75),
			"quantity":  float64(120),
			"createdAt": "2019-01-01T00:00:02.000001Z",
			"updatedAt": "2019-01-01T00:00:02.000001Z",
		},
	}

	if diff := cmp.Diff(want, list); diff != "" {
		t.Fatalf("Response did not match expected. Diff:\n%s", diff)
	}
}

func (p *PodcastTests) PodcastCRUD(t *testing.T) {
	var created map[string]interface{}

	{ // CREATE
		body := strings.NewReader(`{"name":"podcast0","cost":55,"quantity":6}`)

		req := httptest.NewRequest("POST", "/v1/podcasts", body)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		p.app.ServeHTTP(resp, req)

		if http.StatusCreated != resp.Code {
			t.Fatalf("posting: expected status code %v, got %v", http.StatusCreated, resp.Code)
		}

		if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
			t.Fatalf("decoding: %s", err)
		}

		if created["id"] == "" || created["id"] == nil {
			t.Fatal("expected non-empty podcast id")
		}
		if created["createdAt"] == "" || created["createdAt"] == nil {
			t.Fatal("expected non-empty podcast createdAt")
		}
		if created["updatedAt"] == "" || created["updatedAt"] == nil {
			t.Fatal("expected non-empty podcast updatedAt")
		}

		want := map[string]interface{}{
			"id":        created["id"],
			"createdAt": created["createdAt"],
			"updatedAt": created["updatedAt"],
			"name":      "podcast0",
			"cost":      float64(55),
			"quantity":  float64(6),
		}

		if diff := cmp.Diff(want, created); diff != "" {
			t.Fatalf("Response did not match expected. Diff:\n%s", diff)
		}
	}

	{ // READ
		url := fmt.Sprintf("/v1/podcasts/%s", created["id"])
		req := httptest.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		p.app.ServeHTTP(resp, req)

		if http.StatusOK != resp.Code {
			t.Fatalf("retrieving: expected status code %v, got %v", http.StatusOK, resp.Code)
		}

		var fetched map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&fetched); err != nil {
			t.Fatalf("decoding: %s", err)
		}

		// Fetched podcast should match the one we created.
		if diff := cmp.Diff(created, fetched); diff != "" {
			t.Fatalf("Retrieved podcast should match created. Diff:\n%s", diff)
		}
	}
}
