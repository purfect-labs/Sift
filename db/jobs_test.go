package db

import (
	"testing"
)

func setupTestDB(t *testing.T) {
	t.Helper()
	if err := Init(":memory:"); err != nil {
		t.Fatalf("failed to init test DB: %v", err)
	}
}

func TestInit(t *testing.T) {
	err := Init(":memory:")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	if DB == nil {
		t.Fatal("DB should not be nil after Init")
	}
}

func TestInsertAndGetJob(t *testing.T) {
	setupTestDB(t)

	job := &Job{
		ID:          "test-001",
		Title:       "Software Engineer",
		Company:     "Acme Corp",
		Location:    "Remote",
		Description: "Build cool stuff",
		Skills:      "Go, Docker",
		URL:         "https://acme.com/jobs/1",
		Source:      "google_jobs",
		MatchScore:  85.5,
		Status:      "new",
	}

	if err := InsertJob(job); err != nil {
		t.Fatalf("InsertJob failed: %v", err)
	}

	jobs, err := GetAllJobs("new")
	if err != nil {
		t.Fatalf("GetAllJobs failed: %v", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(jobs))
	}
	if jobs[0].Title != "Software Engineer" {
		t.Errorf("wrong title: %s", jobs[0].Title)
	}
	if jobs[0].MatchScore != 85.5 {
		t.Errorf("wrong match score: %f", jobs[0].MatchScore)
	}
}

func TestGetAllJobs_AllStatus(t *testing.T) {
	setupTestDB(t)

	jobs := []*Job{
		{ID: "j1", Title: "Job 1", Company: "A", URL: "http://a.com", Status: "new", MatchScore: 90},
		{ID: "j2", Title: "Job 2", Company: "B", URL: "http://b.com", Status: "saved", MatchScore: 80},
		{ID: "j3", Title: "Job 3", Company: "C", URL: "http://c.com", Status: "applied", MatchScore: 70},
	}
	for _, j := range jobs {
		InsertJob(j)
	}

	all, err := GetAllJobs("all")
	if err != nil {
		t.Fatalf("GetAllJobs(all) failed: %v", err)
	}
	if len(all) != 3 {
		t.Errorf("expected 3 jobs, got %d", len(all))
	}

	saved, err := GetAllJobs("saved")
	if err != nil {
		t.Fatalf("GetAllJobs(saved) failed: %v", err)
	}
	if len(saved) != 1 {
		t.Errorf("expected 1 saved job, got %d", len(saved))
	}
}

func TestUpdateJobStatus(t *testing.T) {
	setupTestDB(t)

	InsertJob(&Job{ID: "j1", Title: "Test", Company: "Co", URL: "http://co.com", Status: "new"})

	if err := UpdateJobStatus("j1", "applied"); err != nil {
		t.Fatalf("UpdateJobStatus failed: %v", err)
	}

	jobs, _ := GetAllJobs("applied")
	if len(jobs) != 1 {
		t.Errorf("expected 1 applied job, got %d", len(jobs))
	}
	if jobs[0].AppliedDate == "" {
		t.Error("applied_date should be set automatically")
	}
}

func TestUpdateJobNotes(t *testing.T) {
	setupTestDB(t)

	InsertJob(&Job{ID: "j1", Title: "Test", Company: "Co", URL: "http://co.com"})

	if err := UpdateJobNotes("j1", "Great opportunity"); err != nil {
		t.Fatalf("UpdateJobNotes failed: %v", err)
	}

	jobs, _ := GetAllJobs("all")
	if jobs[0].Notes != "Great opportunity" {
		t.Errorf("wrong notes: %s", jobs[0].Notes)
	}
}

func TestJobExists(t *testing.T) {
	setupTestDB(t)

	if JobExists("nonexistent") {
		t.Error("JobExists should return false for nonexistent job")
	}

	InsertJob(&Job{ID: "j1", Title: "Test", Company: "Co", URL: "http://co.com"})

	if !JobExists("j1") {
		t.Error("JobExists should return true for inserted job")
	}
}

func TestGetJobCounts(t *testing.T) {
	setupTestDB(t)

	InsertJob(&Job{ID: "j1", Title: "A", Company: "C1", URL: "http://a.com", Status: "new"})
	InsertJob(&Job{ID: "j2", Title: "B", Company: "C2", URL: "http://b.com", Status: "new"})
	InsertJob(&Job{ID: "j3", Title: "C", Company: "C3", URL: "http://c.com", Status: "saved"})

	counts, err := GetJobCounts()
	if err != nil {
		t.Fatalf("GetJobCounts failed: %v", err)
	}
	if counts["new"] != 2 {
		t.Errorf("expected 2 new, got %d", counts["new"])
	}
	if counts["saved"] != 1 {
		t.Errorf("expected 1 saved, got %d", counts["saved"])
	}
}

func TestSaveAndLoadKeywords(t *testing.T) {
	setupTestDB(t)

	kws := []string{"Go", "Docker", "Kubernetes"}
	if err := SaveKeywords("resume.pdf", kws); err != nil {
		t.Fatalf("SaveKeywords failed: %v", err)
	}

	loaded, err := LoadKeywords("resume.pdf")
	if err != nil {
		t.Fatalf("LoadKeywords failed: %v", err)
	}
	if len(loaded) != 3 {
		t.Errorf("expected 3 keywords, got %d", len(loaded))
	}

	// Overwrite with new keywords
	newKws := []string{"Python", "Rust"}
	SaveKeywords("resume.pdf", newKws)
	loaded, _ = LoadKeywords("resume.pdf")
	if len(loaded) != 2 {
		t.Errorf("expected 2 keywords after overwrite, got %d", len(loaded))
	}
}

func TestSaveAndLoadPositions(t *testing.T) {
	setupTestDB(t)

	titles := []string{"Director of Engineering", "Staff Platform Engineer"}
	if err := SavePositions(titles); err != nil {
		t.Fatalf("SavePositions failed: %v", err)
	}

	loaded, err := LoadPositions()
	if err != nil {
		t.Fatalf("LoadPositions failed: %v", err)
	}
	if len(loaded) != 2 {
		t.Errorf("expected 2 positions, got %d", len(loaded))
	}
}

func TestAddAndRemovePosition(t *testing.T) {
	setupTestDB(t)

	AddPosition("DevOps Engineer")
	AddPosition("SRE")

	positions, _ := LoadPositions()
	if len(positions) < 2 {
		t.Errorf("expected at least 2 positions, got %d", len(positions))
	}

	RemovePosition("SRE")
	positions, _ = LoadPositions()
	if len(positions) != 1 {
		t.Errorf("expected 1 position after remove, got %d", len(positions))
	}
}

func TestConfigGetSet(t *testing.T) {
	setupTestDB(t)

	_, err := GetConfig("nonexistent")
	if err == nil {
		t.Error("GetConfig should error for nonexistent key")
	}

	SetConfig("location", "Remote")
	val, err := GetConfig("location")
	if err != nil {
		t.Fatalf("GetConfig failed: %v", err)
	}
	if val != "Remote" {
		t.Errorf("expected 'Remote', got '%s'", val)
	}

	// Overwrite
	SetConfig("location", "San Francisco")
	val, _ = GetConfig("location")
	if val != "San Francisco" {
		t.Errorf("expected 'San Francisco', got '%s'", val)
	}
}

func TestDeleteAndClearJobs(t *testing.T) {
	setupTestDB(t)

	InsertJob(&Job{ID: "j1", Title: "A", Company: "C1", URL: "http://a.com"})
	InsertJob(&Job{ID: "j2", Title: "B", Company: "C2", URL: "http://b.com"})

	DeleteJob("j1")
	jobs, _ := GetAllJobs("all")
	if len(jobs) != 1 {
		t.Errorf("expected 1 job after delete, got %d", len(jobs))
	}

	n, _ := ClearAllJobs()
	if n != 1 {
		t.Errorf("expected 1 cleared, got %d", n)
	}
	jobs, _ = GetAllJobs("all")
	if len(jobs) != 0 {
		t.Errorf("expected 0 jobs after clear, got %d", len(jobs))
	}
}

func TestUpdateOfferDetails(t *testing.T) {
	setupTestDB(t)

	InsertJob(&Job{ID: "j1", Title: "Test", Company: "Co", URL: "http://co.com", Status: "offer"})

	UpdateOfferDetails("j1", "$200k", "Health, 401k", "0.5%")
	jobs, _ := GetAllJobs("offer")
	if len(jobs) != 1 {
		t.Fatal("expected 1 offer job")
	}
	if jobs[0].OfferSalary != "$200k" {
		t.Errorf("wrong salary: %s", jobs[0].OfferSalary)
	}
	if jobs[0].OfferBenefits != "Health, 401k" {
		t.Errorf("wrong benefits: %s", jobs[0].OfferBenefits)
	}
	if jobs[0].OfferEquity != "0.5%" {
		t.Errorf("wrong equity: %s", jobs[0].OfferEquity)
	}
}

func TestInsertJob_Defaults(t *testing.T) {
	setupTestDB(t)

	job := &Job{ID: "j1", Title: "Test", Company: "Co", URL: "http://co.com"}
	if err := InsertJob(job); err != nil {
		t.Fatalf("InsertJob failed: %v", err)
	}
	if job.Status != "new" {
		t.Errorf("default status should be 'new', got '%s'", job.Status)
	}
	if job.CreatedAt == "" {
		t.Error("CreatedAt should be auto-set")
	}
	if job.UpdatedAt == "" {
		t.Error("UpdatedAt should be auto-set")
	}
	if job.ScrapedAt == "" {
		t.Error("ScrapedAt should be auto-set")
	}
}

func TestUpdateJobStatus_DateTracking(t *testing.T) {
	setupTestDB(t)

	InsertJob(&Job{ID: "j1", Title: "Test", Company: "Co", URL: "http://co.com"})

	// applied → should set applied_date
	UpdateJobStatus("j1", "applied")
	jobs, _ := GetAllJobs("all")
	if jobs[0].AppliedDate == "" {
		t.Error("applied_date should be set when status goes to applied")
	}

	// interviewing → should set interview_date
	UpdateJobStatus("j1", "interviewing")
	jobs, _ = GetAllJobs("all")
	if jobs[0].InterviewDate == "" {
		t.Error("interview_date should be set when status goes to interviewing")
	}

	// offer → should set offer_date
	UpdateJobStatus("j1", "offer")
	jobs, _ = GetAllJobs("all")
	if jobs[0].OfferDate == "" {
		t.Error("offer_date should be set when status goes to offer")
	}
}

func TestUpdateJobSkills(t *testing.T) {
	setupTestDB(t)

	InsertJob(&Job{ID: "j1", Title: "Test", Company: "Co", URL: "http://co.com", Skills: "old"})
	UpdateJobSkills("j1", "Go, Python, Rust")
	jobs, _ := GetAllJobs("all")
	if jobs[0].Skills != "Go, Python, Rust" {
		t.Errorf("skills not updated: %s", jobs[0].Skills)
	}
}

func TestUpdateMatchScore(t *testing.T) {
	setupTestDB(t)

	InsertJob(&Job{ID: "j1", Title: "Test", Company: "Co", URL: "http://co.com", MatchScore: 50})
	UpdateMatchScore("j1", 95.5)
	jobs, _ := GetAllJobs("all")
	if jobs[0].MatchScore != 95.5 {
		t.Errorf("match_score not updated: %f", jobs[0].MatchScore)
	}
}

func TestUpdateGapAnalysis(t *testing.T) {
	setupTestDB(t)

	InsertJob(&Job{ID: "j1", Title: "Test", Company: "Co", URL: "http://co.com"})
	UpdateGapAnalysis("j1", "Missing: Kubernetes, Terraform")
	jobs, _ := GetAllJobs("all")
	if jobs[0].GapAnalysis != "Missing: Kubernetes, Terraform" {
		t.Errorf("gap_analysis not updated: %s", jobs[0].GapAnalysis)
	}
}

func TestGetAllJobs_Empty(t *testing.T) {
	setupTestDB(t)

	jobs, err := GetAllJobs("new")
	if err != nil {
		t.Fatalf("GetAllJobs on empty DB failed: %v", err)
	}
	if len(jobs) != 0 {
		t.Errorf("expected 0 jobs, got %d", len(jobs))
	}
}

func TestGetJobCounts_Empty(t *testing.T) {
	setupTestDB(t)

	counts, err := GetJobCounts()
	if err != nil {
		t.Fatalf("GetJobCounts on empty DB failed: %v", err)
	}
	for status, count := range counts {
		if count != 0 {
			t.Errorf("expected 0 for status %s, got %d", status, count)
		}
	}
}

func TestLoadKeywords_Empty(t *testing.T) {
	setupTestDB(t)

	kws, err := LoadKeywords("nonexistent.pdf")
	if err != nil {
		t.Fatalf("LoadKeywords on empty DB failed: %v", err)
	}
	if len(kws) != 0 {
		t.Errorf("expected 0 keywords, got %d", len(kws))
	}
}

func TestLoadPositions_Empty(t *testing.T) {
	setupTestDB(t)

	positions, err := LoadPositions()
	if err != nil {
		t.Fatalf("LoadPositions on empty DB failed: %v", err)
	}
	if len(positions) != 0 {
		t.Errorf("expected 0 positions, got %d", len(positions))
	}
}

func TestUpdateJobStatus_Nonexistent(t *testing.T) {
	setupTestDB(t)

	// Should not error on nonexistent job (just no-op)
	err := UpdateJobStatus("nonexistent", "applied")
	if err != nil {
		t.Errorf("UpdateJobStatus on nonexistent job should not error: %v", err)
	}
}

func TestDeleteJob_Nonexistent(t *testing.T) {
	setupTestDB(t)

	// Should not error on nonexistent job
	err := DeleteJob("nonexistent")
	if err != nil {
		t.Errorf("DeleteJob on nonexistent job should not error: %v", err)
	}
}

func TestClearAllJobs_Empty(t *testing.T) {
	setupTestDB(t)

	n, err := ClearAllJobs()
	if err != nil {
		t.Fatalf("ClearAllJobs on empty DB failed: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 cleared, got %d", n)
	}
}

func TestAddPosition_Duplicate(t *testing.T) {
	setupTestDB(t)

	AddPosition("SRE")
	AddPosition("SRE") // should be ignored (INSERT OR IGNORE)

	positions, _ := LoadPositions()
	if len(positions) != 1 {
		t.Errorf("duplicate position should be ignored, got %d", len(positions))
	}
}
