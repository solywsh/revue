package db

import "gorm.io/gorm"

// SetRevueConfigMusic 设置music开启和关闭
func (gb *GormDb) SetRevueConfigMusic(enable bool) {
	gb.DB.Model(&RevueConfig{}).Where(RevueConfig{ID: 1}).Update("music_enable", enable)
}

// SetRevueConfigReply 设置reply开启和关闭
func (gb *GormDb) SetRevueConfigReply(db *gorm.DB, enable bool) {
	gb.DB.Model(&RevueConfig{}).Where(RevueConfig{ID: 1}).Update("reply_enable", enable)
}

// GetRevueConfig 获取RevueConfig配置
func (gb *GormDb) GetRevueConfig(config *RevueConfig) bool {
	res := gb.DB.Where(RevueConfig{ID: 1}).First(&config)
	if res.RowsAffected >= 1 {
		return true
	}
	return false
}
