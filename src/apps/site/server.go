package site

import (
	"CrestedIbis/src/global"
	"CrestedIbis/src/global/model"
	"container/list"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Insert 新建区域
func (Site) Insert(site Site) (err error) {
	if site.Pid == nil || *site.Pid == 0 {
		// 无父区域，为顶级区域
		site.Pid = nil
		site.Level = 1
	} else {
		// 存在父区域，关联信息
		var parentSite Site
		err = global.Db.Model(&Site{}).Where(&Site{
			IDModel: model.IDModel{
				ID: *site.Pid,
			},
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

// Update 更新区域
func (Site) Update(site Site) (err error) {
	if *site.Pid == site.ID {
		return
	}

	// 更新Level信息
	if *site.Pid == 0 {
		site.Pid = nil
		site.Level = 1
	} else {
		var parentSite Site

		if err = global.Db.Model(&Site{}).Where(&Site{
			IDModel: model.IDModel{
				ID: *site.Pid,
			},
		}).First(&parentSite).Error; err != nil {
			return err
		} else {
			site.Level = parentSite.Level + 1
		}
	}

	// 获取children信息
	var tempSite Site
	if err = global.Db.Model(&Site{}).Preload(clause.Associations, ExpandSitePreload).Where(&Site{
		IDModel: model.IDModel{
			ID: site.ID,
		},
	}).First(&tempSite).Error; err != nil {
		return
	}
	site.Children = tempSite.Children

	// Site 为 树结构，所以采用层次遍历方式遍历站点
	siteQueue := list.New()
	siteQueue.PushBack(site)
	levelQueue := list.New()
	levelQueue.PushBack(site.Level)
	for siteQueue.Len() > 0 {
		elem := siteQueue.Remove(siteQueue.Front()).(Site)
		level := levelQueue.Remove(levelQueue.Front()).(int8)

		for _, child := range elem.Children {
			siteQueue.PushBack(child)
			levelQueue.PushBack(level + 1)
		}
		elem.Level = level
		global.Db.Select("*").Updates(&elem)
	}
	return
}

// Delete 删除区域
func (Site) Delete(idModel model.IDModel) (err error) {
	return global.Db.Delete(&Site{}, idModel.ID).Error
}

// SelectList 根据查询条件，搜索区域
func (Site) SelectList(pid *int64, keywords string) (sites []Site, err error) {
	db := global.Db.Model(&Site{}).Preload(clause.Associations, ExpandSitePreload)

	if keywords != "" {
		db = db.Where("name LIKE ?", "%"+keywords+"%")
	}

	if pid == nil && keywords != "" {
		// 存在keywords，全局搜索
		err = db.Find(&sites).Error
	} else if pid == nil {
		// 都不存在，从一级level开始搜索
		err = db.Where(&Site{
			Level: 1,
		}).Find(&sites).Error
	} else {
		err = db.Where(&Site{
			Pid: pid,
		}).Find(&sites).Error
	}

	return
}

// Deletes 批量删除区域
func (Site) Deletes(idsModel model.IDsModel) (err error) {
	return global.Db.Delete(&Site{}, idsModel.IDs).Error
}

// ExpandSitePreload Preload 自引用预加载，级联查询
func ExpandSitePreload(db *gorm.DB) *gorm.DB {
	return db.Preload(clause.Associations, ExpandSitePreload)
}
