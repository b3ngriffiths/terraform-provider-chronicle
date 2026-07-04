package client

import "testing"

func TestGenerateDefaultBasePaths(t *testing.T) {
	cases := []struct {
		region   string
		wantFeed string
	}{
		{"us", "https://backstory.googleapis.com/v1/feeds"},
		{"europe", "https://europe-backstory.googleapis.com/v1/feeds"},
		{"europe-west3", "https://europe-west3-backstory.googleapis.com/v1/feeds"},
		{"me-central2", "https://me-central2-backstory.googleapis.com/v1/feeds"},
	}

	for _, tc := range cases {
		paths := GenerateDefaultBasePaths(tc.region)
		if got := paths[FeedManagementBasePathKey]; got != tc.wantFeed {
			t.Errorf("region %s: feed base path = %s, want %s", tc.region, got, tc.wantFeed)
		}
	}
}

func TestAllRegionsProduceValidBasePaths(t *testing.T) {
	for _, region := range Regions {
		paths := GenerateDefaultBasePaths(region)
		for key, path := range paths {
			if path == "" {
				t.Errorf("region %s: empty base path for %s", region, key)
			}
		}
	}
}
