package site

import (
	"CrestedIbis/src/global"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// expandSitePreload Preload 自引用预加载，级联查询
func expandSitePreload(db *gorm.DB) *gorm.DB {
	return db.Preload(clause.Associations, expandSitePreload)
}

func selectSites(pid *int64) (sites []Site, err error) {
	if pid == nil {
		// Gorm 默认过滤空值，所以采用 Level 来搜索 pid 为 NULL 的情况
		// 采用级联查询，Preload默认只能搜索1层
		err = global.Db.Debug().Model(&Site{}).Preload(clause.Associations, expandSitePreload).Where(&Site{
			Level: 1,
		}).Find(&sites).Error
	} else {
		err = global.Db.Debug().Model(&Site{}).Preload(clause.Associations, expandSitePreload).Where(&Site{
			Pid: pid,
		}).Find(&sites).Error
	}
	return
}

func insertSite(site Site) (err error) {
	if site.Pid == nil || *site.Pid == 0 {
		site.Pid = nil
		site.Level = 1
	} else {
		var parentSite Site
		err = global.Db.Model(&Site{}).Where(&Site{
			Id: *site.Pid,
		}).First(&parentSite).Error
		if err != nil {
			return errors.New("查询父级区域失败")
		}
		site.Level = parentSite.Level + 1
	}
	err = global.Db.Model(&Site{}).Create(&site).Error
	if err != nil {
		return errors.New("新建区域失败")
	}
	return
}
