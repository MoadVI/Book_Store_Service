package reports

import (
	"Book-Store/internal/models"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type ReportStore struct {
	mu              sync.RWMutex
	reports         []models.SalesReport
	outputDirectory string
}

func NewReportStore(outputDir string) *ReportStore {
	return &ReportStore{
		reports:         make([]models.SalesReport, 0),
		outputDirectory: outputDir,
	}
}

func (rs *ReportStore) SaveReport(report *models.SalesReport) error {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	rs.reports = append(rs.reports, *report)

	if rs.outputDirectory == "" {
		return nil
	}

	if err := os.MkdirAll(rs.outputDirectory, 0755); err != nil {
		return err
	}

	filename := fmt.Sprintf("sales_report_%s.json", report.Timestamp.Format("2006-01-02_15-04-05"))
	filepath := filepath.Join(rs.outputDirectory, filename)

	data, err := json.MarshalIndent(report, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath, data, 0644)
}
