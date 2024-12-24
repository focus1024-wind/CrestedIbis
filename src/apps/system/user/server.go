package user

import (
	"CrestedIbis/src/global"
	"CrestedIbis/src/global/model"
	"CrestedIbis/src/utils"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

// UpdatePassword 修改用户密码
func (SysUser) UpdatePassword(sysUserLogin SysUserLogin) (err error) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(sysUserLogin.Password), bcrypt.DefaultCost)
	err = global.Db.Model(&SysUser{}).Where("username = ?", sysUserLogin.Username).Update("password", string(hashedPassword)).Error
	return
}

// Update 更新用户
func (SysUser) Update(sysUser SysUser) (err error) {
	// 记录更新用户的权限ID信息，用户删除权限组
	roleIds := make(map[int64]bool)
	for _, role := range sysUser.RoleGroups {
		// map默认为false，所以用true记录
		roleIds[role.ID] = true
	}

	if sysUser.Password != "" {
		// 若密码不为空，对密码进行加密
		// 密码加盐加密
		password, _ := bcrypt.GenerateFromPassword([]byte(sysUser.Password), bcrypt.DefaultCost)
		sysUser.Password = string(password)
	}

	// 更新数据
	// 在更新时，会自动添加新增的role信息，所以后续只需要删除对应的role即可
	// gorm在更新时不会更新0值，所以强制sex性别进行更新
	err = global.Db.Updates(&sysUser).Error

	err = global.Db.Where(&SysUser{
		IDModel: model.IDModel{
			ID: sysUser.ID,
		},
	}).Preload("RoleGroups").First(&sysUser).Error

	// 删除对应外键关系
	for _, role := range sysUser.RoleGroups {
		if !roleIds[role.ID] {
			// admin特殊权限，不允许删除
			if sysUser.Username == "admin" && role.RoleName == "admin" {
				continue
			}
			err = global.Db.Model(&sysUser).Association("RoleGroups").Delete(&RoleGroup{
				IDModel: model.IDModel{
					ID: role.ID,
				},
			})
		}
	}
	return
}

// Delete 删除用户
func (SysUser) Delete(idModel model.IDModel) (err error) {
	var user SysUser

	if err = global.Db.Where(&SysUser{
		IDModel: idModel,
	}).First(&user).Error; err != nil {
		return err
	}

	if user.Username == "admin" {
		return errors.New("admin 用户不允许删除")
	} else {
		return global.Db.Delete(&SysUser{}, idModel.ID).Error
	}
}

// Login 用户登陆
func (SysUser) Login(sysUserLogin SysUserLogin) (sysUser SysUser, err error) {
	sysUser.Username = sysUserLogin.Username

	if err = global.Db.Model(&sysUser).Preload("RoleGroups").First(&sysUser).Error; err != nil {
		global.Logger.Errorf("[DataBase] 数据库搜索查询用户失败: %s", err.Error())
		return sysUser, errors.New("查询用户失败")
	} else {
		// 密码校验
		if err = bcrypt.CompareHashAndPassword([]byte(sysUser.Password), []byte(sysUserLogin.Password)); err != nil {
			global.Logger.Errorf("用户名或密码错误: %s", err.Error())
			return sysUser, errors.New("用户名或密码错误")
		} else {
			return sysUser, nil
		}
	}
}

// Insert 新增用户
func (SysUser) Insert(sysUser SysUser) (err error) {
	var count int64
	global.Db.Model(&SysUser{}).Where(&SysUser{
		SysUserLogin: SysUserLogin{
			Username: sysUser.Username,
		},
	}).Count(&count)

	if count > 0 {
		return errors.New("用户已存在")
	} else {
		// 密码加盐加密
		password, _ := bcrypt.GenerateFromPassword([]byte(sysUser.Password), bcrypt.DefaultCost)
		sysUser.Password = string(password)
		err = global.Db.Create(&sysUser).Error
		return err
	}
}

// SelectUsers 根据查询条件获取用户列表
func (SysUser) SelectUsers(page int64, pageSize int64, keywords string) (total int64, sysUsers []SysUser, err error) {
	db := global.Db.Model(SysUser{}).Preload("RoleGroups")

	if keywords != "" {
		db = db.Where("username LIKE ?", "%"+keywords+"%").
			Or("nickname LIKE ?", "%"+keywords+"%").
			Or("email LIKE ?", "%"+keywords+"%").
			Or("phone LIKE ?", "%"+keywords+"%")
	}

	if err = db.Count(&total).Error; err != nil {
		return
	}

	offset := (page - 1) * pageSize

	err = db.Offset(int(offset)).Limit(int(pageSize)).Find(&sysUsers).Error

	return
}

// Deletes 批量删除用户
func (SysUser) Deletes(idsModel model.IDsModel) (err error) {
	var user SysUser
	err = global.Db.Where(&SysUser{
		SysUserLogin: SysUserLogin{
			Username: "admin",
		},
	}).First(&user).Error

	for i := range idsModel.IDs {
		if idsModel.IDs[i] == user.ID {
			idsModel.IDs = append(idsModel.IDs[:i], idsModel.IDs[i+1:]...)
			break
		}
	}

	return global.Db.Delete(&SysUser{}, idsModel.IDs).Error
}

// check 校验是否为特殊权限组
func (RoleGroup) check(id int64) (err error) {
	var role RoleGroup
	err = global.Db.Model(&RoleGroup{}).Where(&RoleGroup{
		IDModel: model.IDModel{
			ID: id,
		},
	}).First(&role).Error
	if err == nil {
		if role.RoleName == "admin" || role.RoleName == "guest" {
			return errors.New("不允许对admin、guest权限组进行修改")
		}
	}
	return
}

// Insert 新增权限组
func (RoleGroup) Insert(roleName string) (err error) {
	return global.Db.Model(&RoleGroup{}).Create(&RoleGroup{
		RoleName: roleName,
	}).Error
}

// Update 更新权限组
func (RoleGroup) Update(roleGroup RoleGroup) (err error) {
	if err = roleGroup.check(roleGroup.ID); err != nil {
		return
	} else {
		return global.Db.Updates(&roleGroup).Error
	}
}

// Delete 删除权限组
func (RoleGroup) Delete(idModel model.IDModel) (err error) {
	err = RoleGroup{}.check(idModel.ID)
	if err != nil {
		return
	} else {
		return global.Db.Delete(&RoleGroup{}, idModel.ID).Error
	}
}

// Select 根据查询条件，搜索权限组
func (RoleGroup) Select(keywords string) (roles []RoleGroup, err error) {
	db := global.Db.Model(&RoleGroup{})

	if keywords != "" {
		db = db.Where("role_name LIKE ?", "%"+keywords+"%")
	}

	err = db.Find(&roles).Error

	return
}

// Deletes 批量删除权限组
func (RoleGroup) Deletes(idsModel model.IDsModel) (err error) {
	for _, id := range idsModel.IDs {
		err = RoleGroup{}.check(id)
		if err != nil {
			return
		}
	}
	return global.Db.Delete(&RoleGroup{}, idsModel.IDs).Error
}

// SelectRules 根据权限组名称获取对应权限
func (RoleGroup) SelectRules(name string) (rules []utils.CasbinRule, err error) {
	err = global.Db.Model(&utils.CasbinRule{}).Where(&utils.CasbinRule{
		RoleKey: name,
	}).Or(&utils.CasbinRule{
		RoleKey: "guest",
	}).Find(&rules).Error
	return
}

// UpdateRules 更新权限组权限信息
func (RoleGroup) UpdateRules(roleRuleUpdateEntity RoleRuleUpdateEntity) (err error) {
	var role RoleGroup

	// 获取对应权限组
	err = global.Db.Model(&RoleGroup{}).Where(&RoleGroup{
		IDModel: model.IDModel{
			ID: roleRuleUpdateEntity.ID,
		},
	}).First(&role).Error

	// 权限组校验
	if err != nil {
		return
	} else if role.RoleName == "admin" || role.RoleName == "guest" {
		return errors.New("不允许对admin、guest权限组进行修改")
	} else {
		// 删除原权限
		err = global.Db.Model(&utils.CasbinRule{}).Where(&utils.CasbinRule{
			RoleKey: role.RoleName,
		}).Delete(&utils.CasbinRule{}).Error

		if len(roleRuleUpdateEntity.Rules) == 0 {
			return
		}

		// 新增权限信息
		for i := range roleRuleUpdateEntity.Rules {
			roleRuleUpdateEntity.Rules[i].Ptype = "p"
			roleRuleUpdateEntity.Rules[i].RoleKey = role.RoleName
		}

		err = global.Db.Model(&utils.CasbinRule{}).Create(&roleRuleUpdateEntity.Rules).Error
		return
	}
}
