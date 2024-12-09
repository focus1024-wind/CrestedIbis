package site

import (
	"CrestedIbis/src/global"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ExpandSitePreload Preload 自引用预加载，级联查询
func ExpandSitePreload(db *gorm.DB) *gorm.DB {
	return db.Preload(clause.Associations, ExpandSitePreload)
}

func SelectSiteById(id int64) (site Site, err error) {
	err = global.Db.Model(&Site{}).Where(&Site{
		Id: id,
	}).First(&site).Error
	return
}

func selectSites(pid *int64, keywords string) (sites []Site, err error) {
	if pid == nil && keywords == "" {
		// Gorm 默认过滤空值，所以采用 Level 来搜索 pid 为 NULL 的情况
		// 采用级联查询，Preload默认只能搜索1层
		err = global.Db.Debug().Model(&Site{}).Preload(clause.Associations, ExpandSitePreload).Where(&Site{
			Level: 1,
		}).Where("name LIKE ?", "%"+keywords+"%").Find(&sites).Error
	} else if pid == nil {
		err = global.Db.Debug().Model(&Site{}).Preload(clause.Associations, ExpandSitePreload).Where("name LIKE ?", "%"+keywords+"%").Find(&sites).Error
	} else {
		err = global.Db.Debug().Model(&Site{}).Preload(clause.Associations, ExpandSitePreload).Where(&Site{
			Pid: pid,
		}).Where("name LIKE ?", "%"+keywords+"%").Find(&sites).Error
	}
	return
}

func updateSiteName(id int64, name string) (err error) {
	return global.Db.Model(&Site{}).Where(&Site{
		Id: id,
	}).Updates(map[string]interface{}{"name": name}).Error
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

func deleteSite(id int64) (err error) {
	return global.Db.Model(&Site{}).Where(&Site{
		Id: id,
	}).Delete(&Site{}).Error
}

func deleteSites(ids []int64) (err error) {
	fmt.Println(ids)
	return global.Db.Model(&Site{}).Delete(&Site{}, ids).Error
}
