package storage

import (
	"time"

	"gorm.io/gorm"
)

// Channels 渠道仓库。
type Channels struct{ db *gorm.DB }

func NewChannels(db *gorm.DB) *Channels { return &Channels{db: db} }

func (r *Channels) Create(c *Channel) error { return r.db.Create(c).Error }
func (r *Channels) Update(c *Channel) error { return r.db.Save(c).Error }
func (r *Channels) Delete(id uint) error    { return r.db.Delete(&Channel{}, id).Error }
func (r *Channels) FindByID(id uint) (*Channel, error) {
	var c Channel
	if err := r.db.First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}
func (r *Channels) List() ([]Channel, error) {
	var list []Channel
	if err := r.db.Order("id ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
func (r *Channels) ListMonitorEnabled() ([]Channel, error) {
	var list []Channel
	if err := r.db.Where("monitor_enabled = ?", true).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
func (r *Channels) ListDueForRefresh(now time.Time) ([]Channel, error) {
	list, err := r.ListMonitorEnabled()
	if err != nil {
		return nil, err
	}
	due := make([]Channel, 0, len(list))
	for i := range list {
		c := list[i]
		if c.LastBalanceAt == nil {
			due = append(due, c)
			continue
		}
		minutes := c.RefreshInterval
		if minutes <= 0 {
			minutes = 1
		}
		nextAt := c.LastBalanceAt.Add(time.Duration(minutes) * time.Minute)
		if !nextAt.After(now) {
			due = append(due, c)
		}
	}
	return due, nil
}
func (r *Channels) UpdateBalance(id uint, balance float64, at any, lastErr string) error {
	return r.db.Model(&Channel{}).Where("id = ?", id).Updates(map[string]any{
		"last_balance":    balance,
		"last_balance_at": at,
		"last_error":      lastErr,
	}).Error
}
func (r *Channels) UpdateConsumption(id uint, today, total float64, at time.Time) error {
	if at.IsZero() {
		at = time.Now()
	}
	return r.db.Model(&Channel{}).Where("id = ?", id).Updates(map[string]any{
		"last_today_consumption": today,
		"last_total_consumption": total,
		"last_consumption_at":    at,
	}).Error
}
func (r *Channels) SetLastError(id uint, msg string) error {
	return r.db.Model(&Channel{}).Where("id = ?", id).Update("last_error", msg).Error
}
