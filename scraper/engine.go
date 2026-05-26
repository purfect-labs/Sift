package scraper

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"sort"
	"strings"
	"time"

	"sift/db"
)

type ScrapeResult struct {
	JobsFound int
	Errors    []string
}

type KeywordCount struct {
	Keyword string `json:"keyword"`
	Count   int    `json:"count"`
}

type LogFunc func(string)

func jobID(title, company, jURL string) string {
	h := sha256.New()
	h.Write([]byte(title + company + jURL))
	return fmt.Sprintf("%x", h.Sum(nil))[:16]
}

func HermesBin() string {
	if p, err := exec.LookPath("hermes"); err == nil {
		return p
	}
	return "hermes"
}

func PythonBin() string {
	if p, err := exec.LookPath("python3"); err == nil {
		return p
	}
	return "python3"
}

func safeTextSlice(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}

func ReadPDFText(pdfPath string) (string, error) {
	script := fmt.Sprintf(`import pymupdf; doc=pymupdf.open(%q); print(" ".join(p.get_text() for p in doc))`, pdfPath)
	cmd := exec.Command(PythonBin(), "-c", script)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("pdf read failed: %v — %s", err, string(out))
	}
	return strings.TrimSpace(string(out)), nil
}

func ExtractKeywordsFromPDF(pdfPath string) ([]string, string, error) {
	fullText, err := ReadPDFText(pdfPath)
	if err != nil {
		return nil, "", err
	}

	prompt := fmt.Sprintf(
		`You are a job matching engine. Read this resume and output ONLY a comma-separated list of 40-60 keywords and key phrases that best represent this candidate's skills, technologies, roles, and industry expertise. Include: job titles, technologies, tools, platforms, methodologies, certifications, and domain-specific terms. Be specific (e.g. "Kubernetes" not "containers", "Director of Engineering" not "management"). Output only the comma list, nothing else.

RESUME:
%s`, safeTextSlice(fullText, 6000))

	cmd := exec.Command(HermesBin(), "-z", prompt)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, "", fmt.Errorf("Hermes extraction failed: %v — output: %s", err, string(out))
	}

	response := strings.TrimSpace(string(out))
	keywords := ParseKeywordResponse(response)

	if len(keywords) < 10 {
		return nil, "", fmt.Errorf("Hermes returned too few keywords (%d). Response: %s", len(keywords), safeTextSlice(response, 200))
	}

	return keywords, fullText, nil
}

// ParseKeywordResponse extracts comma-separated keywords from Hermes output
func ParseKeywordResponse(response string) []string {
	// Strip lines that are clearly meta-commentary from Hermes
	lines := strings.Split(response, "\n")
	var cleanLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		lower := strings.ToLower(line)
		// Skip meta-commentary lines
		if strings.Contains(lower, "deduped") || strings.Contains(lower, "sorted alphabetically") ||
			strings.Contains(lower, "here are") || strings.Contains(lower, "i'll") ||
			strings.Contains(lower, "let me") || strings.HasPrefix(lower, "note:") ||
			strings.HasPrefix(lower, "actually") || strings.HasPrefix(lower, "however") {
			continue
		}
		if line != "" {
			cleanLines = append(cleanLines, line)
		}
	}

	var keywordLine string
	for i := len(cleanLines) - 1; i >= 0; i-- {
		line := cleanLines[i]
		if !strings.HasPrefix(line, "I") && !strings.HasPrefix(line, "Here") && !strings.HasPrefix(line, "Sure") {
			keywordLine = line
			break
		}
	}
	if keywordLine == "" && len(cleanLines) > 0 {
		keywordLine = cleanLines[len(cleanLines)-1]
	}
	if keywordLine == "" {
		keywordLine = response
	}

	parts := strings.Split(keywordLine, ",")
	var keywords []string
	seen := map[string]bool{}
	for _, p := range parts {
		kw := strings.TrimSpace(p)
		kw = strings.Trim(kw, ".-*\"'")
		lk := strings.ToLower(kw)
		if len(kw) > 2 && !seen[lk] {
			seen[lk] = true
			keywords = append(keywords, kw)
		}
	}
	return keywords
}

// computeMatchScore compares job content against resume keywords with multi-level matching
func computeMatchScore(title, company, description string, keywords []string) float64 {
	if len(keywords) == 0 {
		return 0
	}
	text := strings.ToLower(title + " " + company + " " + description)
	var totalScore float64
	var maxPossible float64

	for _, kw := range keywords {
		kwLower := strings.ToLower(strings.TrimSpace(kw))
		maxPossible += 3.0

		// Level 1: exact phrase match (3 points)
		if strings.Contains(text, kwLower) {
			totalScore += 3.0
			continue
		}

		// Level 2: all tokens present (2 points)
		tokens := strings.Fields(kwLower)
		if len(tokens) > 1 {
			allMatch := true
			for _, t := range tokens {
				if len(t) > 2 && !strings.Contains(text, t) {
					allMatch = false
					break
				}
			}
			if allMatch {
				totalScore += 2.0
				continue
			}
		}

		// Level 3: single token of multi-word keyword matches (1 point)
		if len(tokens) > 1 {
			for _, t := range tokens {
				if len(t) > 3 && strings.Contains(text, t) {
					totalScore += 1.0
					break
				}
			}
			continue
		}

		// Level 4: single-word keyword substring match (1 point)
		if len(kwLower) > 3 && strings.Contains(text, kwLower) {
			totalScore += 1.0
		}
	}

	pct := (totalScore / maxPossible) * 100
	if pct > 100 {
		pct = 100
	}
	return float64(int(pct))
}

// AggregateKeywords returns top N most common keywords across all scraped jobs (filler-filtered)
func AggregateKeywords(limit int) ([]KeywordCount, error) {
	jobs, err := db.GetAllJobs("all")
	if err != nil {
		return nil, err
	}

	filler := map[string]bool{
		"the": true, "and": true, "with": true, "for": true, "that": true, "this": true,
		"have": true, "from": true, "are": true, "not": true, "but": true, "or": true,
		"as": true, "at": true, "by": true, "in": true, "of": true, "to": true, "is": true,
		"it": true, "on": true, "be": true, "we": true, "you": true, "they": true, "will": true,
		"can": true, "all": true, "has": true, "been": true, "was": true, "were": true,
		"their": true, "our": true, "your": true, "its": true, "an": true, "if": true,
		"no": true, "so": true, "up": true, "out": true, "just": true, "also": true,
		"about": true, "into": true, "over": true, "after": true, "than": true, "then": true,
		"only": true, "other": true, "some": true, "these": true, "those": true, "each": true,
		"every": true, "any": true, "such": true, "both": true, "few": true, "more": true,
		"most": true, "very": true, "much": true, "many": true, "well": true, "even": true,
		"still": true, "already": true, "always": true, "never": true, "often": true,
		"usually": true, "sometimes": true, "now": true, "new": true, "first": true,
		"last": true, "next": true, "same": true, "own": true, "part": true, "place": true,
		"year": true, "years": true, "time": true, "day": true, "work": true,
		"experience": true, "including": true, "using": true, "required": true, "must": true,
		"need": true, "strong": true, "excellent": true, "good": true, "great": true,
		"ability": true, "knowledge": true, "understanding": true,
		"familiar": true, "plus": true, "bonus": true, "opportunity": true, "looking": true,
		"seeking": true, "join": true, "team": true, "help": true, "build": true, "role": true,
		"position": true, "candidate": true, "ideal": true, "passionate": true,
		"fast": true, "growing": true, "environment": true, "company": true, "companies": true,
	}

	counts := map[string]int{}
	for _, job := range jobs {
		for _, skill := range strings.Split(job.Skills, ",") {
			s := strings.TrimSpace(strings.ToLower(skill))
			if len(s) > 2 && !filler[s] {
				counts[s]++
			}
		}
	}

	// Convert to sorted slice
	var result []KeywordCount
	for kw, count := range counts {
		result = append(result, KeywordCount{Keyword: kw, Count: count})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})

	if limit > 0 && limit < len(result) {
		result = result[:limit]
	}
	return result, nil
}

// AnalyzeJobDemand sends all job descriptions to Hermes for market analysis
func AnalyzeJobDemand(keywords, jobSamples []string) (string, error) {
	sampleText := strings.Join(jobSamples, "\n---\n")
	if len(sampleText) > 8000 {
		sampleText = sampleText[:8000]
	}

	prompt := fmt.Sprintf(
		`You are a job market analyst. Below are resume keywords and sample job descriptions that matched. Output a concise analysis covering:
1. Top 5 most in-demand skills/technologies across these jobs
2. Which of the candidate's keywords have the highest demand
3. Any skill gaps - things companies want that the candidate doesn't list
4. Salary/level trends you notice
5. One actionable recommendation

RESUME KEYWORDS: %s

JOB SAMPLES:
%s`, strings.Join(keywords, ", "), sampleText)

	cmd := exec.Command(HermesBin(), "-z", prompt)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Hermes failed: %v", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func SearchSerpAPI(query, location, apiKey string, keywords []string, logFn LogFunc) (*ScrapeResult, error) {
	result := &ScrapeResult{}

	if apiKey == "" {
		result.Errors = append(result.Errors, "No SerpAPI key configured.")
		return result, nil
	}

	params := url.Values{}
	params.Add("engine", "google_jobs")
	params.Add("q", query)
	params.Add("api_key", apiKey)
	if location != "" && location != "Remote" {
		params.Add("location", location)
	}
	if location == "Remote" {
		query = query + " remote"
		params.Set("q", query)
	}

	apiURL := "https://serpapi.com/search?" + params.Encode()

	if logFn != nil {
		logFn(fmt.Sprintf("GET %s", "https://serpapi.com/search?engine=google_jobs&q="+url.QueryEscape(query)+"&location="+url.QueryEscape(location)))
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(apiURL)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("API request failed: %v", err))
		return result, nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if logFn != nil {
		logFn(fmt.Sprintf("HTTP %d — %d bytes", resp.StatusCode, len(body)))
	}

	if resp.StatusCode != 200 {
		preview := string(body)
		if len(preview) > 300 {
			preview = preview[:300]
		}
		if logFn != nil {
			logFn(fmt.Sprintf("ERROR: %s", preview))
		}
		result.Errors = append(result.Errors, fmt.Sprintf("HTTP %d: %s", resp.StatusCode, preview))
		return result, nil
	}

	var data struct {
		JobsResults []struct {
			Title              string `json:"title"`
			CompanyName        string `json:"company_name"`
			Location           string `json:"location"`
			Description        string `json:"description"`
			DetectedExtensions struct {
				Salary string `json:"salary"`
			} `json:"detected_extensions"`
			ApplyOptions []struct {
				Link string `json:"link"`
			} `json:"apply_options"`
			JobHighlights []struct {
				Items []string `json:"items"`
			} `json:"job_highlights"`
		} `json:"jobs_results"`
		Error string `json:"error"`
	}

	if err := json.Unmarshal(body, &data); err != nil {
		if logFn != nil {
			logFn(fmt.Sprintf("JSON parse error: %v", err))
		}
		return result, nil
	}

	if data.Error != "" {
		if logFn != nil {
			logFn(fmt.Sprintf("API error: %s", data.Error))
		}
		result.Errors = append(result.Errors, data.Error)
		return result, nil
	}

	if logFn != nil {
		logFn(fmt.Sprintf("Found %d job results", len(data.JobsResults)))
	}

	for _, j := range data.JobsResults {
		jobURL := ""
		if len(j.ApplyOptions) > 0 {
			jobURL = j.ApplyOptions[0].Link
		}

		var skills []string
		for _, h := range j.JobHighlights {
			skills = append(skills, h.Items...)
		}

		id := jobID(j.Title, j.CompanyName, jobURL)
		if db.JobExists(id) {
			continue
		}

		job := &db.Job{
			ID:          id,
			Title:       j.Title,
			Company:     j.CompanyName,
			Location:    j.Location,
			Salary:      j.DetectedExtensions.Salary,
			Description: j.Description,
			Skills:      strings.Join(skills, ", "),
			URL:         jobURL,
			Source:      "google_jobs",
			MatchScore:  computeMatchScore(j.Title, j.CompanyName, j.Description, keywords),
			Status:      "new",
			ScrapedAt:   time.Now().Format(time.RFC3339),
		}

		if err := db.InsertJob(job); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("DB insert failed for %s: %v", j.Title, err))
			continue
		}
		result.JobsFound++
	}

	return result, nil
}

func SearchMultiple(queries []string, location, apiKey string, keywords []string, logFn LogFunc) (*ScrapeResult, error) {
	total := &ScrapeResult{}
	for _, q := range queries {
		if logFn != nil {
			logFn(fmt.Sprintf("--- Query: %s ---", q))
		}
		r, err := SearchSerpAPI(q, location, apiKey, keywords, logFn)
		if err != nil {
			return total, err
		}
		total.JobsFound += r.JobsFound
		total.Errors = append(total.Errors, r.Errors...)
	}
	return total, nil
}

func DefaultQueries() []string {
	return []string{
		"director of engineering platform",
		"head of engineering devops platform",
		"staff platform engineer infrastructure",
		"principal devops engineer sre",
		"ai engineering manager platform",
	}
}
