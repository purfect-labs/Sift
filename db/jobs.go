package db

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Init(dbPath string) error {
	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	schema := `
	CREATE TABLE IF NOT EXISTS jobs (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		company TEXT NOT NULL,
		location TEXT DEFAULT '',
		sector TEXT DEFAULT '',
		salary TEXT DEFAULT '',
		description TEXT DEFAULT '',
		skills TEXT DEFAULT '',
		url TEXT NOT NULL,
		source TEXT DEFAULT '',
		match_score REAL DEFAULT 0,
		status TEXT DEFAULT 'new',
		notes TEXT DEFAULT '',
		applied_date TEXT DEFAULT '',
		interview_date TEXT DEFAULT '',
		offer_date TEXT DEFAULT '',
		offer_salary TEXT DEFAULT '',
		offer_benefits TEXT DEFAULT '',
		offer_equity TEXT DEFAULT '',
		gap_analysis TEXT DEFAULT '',
		scraped_at TEXT NOT NULL,
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS resume_keywords (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		resume_path TEXT NOT NULL,
		keyword TEXT NOT NULL,
		created_at TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS config (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS target_positions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL UNIQUE,
		source TEXT DEFAULT 'manual',
		created_at TEXT NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_jobs_status ON jobs(status);
	CREATE INDEX IF NOT EXISTS idx_jobs_match_score ON jobs(match_score);
	CREATE INDEX IF NOT EXISTS idx_jobs_company ON jobs(company);
	CREATE INDEX IF NOT EXISTS idx_jobs_scraped_at ON jobs(scraped_at);
	CREATE INDEX IF NOT EXISTS idx_keywords_resume ON resume_keywords(resume_path);
	`

	_, err = DB.Exec(schema)
	if err != nil {
		return err
	}

	// Migration: add columns if missing from old DB
	migrations := []string{
		`ALTER TABLE jobs ADD COLUMN applied_date TEXT DEFAULT ''`,
		`ALTER TABLE jobs ADD COLUMN interview_date TEXT DEFAULT ''`,
		`ALTER TABLE jobs ADD COLUMN offer_date TEXT DEFAULT ''`,
		`ALTER TABLE jobs ADD COLUMN offer_salary TEXT DEFAULT ''`,
		`ALTER TABLE jobs ADD COLUMN offer_benefits TEXT DEFAULT ''`,
		`ALTER TABLE jobs ADD COLUMN offer_equity TEXT DEFAULT ''`,
		`ALTER TABLE jobs ADD COLUMN gap_analysis TEXT DEFAULT ''`,
	}
	for _, m := range migrations {
		DB.Exec(m) // ignore errors — column might already exist
	}

	return nil
}

type Job struct {
	ID             string  `json:"id"`
	Title          string  `json:"title"`
	Company        string  `json:"company"`
	Location       string  `json:"location"`
	Sector         string  `json:"sector"`
	Salary         string  `json:"salary"`
	Description    string  `json:"description"`
	Skills         string  `json:"skills"`
	URL            string  `json:"url"`
	Source         string  `json:"source"`
	MatchScore     float64 `json:"match_score"`
	Status         string  `json:"status"`
	Notes          string  `json:"notes"`
	AppliedDate    string  `json:"applied_date"`
	InterviewDate  string  `json:"interview_date"`
	OfferDate      string  `json:"offer_date"`
	OfferSalary    string  `json:"offer_salary"`
	OfferBenefits  string  `json:"offer_benefits"`
	OfferEquity    string  `json:"offer_equity"`
	GapAnalysis    string  `json:"gap_analysis"`
	ScrapedAt      string  `json:"scraped_at"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

func InsertJob(job *Job) error {
	now := time.Now().Format(time.RFC3339)
	if job.CreatedAt == "" {
		job.CreatedAt = now
	}
	job.UpdatedAt = now
	if job.ScrapedAt == "" {
		job.ScrapedAt = now
	}
	if job.Status == "" {
		job.Status = "new"
	}

	_, err := DB.Exec(`
		INSERT OR REPLACE INTO jobs (id, title, company, location, sector, salary, description, skills, url, source, match_score, status, notes, applied_date, interview_date, offer_date, offer_salary, offer_benefits, offer_equity, gap_analysis, scraped_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, job.ID, job.Title, job.Company, job.Location, job.Sector, job.Salary,
		job.Description, job.Skills, job.URL, job.Source, job.MatchScore,
		job.Status, job.Notes, job.AppliedDate, job.InterviewDate, job.OfferDate,
		job.OfferSalary, job.OfferBenefits, job.OfferEquity, job.GapAnalysis,
		job.ScrapedAt, job.CreatedAt, job.UpdatedAt)
	return err
}

func GetAllJobs(status string) ([]Job, error) {
	var rows *sql.Rows
	var err error

	if status == "" || status == "all" {
		rows, err = DB.Query(`SELECT * FROM jobs ORDER BY match_score DESC, scraped_at DESC`)
	} else {
		rows, err = DB.Query(`SELECT * FROM jobs WHERE status = ? ORDER BY match_score DESC, scraped_at DESC`, status)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []Job
	for rows.Next() {
		var j Job
		err := rows.Scan(&j.ID, &j.Title, &j.Company, &j.Location, &j.Sector,
			&j.Salary, &j.Description, &j.Skills, &j.URL, &j.Source,
			&j.MatchScore, &j.Status, &j.Notes, &j.AppliedDate, &j.InterviewDate, &j.OfferDate,
			&j.OfferSalary, &j.OfferBenefits, &j.OfferEquity, &j.GapAnalysis,
			&j.ScrapedAt, &j.CreatedAt, &j.UpdatedAt)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}
	return jobs, nil
}

// GetJobStatus returns the current status of a job by ID.
func GetJobStatus(id string) (string, error) {
	var status string
	err := DB.QueryRow(`SELECT status FROM jobs WHERE id = ?`, id).Scan(&status)
	if err != nil {
		return "", err
	}
	return status, nil
}

func UpdateJobStatus(id, status string) error {
	now := time.Now().Format(time.RFC3339)
	_, err := DB.Exec(`UPDATE jobs SET status = ?, updated_at = ? WHERE id = ?`, status, now, id)
	if err == nil {
		// Auto-set date fields when status changes
		switch status {
		case "applied":
			DB.Exec(`UPDATE jobs SET applied_date = ? WHERE id = ? AND applied_date = ''`, now, id)
		case "interviewing":
			DB.Exec(`UPDATE jobs SET interview_date = ? WHERE id = ? AND interview_date = ''`, now, id)
		case "offer":
			DB.Exec(`UPDATE jobs SET offer_date = ? WHERE id = ? AND offer_date = ''`, now, id)
		}
	}
	return err
}

func UpdateJobNotes(id, notes string) error {
	now := time.Now().Format(time.RFC3339)
	_, err := DB.Exec(`UPDATE jobs SET notes = ?, updated_at = ? WHERE id = ?`, notes, now, id)
	return err
}

func UpdateJobSkills(id, skills string) error {
	now := time.Now().Format(time.RFC3339)
	_, err := DB.Exec(`UPDATE jobs SET skills = ?, updated_at = ? WHERE id = ?`, skills, now, id)
	return err
}

func UpdateMatchScore(id string, score float64) error {
	now := time.Now().Format(time.RFC3339)
	_, err := DB.Exec(`UPDATE jobs SET match_score = ?, updated_at = ? WHERE id = ?`, score, now, id)
	return err
}

func UpdateGapAnalysis(id, gaps string) error {
	now := time.Now().Format(time.RFC3339)
	_, err := DB.Exec(`UPDATE jobs SET gap_analysis = ?, updated_at = ? WHERE id = ?`, gaps, now, id)
	return err
}

func GetJobCounts() (map[string]int, error) {
	rows, err := DB.Query(`SELECT status, COUNT(*) FROM jobs GROUP BY status`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := map[string]int{"new": 0, "saved": 0, "applied": 0, "interviewing": 0, "offer": 0, "rejected": 0, "archived": 0, "not_remote": 0}
	for rows.Next() {
		var status string
		var count int
		rows.Scan(&status, &count)
		counts[status] = count
	}
	return counts, nil
}

func JobExists(id string) bool {
	var exists int
	DB.QueryRow(`SELECT COUNT(*) FROM jobs WHERE id = ?`, id).Scan(&exists)
	return exists > 0
}

// SaveKeywords persists extracted keywords to DB, clearing old ones for this resume
func SaveKeywords(resumePath string, keywords []string) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`DELETE FROM resume_keywords WHERE resume_path = ?`, resumePath)
	if err != nil {
		return err
	}

	now := time.Now().Format(time.RFC3339)
	for _, kw := range keywords {
		_, err = tx.Exec(`INSERT INTO resume_keywords (resume_path, keyword, created_at) VALUES (?, ?, ?)`,
			resumePath, kw, now)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

// LoadKeywords retrieves saved keywords for a resume path
func LoadKeywords(resumePath string) ([]string, error) {
	rows, err := DB.Query(`SELECT keyword FROM resume_keywords WHERE resume_path = ? ORDER BY id`, resumePath)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keywords []string
	for rows.Next() {
		var kw string
		rows.Scan(&kw)
		keywords = append(keywords, kw)
	}
	return keywords, nil
}

// SavePositions persists target job titles
func SavePositions(titles []string) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`DELETE FROM target_positions`)
	if err != nil {
		return err
	}

	now := time.Now().Format(time.RFC3339)
	for _, t := range titles {
		_, err = tx.Exec(`INSERT INTO target_positions (title, source, created_at) VALUES (?, 'ai', ?)`, t, now)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

// AddPosition adds a single position (manual entry)
func AddPosition(title string) error {
	now := time.Now().Format(time.RFC3339)
	_, err := DB.Exec(`INSERT OR IGNORE INTO target_positions (title, source, created_at) VALUES (?, 'manual', ?)`, title, now)
	return err
}

// RemovePosition removes a position by title
func RemovePosition(title string) error {
	_, err := DB.Exec(`DELETE FROM target_positions WHERE title = ?`, title)
	return err
}

// LoadPositions retrieves all saved target job titles
func LoadPositions() ([]string, error) {
	rows, err := DB.Query(`SELECT title FROM target_positions ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var titles []string
	for rows.Next() {
		var t string
		rows.Scan(&t)
		titles = append(titles, t)
	}
	return titles, nil
}

// GetConfig retrieves a config value
func GetConfig(key string) (string, error) {
	var value string
	err := DB.QueryRow(`SELECT value FROM config WHERE key = ?`, key).Scan(&value)
	if err != nil {
		return "", err
	}
	return value, nil
}

// SetConfig stores a config key/value pair
func SetConfig(key, value string) error {
	_, err := DB.Exec(`INSERT OR REPLACE INTO config (key, value) VALUES (?, ?)`, key, value)
	return err
}

// DeleteJob removes a single job
func DeleteJob(id string) error {
	_, err := DB.Exec(`DELETE FROM jobs WHERE id = ?`, id)
	return err
}

// ClearAllJobs removes all jobs
func ClearAllJobs() (int, error) {
	result, err := DB.Exec(`DELETE FROM jobs`)
	if err != nil {
		return 0, err
	}
	n, _ := result.RowsAffected()
	return int(n), nil
}

func UpdateOfferDetails(id, salary, benefits, equity string) error {
	now := time.Now().Format(time.RFC3339)
	_, err := DB.Exec(`UPDATE jobs SET offer_salary = ?, offer_benefits = ?, offer_equity = ?, updated_at = ? WHERE id = ?`,
		salary, benefits, equity, now, id)
	return err
}
