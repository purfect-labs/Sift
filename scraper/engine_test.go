package scraper

import (
	"testing"
)

func TestJobID(t *testing.T) {
	id1 := jobID("Software Engineer", "Google", "https://google.com/jobs/1")
	id2 := jobID("Software Engineer", "Google", "https://google.com/jobs/1")
	id3 := jobID("Software Engineer", "Meta", "https://google.com/jobs/1")

	if id1 != id2 {
		t.Errorf("same inputs should produce same ID: %s != %s", id1, id2)
	}
	if id1 == id3 {
		t.Errorf("different company should produce different ID")
	}
	if len(id1) != 16 {
		t.Errorf("job ID should be 16 chars, got %d", len(id1))
	}
}

func TestDefaultQueries(t *testing.T) {
	queries := DefaultQueries()
	if len(queries) == 0 {
		t.Error("DefaultQueries should return at least one query")
	}
	for _, q := range queries {
		if q == "" {
			t.Error("DefaultQueries should not contain empty strings")
		}
	}
}

func TestParseKeywordResponse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		minCount int
	}{
		{
			name:     "simple comma list",
			input:    "Kubernetes, Docker, Golang, Python, AWS, Terraform",
			minCount: 6,
		},
		{
			name:     "with meta commentary",
			input:    "Here are the keywords:\nKubernetes, Docker, Golang\nI'll sort them alphabetically.",
			minCount: 3,
		},
		{
			name:     "multi-line with spaces",
			input:    "  Kubernetes ,  Docker  ,  Golang  ,  Python",
			minCount: 4,
		},
		{
			name:     "single keyword",
			input:    "Kubernetes",
			minCount: 1,
		},
		{
			name:     "empty response",
			input:    "",
			minCount: 0,
		},
		{
			name:     "deduped note prefix",
			input:    "Deduped and sorted alphabetically:\nKubernetes, Docker, Python",
			minCount: 3,
		},
		{
			name:     "note prefix",
			input:    "Note: here are keywords\nKubernetes, Docker, Python",
			minCount: 3,
		},
		{
			name:     "let me prefix",
			input:    "Let me extract those:\nKubernetes, Docker, Python",
			minCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := ParseKeywordResponse(tt.input)
			if len(results) < tt.minCount {
				t.Errorf("expected at least %d keywords, got %d: %v", tt.minCount, len(results), results)
			}
			for _, kw := range results {
				if kw == "" {
					t.Error("should not return empty keyword")
				}
			}
		})
	}
}

func TestParseKeywordResponse_Dedup(t *testing.T) {
	results := ParseKeywordResponse("Go, go, GO, Docker, docker")
	if len(results) > 3 {
		t.Errorf("expected dedup to remove case duplicates, got %d: %v", len(results), results)
	}
}

func TestParseKeywordResponse_AllCommentary(t *testing.T) {
	// "Sure thing" survives because the backwards loop skips lines starting with
	// "Sure" but the fallback takes the last clean line anyway. Real Hermes output
	// would never be purely commentary — this tests the safety fallback.
	results := ParseKeywordResponse("Here are the keywords\nI'll list them\nSure thing")
	if len(results) == 0 {
		t.Error("fallback should return last line even if it looks like commentary")
	}
}

func TestComputeMatchScore(t *testing.T) {
	keywords := []string{"Kubernetes", "Docker", "Go", "AWS", "Terraform"}

	tests := []struct {
		name        string
		title       string
		company     string
		description string
		minScore    float64
		maxScore    float64
	}{
		{
			name:        "exact match on all keywords",
			title:       "Kubernetes Engineer",
			company:     "AWS",
			description: "Must know Docker, Go, and Terraform",
			minScore:    60,
			maxScore:    100,
		},
		{
			name:        "partial match",
			title:       "Kubernetes Engineer",
			company:     "Acme Corp",
			description: "Experience with container orchestration",
			minScore:    0,
			maxScore:    30,
		},
		{
			name:        "no match",
			title:       "Baker",
			company:     "Bakery Inc",
			description: "Must know how to bake bread",
			minScore:    0,
			maxScore:    0,
		},
		{
			name:        "empty keywords",
			title:       "Engineer",
			company:     "Tech",
			description: "Stuff",
			minScore:    0,
			maxScore:    0,
		},
		{
			name:        "single token match",
			title:       "DevOps Engineer",
			company:     "Cloud Co",
			description: "Experience with cloud platforms and containers",
			minScore:    0,
			maxScore:    15,
		},
		{
			name:        "substring match",
			title:       "Go Developer",
			company:     "Startup",
			description: "Building microservices",
			minScore:    15,
			maxScore:    25,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kws := keywords
			if tt.name == "empty keywords" {
				kws = nil
			}
			score := computeMatchScore(tt.title, tt.company, tt.description, kws)
			if score < tt.minScore || score > tt.maxScore {
				t.Errorf("score %f not in range [%f, %f]", score, tt.minScore, tt.maxScore)
			}
		})
	}
}

func TestComputeMatchScore_MultiWord(t *testing.T) {
	keywords := []string{"Director of Engineering", "Platform Engineering"}
	score := computeMatchScore(
		"Director of Platform Engineering",
		"Tech Corp",
		"Leading the platform engineering team",
		keywords,
	)
	if score < 50 {
		t.Errorf("multi-word keyword should match well, got %f", score)
	}
}

func TestComputeMatchScore_ExactVsPartial(t *testing.T) {
	keywords := []string{"Kubernetes"}

	exact := computeMatchScore("Kubernetes Engineer", "", "Must know Kubernetes deeply", keywords)
	partial := computeMatchScore("Container Engineer", "", "Orchestration experience", keywords)
	if exact <= partial {
		t.Errorf("exact match score (%f) should be higher than partial (%f)", exact, partial)
	}
}

func TestScrapeResult_Defaults(t *testing.T) {
	r := &ScrapeResult{}
	if r.JobsFound != 0 {
		t.Error("new ScrapeResult should have 0 JobsFound")
	}
	if len(r.Errors) != 0 {
		t.Error("new ScrapeResult should have no errors")
	}
}

func TestKeywordCount_Defaults(t *testing.T) {
	kc := &KeywordCount{}
	if kc.Keyword != "" {
		t.Error("Keyword should default to empty")
	}
	if kc.Count != 0 {
		t.Error("Count should default to 0")
	}
}
