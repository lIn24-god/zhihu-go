package dao

import (
	"context"
	"zhihu-go/internal/model"

	"gorm.io/gorm"
)

type TimelineDAO interface {
	// BatchInsert 批量插入 timeline 记录
	BatchInsert(ctx context.Context, timelines []*model.Timeline) error
	// GetUserTimeline 分页获取用户的 timeline，按时间倒序
	GetUserTimeline(ctx context.Context, userID uint, limit, offset int) ([]model.Timeline, error)
	// CountUserTimeline 获取用户 timeline 总记录数（用于分页）
	CountUserTimeline(ctx context.Context, userID uint) (int64, error)
}

type timelineDAO struct {
	db *gorm.DB
}

func NewTimelineDAO(db *gorm.DB) TimelineDAO {
	return &timelineDAO{db: db}
}

func (dao *timelineDAO) BatchInsert(ctx context.Context, timelines []*model.Timeline) error {
	if len(timelines) == 0 {
		return nil
	}
	// 批量插入，每批 100 条，避免单次 SQL 过大
	return dao.db.WithContext(ctx).CreateInBatches(timelines, 100).Error
}

func (dao *timelineDAO) GetUserTimeline(ctx context.Context, userID uint, limit, offset int) ([]model.Timeline, error) {
	var list []model.Timeline
	err := dao.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&list).Error
	return list, err
}

func (dao *timelineDAO) CountUserTimeline(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := dao.db.WithContext(ctx).
		Model(&model.Timeline{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}
