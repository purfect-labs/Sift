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
