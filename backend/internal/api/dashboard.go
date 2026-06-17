package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// registerDashboard 提供首页所需聚合视图。
func registerDashboard(g *gin.RouterGroup, d *Deps) {
	g.GET("/dashboard/summary", func(c *gin.Context) { dashboardSummary(c, d) })
	g.GET("/dashboard/balance-trend", func(c *gin.Context) { dashboardBalanceTrend(c, d) })
}

type dashboardLowest struct {
	ChannelID uint     `json:"channel_id"`
	Name      string   `json:"name"`
	Balance   *float64 `json:"balance"`
}

type dashboardChannelStat struct {
	ID                   uint     `json:"id"`
	Name                 string   `json:"name"`
	Type                 string   `json:"type"`
	MonitorEnabled       bool     `json:"monitor_enabled"`
	RechargeRatio        float64  `json:"recharge_ratio"`
	RechargeURL          string   `json:"recharge_url,omitempty"`
	LastBalance          *float64 `json:"last_balance,omitempty"`
	LastError            string   `json:"last_error,omitempty"`
	LastTodayConsumption *float64 `json:"last_today_consumption,omitempty"`
	LastTotalConsumption *float64 `json:"last_total_consumption,omitempty"`
	LastConsumptionAt    *string  `json:"last_consumption_at,omitempty"`
}

func effectiveRechargeRatio(v float64) float64 {
	if v <= 0 {
		return 1
	}
	return v
}

func dashboardSummary(c *gin.Context, d *Deps) {
	channels, err := d.Channels.List()
	if err != nil {
		fail(c, http.StatusInternalServerError, err)
		return
	}

	stats := make([]dashboardChannelStat, 0, len(channels))
	var totalBalance float64
	var lowest *dashboardLowest
	var activeCount, failedCount int

	for _, ch := range channels {
		rechargeRatio := effectiveRechargeRatio(ch.RechargeRatio)
		stat := dashboardChannelStat{
			ID:             ch.ID,
			Name:           ch.Name,
			Type:           string(ch.Type),
			MonitorEnabled: ch.MonitorEnabled,
			RechargeRatio:  rechargeRatio,
			RechargeURL:    ch.RechargeURL,
			LastError:      ch.LastError,
		}
		if ch.LastBalance != nil {
			converted := *ch.LastBalance / rechargeRatio
			stat.LastBalance = &converted
		}
		if ch.LastTodayConsumption != nil {
			converted := *ch.LastTodayConsumption / rechargeRatio
			stat.LastTodayConsumption = &converted
		}
		if ch.LastTotalConsumption != nil {
			converted := *ch.LastTotalConsumption / rechargeRatio
			stat.LastTotalConsumption = &converted
		}
		if ch.LastConsumptionAt != nil {
			formatted := ch.LastConsumptionAt.Format(time.RFC3339)
			stat.LastConsumptionAt = &formatted
		}
		stats = append(stats, stat)
		if ch.LastError != "" {
			failedCount++
		} else if ch.MonitorEnabled {
			activeCount++
		}
		if ch.LastBalance != nil {
			bal := *ch.LastBalance / rechargeRatio
			totalBalance += bal
			if lowest == nil || (lowest.Balance == nil) || (bal < *lowest.Balance) {
				lowest = &dashboardLowest{ChannelID: ch.ID, Name: ch.Name, Balance: &bal}
			}
		}
	}

	recentChanges, err := d.Rates.ListChanges(0, 10)
	if err != nil {
		fail(c, http.StatusInternalServerError, err)
		return
	}
	recentNotifs, err := d.Notifies.ListLogs(10)
	if err != nil {
		fail(c, http.StatusInternalServerError, err)
		return
	}
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		loc = time.Local
	}
	now := time.Now().In(loc)
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	consumption, err := d.Rates.AggregateConsumption(todayStart)
	if err != nil {
		fail(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"total_channels":           len(channels),
			"active_channels":          activeCount,
			"failed_channels":          failedCount,
			"total_balance":            totalBalance,
			"today_consumption":        consumption.Today,
			"total_consumption":        consumption.Total,
			"lowest_balance":           lowest,
			"channels":                 stats,
			"recent_rate_changes":      recentChanges,
			"recent_notification_logs": recentNotifs,
		},
	})
}

func dashboardBalanceTrend(c *gin.Context, d *Deps) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))
	if days <= 0 {
		days = 7
	}
	trend, err := d.Rates.AggregateBalanceTrend(days)
	if err != nil {
		fail(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": trend})
}
