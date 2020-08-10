package podcast_test

import (
	"context"
	"testing"
	"time"

	"github.com/dapperAuteur/dashboard-go-api/internal/podcast"
	"github.com/dapperAuteur/dashboard-go-api/internal/schema"
	"github.com/dapperAuteur/dashboard-go-api/internal/tests"
	"github.com/google/go-cmp/cmp"
)

func TestPodcasts(t *testing.T) {
	db, teardown := tests.NewUnit(t)
	defer teardown()

	ctx := context.Background()

	newP := podcast.NewPodcast{
		Title:     "Comic Book",
		Author:     "10",
		Tags: ["55","00"],
	}
	now := time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)

	p0, err := podcast.CreatePodcast(ctx, db, newP, now)
	if err != nil {
		t.Fatalf("creating podcast p0: %s", err)
	}

	p1, err := podcast.Retrieve(ctx, db, p0.ID)
	if err != nil {
		t.Fatalf("getting podcast p0: %s", err)
	}

	if diff := cmp.Diff(p1, p0); diff != "" {
		t.Fatalf("fetched != created:\n%s", diff)
	}
}

func TestPodcastList(t *testing.T) {
	db, teardown := tests.NewUnit(t)
	defer teardown()

	ctx := context.Background()

	if err := schema.Seed(ctx, db); err != nil {
		t.Fatal(err)
	}

	ps, err := podcast.List(ctx, db)
	if err != nil {
		t.Fatalf("listing products: %s", err)
	}
	if exp, got := 2, len(ps); exp != got {
		t.Fatalf("expected podcast list size %v, got %v", exp, got)
	}
}
