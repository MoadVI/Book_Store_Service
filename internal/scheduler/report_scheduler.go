package scheduler

import (
	"Book-Store/internal/config"
	"Book-Store/internal/http/middleware"
	"Book-Store/internal/reports"
	"Book-Store/internal/store"
	"context"
	"log"
	"sync"
	"time"
)

type ReportScheduler struct {
	orderStore  store.OrderStore
	reportStore *reports.ReportStore
	metrics     middleware.Metrics
	interval    time.Duration
	ticker      *time.Ticker
	stopChan    chan struct{}
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewReportScheduler(orderStore store.OrderStore, reportStore *reports.ReportStore, metrics middleware.Metrics, interval time.Duration) *ReportScheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &ReportScheduler{
		orderStore:  orderStore,
		reportStore: reportStore,
		metrics:     metrics,
		interval:    interval,
		stopChan:    make(chan struct{}),
		ctx:         ctx,
		cancel:      cancel,
	}
}

func (rs *ReportScheduler) Start() {
	rs.ticker = time.NewTicker(rs.interval)

	rs.wg.Go(func() {
		log.Println("Report scheduler started")

		rs.generateAndSaveReport()

		for {
			select {
			case <-rs.ticker.C:
				rs.generateAndSaveReport()
			case <-rs.stopChan:
				log.Println("Report scheduler stopping...")
				return
			case <-rs.ctx.Done():
				log.Println("Report scheduler context cancelled")
				return
			}
		}
	})
}

func (rs *ReportScheduler) Stop() {
	close(rs.stopChan)
	rs.cancel()
	if rs.ticker != nil {
		rs.ticker.Stop()
	}
	rs.wg.Wait()
	log.Println("Report scheduler stopped")
}

func (rs *ReportScheduler) generateAndSaveReport() {
	log.Println("Generating sales report...")

	report, err := reports.GenerateSalesReport(rs.ctx, rs.orderStore)
	if err != nil {
		log.Printf("Error generating report: %v", err)
		return
	}

	if err := rs.reportStore.SaveReport(report); err != nil {
		log.Printf("Error saving report: %v", err)
		return
	}

	rs.metrics.ResetHits()

	output_dir := config.LoadConfig().ReportOutputDirectory
	log.Printf("Sales report generated successfully at %s", output_dir)
	log.Printf("Total Revenue: $%.2f, Total Orders: %d", report.TotalRevenue, report.TotalOrders)
}
