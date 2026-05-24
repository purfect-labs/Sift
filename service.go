package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"jobdash/db"
	"jobdash/scraper"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type Config struct {
	SerpAPIKey string   `json:"serpapi_key"`
	Location   string   `json:"location"`
	Queries    []string `json:"queries"`
	ResumePath string   `json:"resume_path"`
	ResumeText string   `json:"resume_text"`
	Keywords   []string `json:"keywords"`
}

type JobService struct {
	config   Config
	keywords []string
	logs     []string
	app      *application.App
}

func loadEnv(path string) map[string]string {
	env := map[string]string{}
	data, err := os.ReadFile(path)
	if err != nil {
		return env
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			env[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return env
}

func safeSlice(s string, n int) string {
	if n <= 0 {
		return ""
	}
	if len(s) <= n {
		return s
	}
	return s[:n]
}

func (s *JobService) log(msg string) {
	s.logs = append(s.logs, msg)
	if s.app != nil {
		s.app.Event.Emit("scrape_log", msg)
	}
}

func (s *JobService) Init() error {
	home, _ := os.UserHomeDir()
	dbPath := filepath.Join(home, ".jobdash", "jobs.db")
	os.MkdirAll(filepath.Dir(dbPath), 0700)

	if err := db.Init(dbPath); err != nil {
		return err
	}

	env := loadEnv(filepath.Join(home, ".jobdash", ".env"))
	if key := env["SERP_API_KEY"]; key != "" {
		s.config.SerpAPIKey = key
	}

	kw, _ := db.LoadKeywords("resume")
	if len(kw) > 0 {
		s.keywords = kw
		s.config.Keywords = kw
	}
	if v, err := db.GetConfig("resume_path"); err == nil && v != "" {
		s.config.ResumePath = v
	}
	if v, err := db.GetConfig("location"); err == nil && v != "" {
		s.config.Location = v
	}

	return nil
}

func (s *JobService) GetJobs(status string) ([]db.Job, error) {
	return db.GetAllJobs(status)
}

func (s *JobService) UpdateStatus(id, status string) error {
	return db.UpdateJobStatus(id, status)
}

func (s *JobService) UpdateNotes(id, notes string) error {
	return db.UpdateJobNotes(id, notes)
}

func (s *JobService) GetCounts() (map[string]int, error) {
	return db.GetJobCounts()
}

func (s *JobService) SetConfig(key, value string) error {
	switch key {
	case "location":
		s.config.Location = value
		db.SetConfig("location", value)
	case "resume_path":
		s.config.ResumePath = value
		db.SetConfig("resume_path", value)
		kw, _ := db.LoadKeywords("resume")
		if len(kw) > 0 {
			s.keywords = kw
			s.config.Keywords = kw
		}
	case "serpapi_key":
		s.config.SerpAPIKey = value
		home, _ := os.UserHomeDir()
		envPath := filepath.Join(home, ".jobdash", ".env")
		os.MkdirAll(filepath.Dir(envPath), 0700)
		// Read existing .env, update SERP_API_KEY line, write back
		existing, _ := os.ReadFile(envPath)
		lines := strings.Split(string(existing), "\n")
		found := false
		for i, line := range lines {
			if strings.HasPrefix(strings.TrimSpace(line), "SERP_API_KEY=") {
				lines[i] = "SERP_API_KEY=" + value
				found = true
				break
			}
		}
		if !found {
			lines = append(lines, "SERP_API_KEY="+value)
		}
		return os.WriteFile(envPath, []byte(strings.Join(lines, "\n")), 0600)
	}
	return nil
}

func (s *JobService) GetConfig() Config {
	return s.config
}

func (s *JobService) GetPositions() ([]string, error) {
	return db.LoadPositions()
}

func (s *JobService) SavePositions(positions []string) error {
	return db.SavePositions(positions)
}

func (s *JobService) RecommendPositions() ([]string, error) {
	if s.config.ResumePath == "" {
		return nil, fmt.Errorf("Set resume path first")
	}

	s.log("Asking Hermes to recommend job titles...")

	fullText, err := scraper.ReadPDFText(s.config.ResumePath)
	if err != nil {
		return nil, err
	}

	prompt := fmt.Sprintf(
		`Read this resume and output ONLY a comma-separated list of 5-8 specific job titles this person should target. Be specific: use real titles like "Director of Platform Engineering" not generic ones. Output only the comma list.

RESUME:
%s`, safeSlice(fullText, 5000))

	cmd := exec.Command(scraper.HermesBin(), "-z", prompt)
	out, err := cmd.CombinedOutput()
	if err != nil {
		s.log(fmt.Sprintf("Hermes error: %v", err))
		return nil, fmt.Errorf("Hermes failed: %v", err)
	}

	response := strings.TrimSpace(string(out))
	var titles []string
	for _, t := range strings.Split(response, ",") {
		t = strings.TrimSpace(t)
		t = strings.Trim(t, ".-*\"'")
		if len(t) > 3 {
			titles = append(titles, t)
		}
	}

	if len(titles) < 3 {
		return nil, fmt.Errorf("Hermes returned too few titles: %s", safeSlice(response, 200))
	}

	db.SavePositions(titles)
	s.log(fmt.Sprintf("Hermes recommended %d positions: %s", len(titles), strings.Join(titles, ", ")))
	return titles, nil
}

func (s *JobService) ExtractResume(resumePath string) (*ScrapeResult, error) {
	s.logs = nil // clear previous logs
	s.log(fmt.Sprintf("Reading resume: %s", resumePath))
	s.log("Sending to Hermes for keyword extraction...")

	keywords, text, err := scraper.ExtractKeywordsFromPDF(resumePath)
	if err != nil {
		s.log(fmt.Sprintf("ERROR: %v", err))
		return &ScrapeResult{Errors: []string{err.Error()}}, nil
	}

	db.SaveKeywords("resume", keywords)

	s.keywords = keywords
	s.config.Keywords = keywords
	s.config.ResumeText = text
	s.config.ResumePath = resumePath

	s.log(fmt.Sprintf("Extracted %d keywords: %s", len(keywords), strings.Join(keywords[:min(10, len(keywords))], ", ")+"..."))
	return &ScrapeResult{JobsFound: len(keywords), Errors: s.logs}, nil
}

// GetExtractLog returns the accumulated log from the last extraction
func (s *JobService) GetExtractLog() []string {
	return s.logs
}

type ScrapeResult struct {
	JobsFound int      `json:"jobs_found"`
	Errors    []string `json:"errors"`
}

func (s *JobService) ScrapeJobs(customQuery string) (*ScrapeResult, error) {
	s.log("=== Starting job scrape ===")

	if len(s.keywords) == 0 {
		kw, _ := db.LoadKeywords("resume")
		if len(kw) > 0 {
			s.keywords = kw
			s.config.Keywords = kw
			s.log(fmt.Sprintf("Loaded %d keywords from DB", len(kw)))
		}
	}

	if len(s.keywords) == 0 && s.config.ResumePath != "" {
		s.log("No keywords found, auto-extracting from resume...")
		if _, err := s.ExtractResume(s.config.ResumePath); err != nil {
			return &ScrapeResult{Errors: []string{fmt.Sprintf("Auto-extract failed: %v", err)}}, nil
		}
	}

	if len(s.keywords) == 0 {
		return &ScrapeResult{Errors: []string{"No resume keywords. Set resume path in Settings and click Extract Keywords."}}, nil
	}

	if s.config.SerpAPIKey == "" {
		return &ScrapeResult{Errors: []string{"No SerpAPI key. Add SERP_API_KEY=your-key to ~/.jobdash/.env"}}, nil
	}

	queries := s.config.Queries
	if len(queries) == 0 {
		positions, _ := db.LoadPositions()
		if len(positions) > 0 {
			queries = positions
		} else {
			queries = scraper.DefaultQueries()
		}
	}
	if customQuery != "" {
		queries = []string{customQuery}
	}

	location := s.config.Location
	if location == "" {
		location = "Remote"
	}

	s.log("Queries:")
	for _, q := range queries {
		s.log(fmt.Sprintf("  - %s", q))
	}
	s.log(fmt.Sprintf("Location: %s", location))
	s.log(fmt.Sprintf("Keywords: %d loaded", len(s.keywords)))
	s.log("Calling SerpAPI...")

	result, err := scraper.SearchMultiple(queries, location, s.config.SerpAPIKey, s.keywords, s.log)
	if err != nil {
		s.log(fmt.Sprintf("ERROR: %v", err))
		return nil, err
	}

	s.log(fmt.Sprintf("Found %d new jobs", result.JobsFound))
	for _, e := range result.Errors {
		s.log(fmt.Sprintf("WARN: %s", e))
	}

	// Auto-run Hermes keyword extraction + gap analysis on new jobs
	if result.JobsFound > 0 {
		s.log("Running Hermes keyword extraction + match rating on new jobs...")
		s.analyzeNewJobs()
		s.log("Running gap analysis against your resume...")
		s.analyzeGapsForNewJobs()
	}

	s.log("=== Scrape complete ===")
	return &ScrapeResult{JobsFound: result.JobsFound, Errors: result.Errors}, nil
}

func (s *JobService) GetDefaultQueries() []string {
	return scraper.DefaultQueries()
}

var fillerWords = map[string]bool{
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
	"ability": true, "skills": true, "knowledge": true, "understanding": true,
	"familiar": true, "plus": true, "bonus": true, "opportunity": true, "looking": true,
	"seeking": true, "join": true, "team": true, "help": true, "build": true, "role": true,
	"position": true, "candidate": true, "ideal": true, "passionate": true,
	"fast": true, "growing": true, "environment": true, "company": true, "companies": true,
}

func filterFillerWords(keywords []string) []string {
	var filtered []string
	for _, kw := range keywords {
		kw = strings.TrimSpace(kw)
		lower := strings.ToLower(kw)
		if !fillerWords[lower] && len(kw) > 2 {
			filtered = append(filtered, kw)
		}
	}
	return filtered
}

func (s *JobService) DeleteJob(id string) error {
	return db.DeleteJob(id)
}

func (s *JobService) ClearAllJobs() error {
	n, err := db.ClearAllJobs()
	if err != nil {
		s.log(fmt.Sprintf("Clear failed: %v", err))
		return err
	}
	s.log(fmt.Sprintf("Cleared %d jobs", n))
	return nil
}

func (s *JobService) UpdateOffer(id, salary, benefits, equity string) error {
	return db.UpdateOfferDetails(id, salary, benefits, equity)
}

func (s *JobService) GetGapAnalysis(jobID string) (string, error) {
	jobs, err := db.GetAllJobs("all")
	if err != nil {
		return "", err
	}
	var job *db.Job
	for i := range jobs {
		if jobs[i].ID == jobID {
			job = &jobs[i]
			break
		}
	}
	if job == nil || len(s.keywords) == 0 {
		return "", nil
	}

	resumeSummary := strings.Join(s.keywords, ", ")
	jobText := fmt.Sprintf("TITLE: %s\nCOMPANY: %s\nSKILLS: %s\nDESC: %s", job.Title, job.Company, job.Skills, job.Description[:1500])

	prompt := fmt.Sprintf(
		`Compare this candidate's resume against a job listing. Identify ONLY genuine skill gaps — things the job requires that the candidate truly does NOT have. Use common sense: if the candidate lists "Kubernetes", do NOT flag "container orchestration" as a gap. If they list "Python", do NOT flag "scripting" or "programming". Only flag specific technologies, tools, or domain expertise that is clearly missing.

Output: a short comma-separated list of 3-5 genuine gaps. If the candidate covers everything well, say "None - strong match". Keep it brief.

RESUME: %s

JOB: %s`, resumeSummary, jobText)

	cmd := exec.Command(scraper.HermesBin(), "-z", prompt)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", nil
	}
	return strings.TrimSpace(string(out)), nil
}

func (s *JobService) GetKeywordStats() ([]scraper.KeywordCount, error) {
	return scraper.AggregateKeywords(100)
}

func (s *JobService) GetMarketAnalysis() (string, error) {
	jobs, err := db.GetAllJobs("all")
	if err != nil {
		return "", err
	}
	var samples []string
	for i, j := range jobs {
		if i >= 20 {
			break
		}
		if j.Description != "" {
			samples = append(samples, fmt.Sprintf("%s at %s: %s", j.Title, j.Company, j.Description[:300]))
		}
	}
	if len(samples) == 0 {
		return "No job data to analyze. Scrape some jobs first.", nil
	}
	return scraper.AnalyzeJobDemand(s.keywords, samples)
}

// remoteKeywords: positive = remote signal, negative = not remote
var remotePositive = []string{
	"remote", "work from home", "wfh", "fully remote", "100% remote",
	"100% remote", "work remotely", "telecommute", "telecommuting",
	"anywhere", "distributed team", "distributed workforce",
	"home office", "home based", "home-based", "virtual position",
	"virtual role", "virtual office", "remote first", "remote-first",
	"remote friendly", "remote ok", "remoteok", "location independent",
	"location-independent", "no location requirement", "work anywhere",
	"no office", "no relocation", "full time remote", "permanent remote",
	"remote eligible", "remote position", "remote work", "remote job",
	"remote option", "remote opportunity", "flexible location",
	"work from anywhere", "wfa", "digital nomad",
}

var remoteNegative = []string{
	"hybrid", "on-site", "onsite", "on site", "in-office", "in office",
	"must relocate", "relocation required", "relocation assistance",
	"must be in", "must be based in", "must reside in", "must live in",
	"required in office", "required to be in", "no remote",
	"not remote", "in-person", "in person", "on location",
	"office based", "office-based", "on campus", "on-campus",
}

// FilterRemoteJobs scans all new/unfiltered jobs and marks them remote/not_remote using keyword matching
func (s *JobService) FilterRemoteJobs() (*ScrapeResult, error) {
	jobs, err := db.GetAllJobs("new")
	if err != nil {
		return nil, err
	}
	if len(jobs) == 0 {
		allJobs, _ := db.GetAllJobs("all")
		for _, j := range allJobs {
			if j.Status != "not_remote" && j.Status != "saved" && j.Status != "applied" && j.Status != "interviewing" && j.Status != "offer" {
				jobs = append(jobs, j)
			}
		}
	}

	var remoteCount, notRemoteCount int
	for _, j := range jobs {
		text := strings.ToLower(j.Title + " " + j.Location + " " + j.Description)
		score := 0
		for _, kw := range remotePositive {
			if strings.Contains(text, kw) {
				score += 2
			}
		}
		for _, kw := range remoteNegative {
			if strings.Contains(text, kw) {
				score -= 3
			}
		}
		if score > 0 {
			remoteCount++
			db.UpdateJobStatus(j.ID, "saved")
		} else {
			notRemoteCount++
			db.UpdateJobStatus(j.ID, "not_remote")
		}
	}

	s.log(fmt.Sprintf("Pre-filter: %d remote, %d filtered out of %d jobs", remoteCount, notRemoteCount, len(jobs)))

	// Step 2: Run Hermes on remote jobs for smart keyword extraction
	if remoteCount > 0 {
		s.log("Running Hermes keyword analysis on remote jobs...")
		s.analyzeRemoteJobs()
	}

	return &ScrapeResult{
		JobsFound: remoteCount,
		Errors:    []string{fmt.Sprintf("%d remote kept, %d filtered out", remoteCount, notRemoteCount)},
	}, nil
}

func (s *JobService) analyzeRemoteJobs() {
	jobs, err := db.GetAllJobs("saved")
	if err != nil || len(jobs) == 0 {
		return
	}
	s.analyzeBatch(jobs)
}

func (s *JobService) analyzeNewJobs() {
	jobs, err := db.GetAllJobs("new")
	if err != nil || len(jobs) == 0 {
		return
	}
	s.analyzeBatch(jobs)
	s.rateMatchBatch(jobs)
}

func (s *JobService) rateMatchBatch(jobs []db.Job) {
	if len(s.keywords) == 0 {
		return
	}
	resumeSummary := strings.Join(s.keywords, ", ")

	for i := 0; i < len(jobs); i += 10 {
		end := i + 10
		if end > len(jobs) {
			end = len(jobs)
		}
		batch := jobs[i:end]

		var jobList []string
		for _, j := range batch {
			desc := j.Description
			if len(desc) > 400 {
				desc = desc[:400]
			}
			jobList = append(jobList, fmt.Sprintf("ID:%s | TITLE:%s | COMPANY:%s | DESC:%s", j.ID, j.Title, j.Company, desc))
		}

		prompt := fmt.Sprintf(
			`You are a job match rater. Compare this candidate's resume keywords against each job listing. For each job ID, output a match percentage (0-100) based on how well the candidate's skills and experience fit the role. Be honest — don't inflate scores. Output ONLY a JSON array: [{"id":"...","score":65},...]

RESUME KEYWORDS: %s

JOBS:
%s`, resumeSummary, strings.Join(jobList, "\n---\n"))

		cmd := exec.Command(scraper.HermesBin(), "-z", prompt)
		out, err := cmd.CombinedOutput()
		if err != nil {
			s.log(fmt.Sprintf("Rating batch %d failed: %v", i/10, err))
			continue
		}

		response := strings.TrimSpace(string(out))
		jStart := strings.Index(response, "[")
		jEnd := strings.LastIndex(response, "]")
		if jStart == -1 || jEnd == -1 {
			continue
		}

		var ratings []struct {
			ID    string  `json:"id"`
			Score float64 `json:"score"`
		}
		if err := json.Unmarshal([]byte(response[jStart:jEnd+1]), &ratings); err != nil {
			continue
		}

		for _, r := range ratings {
			if r.Score > 0 {
				db.UpdateMatchScore(r.ID, r.Score)
			}
		}
		s.log(fmt.Sprintf("Rating batch %d: rated %d jobs", i/10, len(batch)))
	}
}

func (s *JobService) analyzeGapsForNewJobs() {
	jobs, err := db.GetAllJobs("new")
	if err != nil || len(jobs) == 0 {
		return
	}
	if len(s.keywords) == 0 {
		return
	}

	resumeSummary := strings.Join(s.keywords, ", ")

	for i := 0; i < len(jobs); i += 10 {
		end := i + 10
		if end > len(jobs) {
			end = len(jobs)
		}
		batch := jobs[i:end]

		var jobList []string
		for _, j := range batch {
			desc := j.Description
			if len(desc) > 300 {
				desc = desc[:300]
			}
			jobList = append(jobList, fmt.Sprintf("ID:%s | %s at %s | SKILLS:%s | DESC:%s", j.ID, j.Title, j.Company, j.Skills, desc))
		}

		prompt := fmt.Sprintf(
			`Compare this candidate's resume against each job listing. For each job ID, identify only GENUINE skill gaps — things the job requires that the candidate truly lacks. Use common sense: if resume has "Kubernetes", do NOT flag "container orchestration". If resume has "Python", do NOT flag "scripting". Only flag specific technologies/tools/domain expertise clearly missing. Output a JSON array: [{"id":"...","gaps":"short list or 'none'"},...]

RESUME: %s

JOBS:
%s`, resumeSummary, strings.Join(jobList, "\n---\n"))

		cmd := exec.Command(scraper.HermesBin(), "-z", prompt)
		out, err := cmd.CombinedOutput()
		if err != nil {
			continue
		}

		response := strings.TrimSpace(string(out))
		jStart := strings.Index(response, "[")
		jEnd := strings.LastIndex(response, "]")
		if jStart == -1 || jEnd == -1 {
			continue
		}

		var results []struct {
			ID   string `json:"id"`
			Gaps string `json:"gaps"`
		}
		if err := json.Unmarshal([]byte(response[jStart:jEnd+1]), &results); err != nil {
			continue
		}

		for _, r := range results {
			if r.Gaps != "" && r.Gaps != "none" {
				db.UpdateGapAnalysis(r.ID, r.Gaps)
			}
		}
		s.log(fmt.Sprintf("Gap batch %d: analyzed %d jobs", i/10, len(batch)))
	}
}

func (s *JobService) analyzeBatch(jobs []db.Job) {

	batchSize := 20
	for i := 0; i < len(jobs); i += batchSize {
		end := i + batchSize
		if end > len(jobs) {
			end = len(jobs)
		}
		batch := jobs[i:end]

		var jobList []string
		for _, j := range batch {
			desc := j.Description
			if len(desc) > 400 {
				desc = desc[:400]
			}
			jobList = append(jobList, fmt.Sprintf("TITLE: %s | COMPANY: %s | DESC: %s", j.Title, j.Company, desc))
		}

		prompt := fmt.Sprintf(
			`Extract STRICTLY technical skills, tools, platforms, and job requirements from these job listings. Output ONLY a comma-separated list of keywords.

WHAT TO OUTPUT (examples): Kubernetes, Terraform, Python, AWS, CI/CD pipelines, Docker, Go, machine learning, TypeScript, PostgreSQL, Redis, microservices, gRPC, React, Svelte, Helm, GitLab CI, Jenkins, Prometheus, Grafana, Kafka, Elasticsearch, MongoDB, GraphQL, REST APIs, Linux, system design, team leadership

WHAT TO NEVER OUTPUT: the, and, with, for, that, this, have, from, are, not, but, or, as, at, by, in, of, to, is, it, on, be, we, you, they, will, can, all, has, been, was, were, their, our, your, its, an, if, no, so, up, out, just, also, about, into, over, after, than, then, only, other, some, these, those, each, every, any, such, both, few, more, most, very, much, many, well, even, still, already, always, never, often, usually, sometimes, now, new, first, last, next, same, own, part, place, year, years, time, day, work, experience, including, using, required, must, need, strong, excellent, good, great, ability, skills, knowledge, understanding, familiar, plus, bonus, opportunity, looking, seeking, join

OUTPUT FORMAT: Kubernetes, AWS, Terraform, Python, CI/CD (comma-separated ONLY, no other text)

JOB LISTINGS:
%s`, strings.Join(jobList, "\n---\n"))

		cmd := exec.Command(scraper.HermesBin(), "-z", prompt)
		out, err := cmd.CombinedOutput()
		if err != nil {
			s.log(fmt.Sprintf("Hermes batch %d failed: %v", i/batchSize, err))
			continue
		}

		response := strings.TrimSpace(string(out))
		keywords := scraper.ParseKeywordResponse(response)

		// Hard filter: strip known filler/soft words that Hermes might leak
		keywords = filterFillerWords(keywords)

		if len(keywords) > 0 {
			// Store extracted keywords in job skills field
			enriched := strings.Join(keywords, ", ")
			for _, j := range batch {
				db.UpdateJobSkills(j.ID, enriched)
			}
			s.log(fmt.Sprintf("Batch %d: extracted %d keywords", i/batchSize, len(keywords)))
		}
	}
}

type JobAnalysis struct {
	SalaryRange  string `json:"salary_range"`
	LocationType string `json:"location_type"`
	Experience   string `json:"experience"`
	MustHave     string `json:"must_have"`
	NiceToHave   string `json:"nice_to_have"`
	RedFlags     string `json:"red_flags"`
	Summary      string `json:"summary"`
	Error        string `json:"error,omitempty"`
}

func (s *JobService) AnalyzeJob(title, company, description string) (*JobAnalysis, error) {
	if description == "" {
		return &JobAnalysis{Error: "No description to analyze"}, nil
	}

	prompt := fmt.Sprintf(
		`Analyze this job posting. Output ONLY valid JSON with these fields (use "N/A" if not found):
{
  "salary_range": "extracted salary or range",
  "location_type": "remote, hybrid, or onsite",
  "experience": "years required",
  "must_have": "top 3 must-have skills or requirements",
  "nice_to_have": "top 3 nice-to-have skills",
  "red_flags": "any red flags like on-call, legacy tech, or vague requirements",
  "summary": "one sentence TL;DR of what this role actually is"
}

JOB: %s at %s

DESCRIPTION: %s`, title, company, description[:3000])

	cmd := exec.Command(scraper.HermesBin(), "-z", prompt)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return &JobAnalysis{Error: fmt.Sprintf("Hermes failed: %v", err)}, nil
	}

	// Extract JSON from response
	response := strings.TrimSpace(string(out))
	// Find JSON block
	start := strings.Index(response, "{")
	end := strings.LastIndex(response, "}")
	if start == -1 || end == -1 {
		return &JobAnalysis{Error: "Could not parse analysis"}, nil
	}

	jsonStr := response[start : end+1]
	// Fix common JSON issues from LLM output
	jsonStr = strings.ReplaceAll(jsonStr, "`", "")

	var analysis JobAnalysis
	if err := json.Unmarshal([]byte(jsonStr), &analysis); err != nil {
		return &JobAnalysis{Error: fmt.Sprintf("Parse error: %v", err)}, nil
	}

	return &analysis, nil
}
