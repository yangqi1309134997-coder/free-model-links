package healthcheck

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/ch3nnn/webstack-go/internal/dal/model"
	"github.com/ch3nnn/webstack-go/pkg/log"
)

type HealthStatus struct {
	Uptime7d     float64 `json:"uptime_7d"`
	AvgLatencyMs int     `json:"avg_latency_ms"`
	LastCheck    string  `json:"last_check"`
	IsOnline     bool    `json:"is_online"`
}

type CheckRecord struct {
	Online    bool
	LatencyMs int
	CheckedAt time.Time
}

type Checker struct {
	db       *gorm.DB
	logger   *log.Logger
	interval time.Duration
	client   *http.Client
	mu       sync.RWMutex
	records  map[int][]CheckRecord
	maxFails int
}

func NewChecker(db *gorm.DB, logger *log.Logger, intervalMinutes int, maxConsecutiveFails int) *Checker {
	if intervalMinutes <= 0 {
		intervalMinutes = 30
	}
	if maxConsecutiveFails <= 0 {
		maxConsecutiveFails = 10
	}

	return &Checker{
		db:       db,
		logger:   logger,
		interval: time.Duration(intervalMinutes) * time.Minute,
		client: &http.Client{
			Timeout: 15 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
			},
		},
		records:  make(map[int][]CheckRecord),
		maxFails: maxConsecutiveFails,
	}
}

func (c *Checker) Start() {
	go func() {
		c.logger.Info("Health checker starting", zap.Duration("interval", c.interval))
		c.runOnce()
		ticker := time.NewTicker(c.interval)
		defer ticker.Stop()
		for range ticker.C {
			// Check if enabled
			var config model.SysConfig
			c.db.First(&config)
			if !config.HealthCheckEnabled {
				continue
			}
			c.runOnce()
		}
	}()
}

func (c *Checker) RunOnce() {
	c.runOnce()
}

func (c *Checker) runOnce() {
	var sites []model.StSite
	c.db.Where("is_used = ? AND deleted_at IS NULL", true).Find(&sites)

	if len(sites) == 0 {
		return
	}

	c.logger.Info("Running health checks", zap.Int("site_count", len(sites)))

	var wg sync.WaitGroup
	sem := make(chan struct{}, 5)

	for _, site := range sites {
		wg.Add(1)
		sem <- struct{}{}
		go func(s model.StSite) {
			defer wg.Done()
			defer func() { <-sem }()
			c.checkSite(s)
		}(site)
	}
	wg.Wait()
}

func (c *Checker) checkSite(site model.StSite) {
	start := time.Now()
	online := false
	latencyMs := 0

	req, err := http.NewRequest("GET", site.URL, nil)
	if err != nil {
		c.logger.Warn("Health check: invalid URL", zap.Int("site_id", site.ID), zap.String("url", site.URL))
		c.recordCheck(site.ID, false, 0)
		c.updateSiteHealth(site.ID)
		return
	}

	req.Header.Set("User-Agent", "WebStack-Go-HealthChecker/1.0")

	resp, err := c.client.Do(req)
	latencyMs = int(time.Since(start).Milliseconds())

	if err != nil {
		c.logger.Debug("Health check failed", zap.Int("site_id", site.ID), zap.Error(err))
	} else {
		defer resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode < 500 {
			online = true
		}
	}

	c.recordCheck(site.ID, online, latencyMs)
	c.updateSiteHealth(site.ID)
	c.checkAutoDelist(site.ID)
}

func (c *Checker) recordCheck(siteID int, online bool, latencyMs int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	record := CheckRecord{
		Online:    online,
		LatencyMs: latencyMs,
		CheckedAt: time.Now(),
	}

	c.records[siteID] = append(c.records[siteID], record)

	cutoff := time.Now().Add(-7 * 24 * time.Hour)
	records := c.records[siteID]
	filtered := records[:0]
	for _, r := range records {
		if r.CheckedAt.After(cutoff) {
			filtered = append(filtered, r)
		}
	}
	c.records[siteID] = filtered
}

func (c *Checker) updateSiteHealth(siteID int) {
	c.mu.RLock()
	records := c.records[siteID]
	c.mu.RUnlock()

	if len(records) == 0 {
		return
	}

	totalChecks := len(records)
	onlineCount := 0
	totalLatency := 0
	latencyCount := 0

	for _, r := range records {
		if r.Online {
			onlineCount++
			totalLatency += r.LatencyMs
			latencyCount++
		}
	}

	uptime := float64(onlineCount) / float64(totalChecks) * 100
	avgLatency := 0
	if latencyCount > 0 {
		avgLatency = totalLatency / latencyCount
	}

	lastRecord := records[len(records)-1]

	status := HealthStatus{
		Uptime7d:     float64(int(uptime*10)) / 10,
		AvgLatencyMs: avgLatency,
		LastCheck:    lastRecord.CheckedAt.Format(time.DateTime),
		IsOnline:     lastRecord.Online,
	}

	statusJSON, _ := json.Marshal(status)
	c.db.Model(&model.StSite{}).Where("id = ?", siteID).Update("health_status", string(statusJSON))
}

func (c *Checker) checkAutoDelist(siteID int) {
	c.mu.RLock()
	records := c.records[siteID]
	c.mu.RUnlock()

	if len(records) < c.maxFails {
		return
	}

	tail := records[len(records)-c.maxFails:]
	allFailed := true
	for _, r := range tail {
		if r.Online {
			allFailed = false
			break
		}
	}

	if allFailed {
		c.logger.Warn("Auto-delisting site due to consecutive failures",
			zap.Int("site_id", siteID),
			zap.Int("consecutive_fails", c.maxFails))
		c.db.Model(&model.StSite{}).Where("id = ?", siteID).Update("is_used", false)
		fmt.Printf("[HealthCheck] Site %d auto-delisted after %d consecutive failures\n", siteID, c.maxFails)
	}
}
