package reports

import (
	"Book-Store/internal/models"
	"fmt"
	"time"
)

func (rs *ReportStore) GetLatestReport() (*models.SalesReport, error) {
	rs.mu.RLock()
	defer rs.mu.RUnlock()

	if len(rs.reports) == 0 {
		return nil, fmt.Errorf("no reports available")
	}

	return &rs.reports[len(rs.reports)-1], nil
}

func (rs *ReportStore) GetReportsSince(since time.Time) []models.SalesReport {
	rs.mu.RLock()
	defer rs.mu.RUnlock()

	filtered := make([]models.SalesReport, 0)
	for _, report := range rs.reports {
		if report.Timestamp.After(since) {
			filtered = append(filtered, report)
		}
	}
	return filtered
}

func (rs *ReportStore) ListAllReports() []models.SalesReport {
	rs.mu.RLock()
	defer rs.mu.RUnlock()

	return rs.reports
}
