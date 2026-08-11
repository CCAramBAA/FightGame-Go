package cron

import (
	"log"
	"sync"
	"time"

	"gorm.io/gorm"
)

// TaskManager 定时任务管理器
type TaskManager struct {
	db     *gorm.DB
	stopCh chan struct{}
	wg     sync.WaitGroup
}

// RoomMgrInterface 房间管理器接口
type RoomMgrInterface interface {
	RemoveIdleRooms(time.Duration) int
}

func NewTaskManager(db *gorm.DB) *TaskManager {
	return &TaskManager{
		db:     db,
		stopCh: make(chan struct{}),
	}
}

// Start 启动所有定时任务
func (t *TaskManager) Start(roomMgr RoomMgrInterface) {
	t.wg.Add(3)

	// 1. 每小时检查一次，凌晨 3 点清理超过 30 天的对局回放
	go func() {
		defer t.wg.Done()
		ticker := time.NewTicker(15 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-t.stopCh:
				return
			case now := <-ticker.C:
				if now.Hour() == 3 {
					t.cleanupOldReplays()
				}
			}
		}
	}()

	// 2. 每日凌晨 4 点提醒备份
	go func() {
		defer t.wg.Done()
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-t.stopCh:
				return
			case now := <-ticker.C:
				if now.Hour() == 4 {
					log.Println("[Cron] 数据库备份提醒：请定期执行 mysqldump 备份")
				}
			}
		}
	}()

	// 3. 每 5 分钟清理闲置超过 30 分钟的房间
	go func() {
		defer t.wg.Done()
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-t.stopCh:
				return
			case <-ticker.C:
				if roomMgr != nil {
					count := roomMgr.RemoveIdleRooms(30 * time.Minute)
					if count > 0 {
						log.Printf("[Cron] 清理了 %d 个闲置房间", count)
					}
				}
			}
		}
	}()
}

// Stop 停止所有定时任务
func (t *TaskManager) Stop() {
	close(t.stopCh)
	t.wg.Wait()
	log.Println("[Cron] 所有定时任务已停止")
}

// cleanupOldReplays 清理超过 30 天的对局回放数据
func (t *TaskManager) cleanupOldReplays() {
	cutoff := time.Now().Add(-30 * 24 * time.Hour)

	// 清空过期回放帧数据
	result := t.db.Table("battle_records").
		Where("created_at < ? AND frame_data IS NOT NULL AND frame_data != ''", cutoff).
		Update("frame_data", "")

	if result.Error != nil {
		log.Printf("[Cron] 清理回放数据失败: %v", result.Error)
		return
	}
	log.Printf("[Cron] 清理了 %d 条过期回放帧数据", result.RowsAffected)
}
