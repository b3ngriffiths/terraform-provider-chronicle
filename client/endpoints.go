package client

import (
	"fmt"
	"time"

	"golang.org/x/time/rate"
)

type ClientRateLimiters struct {
	FeedManagementCreateFeed *rate.Limiter
	FeedManagementGetFeed    *rate.Limiter
	FeedManagementListFeeds  *rate.Limiter
	FeedManagementUpdateFeed *rate.Limiter
	FeedManagementDeleteFeed *rate.Limiter
	FeedManagementEnableFeed *rate.Limiter

	DetectionCreateRule         *rate.Limiter
	DetectionCreateRuleVersion  *rate.Limiter
	DetectionGetRule            *rate.Limiter
	DetectionUpdateRule         *rate.Limiter
	DetectionDeleteRule         *rate.Limiter
	DetectionEnableLiveRule     *rate.Limiter
	DetectionEnableAlertingRule *rate.Limiter
	DetectionVerifyYARARule     *rate.Limiter

	RBACCreateSubject *rate.Limiter
	RBACGetSubject    *rate.Limiter
	RBACUpdateSubject *rate.Limiter
	RBACDeleteSubject *rate.Limiter

	ReferenceListsCreateList *rate.Limiter
	ReferenceListsGetList    *rate.Limiter
	ReferenceListsUpdateList *rate.Limiter
}

func NewClientRateLimiters() *ClientRateLimiters {
	return &ClientRateLimiters{
		FeedManagementCreateFeed: rate.NewLimiter(rate.Every(time.Second), 1),
		FeedManagementGetFeed:    rate.NewLimiter(rate.Every(time.Second), 1),
		FeedManagementListFeeds:  rate.NewLimiter(rate.Every(time.Second), 1),
		FeedManagementUpdateFeed: rate.NewLimiter(rate.Every(time.Second), 1),
		FeedManagementDeleteFeed: rate.NewLimiter(rate.Every(time.Second), 1),
		FeedManagementEnableFeed: rate.NewLimiter(rate.Every(time.Second), 1),

		DetectionCreateRule:         rate.NewLimiter(rate.Every(time.Second), 1),
		DetectionCreateRuleVersion:  rate.NewLimiter(rate.Every(time.Second), 1),
		DetectionGetRule:            rate.NewLimiter(rate.Every(time.Second), 1),
		DetectionUpdateRule:         rate.NewLimiter(rate.Every(time.Second), 1),
		DetectionDeleteRule:         rate.NewLimiter(rate.Every(time.Second), 1),
		DetectionEnableLiveRule:     rate.NewLimiter(rate.Every(time.Second), 1),
		DetectionEnableAlertingRule: rate.NewLimiter(rate.Every(time.Second), 1),
		DetectionVerifyYARARule:     rate.NewLimiter(rate.Every(time.Second), 1),

		RBACCreateSubject: rate.NewLimiter(rate.Every(time.Second), 1),
		RBACGetSubject:    rate.NewLimiter(rate.Every(time.Second), 1),
		RBACUpdateSubject: rate.NewLimiter(rate.Every(time.Second), 1),
		RBACDeleteSubject: rate.NewLimiter(rate.Every(time.Second), 1),

		ReferenceListsCreateList: rate.NewLimiter(rate.Every(time.Second), 1),
		ReferenceListsGetList:    rate.NewLimiter(rate.Every(time.Second), 1),
		ReferenceListsUpdateList: rate.NewLimiter(rate.Every(time.Second), 1),
	}
}

const (
	BigQueryAPIEnvVar     = "CHRONICLE_BIGQUERY_CREDENTIALS"
	BackstoryAPIEnvVar    = "CHRONICLE_BACKSTORY_CREDENTIALS"
	IngestionAPIEnvVar    = "CHRONICLE_INGESTION_CREDENTIALS"
	ForwarderAPIEnvVar    = "CHRONICLE_FORWARDER_CREDENTIALS"
	ChronicleRegionEnvVar = "CHRONICLE_REGION"
)

var EnvAPICrendetialsVar = []string{BigQueryAPIEnvVar, BackstoryAPIEnvVar, IngestionAPIEnvVar, ForwarderAPIEnvVar}

const (
	RegionUS             = "us"
	RegionEurope         = "europe"
	RegionEuropeWest2    = "europe-west2"
	RegionAsiaSouthEast1 = "asia-southeast1"
)

// Regions holds the Chronicle regions with regionalized legacy API endpoints
// ({region}-backstory.googleapis.com); "us" uses the global endpoints.
// See https://cloud.google.com/chronicle/docs/reference/feed-management-api.
var Regions = []string{
	RegionUS,
	RegionEurope,
	"africa-south1",
	"asia-northeast1",
	"asia-south1",
	RegionAsiaSouthEast1,
	"asia-southeast2",
	"australia-southeast1",
	"europe-central2",
	RegionEuropeWest2,
	"europe-west3",
	"europe-west6",
	"europe-west9",
	"europe-west12",
	"me-central1",
	"me-central2",
	"me-west1",
	"northamerica-northeast2",
	"southamerica-east1",
}

const APIDomain = "googleapis.com"

const backstorySubDomain = "backstory"

// regionalSubDomain derives the regional endpoint subdomain for a given API
// subdomain: the US region uses the global endpoint, every other region
// prefixes it (e.g. europe-backstory, asia-southeast1-malachiteingestion-pa).
func regionalSubDomain(subDomain, region string) string {
	if region == RegionUS || region == "" {
		return subDomain
	}
	return fmt.Sprintf("%s-%s", region, subDomain)
}

const (
	EventsBasePathKey   = "Events"
	AlertBasePathKey    = "Alert"
	ArtifactBasePathKey = "Artifact"
	AliasBasePathKey    = "Alias"
	AssetBasePathKey    = "Asset"
	IOCBasePathKey      = "IOC"

	RuleBasePathKey           = "rules"
	FeedManagementBasePathKey = "Feed"

	SubjectsBasePathKey = "Subjects"

	ReferenceListsPathKey = "ReferenceLists"
)

func GenerateDefaultBasePaths(region string) map[string]string {
	backstory := regionalSubDomain(backstorySubDomain, region)

	var DefaultBasePaths = map[string]string{
		EventsBasePathKey:   getBasePathFromDomainsAndPath("/v1/events", backstory),
		AlertBasePathKey:    getBasePathFromDomainsAndPath("/v1/alert", backstory),
		ArtifactBasePathKey: getBasePathFromDomainsAndPath("/v1/artifact", backstory),
		AliasBasePathKey:    getBasePathFromDomainsAndPath("/v1/alias", backstory),
		AssetBasePathKey:    getBasePathFromDomainsAndPath("/v1/asset", backstory),
		IOCBasePathKey:      getBasePathFromDomainsAndPath("/v1/ioc", backstory),

		RuleBasePathKey:           getBasePathFromDomainsAndPath("/v2/detect/rules", backstory),
		FeedManagementBasePathKey: getBasePathFromDomainsAndPath("/v1/feeds", backstory),

		SubjectsBasePathKey: getBasePathFromDomainsAndPath("/v1/subjects", backstory),

		ReferenceListsPathKey: getBasePathFromDomainsAndPath("/v2/lists", backstory),
	}

	return DefaultBasePaths
}

func getBasePathFromDomainsAndPath(basePath string, domain string) string {
	return fmt.Sprintf("https://%s.%s%s", domain, APIDomain, basePath)
}
