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
		roleIds[role.RoleId] = true
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
		if !roleIds[role.RoleId] {
			// admin特殊权限，不允许删除
			if sysUser.Username == "admin" && role.RoleName == "admin" {
				continue
			}
			err = global.Db.Model(&sysUser).Association("RoleGroups").Delete(&RoleGroup{
				RoleId: role.RoleId,
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
		return global.Db.Where(&SysUser{
			IDModel: idModel,
		}).Delete(&SysUser{}).Error
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

	return global.Db.Model(&SysUser{}).Delete(&SysUser{}, idsModel.IDs).Error
}

func checkRole(id int64) (err error) {
	var role RoleGroup
	err = global.Db.Model(&RoleGroup{}).Where(&RoleGroup{
		RoleId: id,
	}).First(&role).Error
	if err == nil {
		if role.RoleName == "admin" || role.RoleName == "guest" {
			return errors.New("不允许对admin、guest权限组进行修改")
		}
	}
	return
}

func getRoleById(id int64) (role RoleGroup, err error) {
	err = global.Db.Model(&RoleGroup{}).Where(&RoleGroup{
		RoleId: id,
	}).First(&role).Error
	return
}

func selectAllRoles(keywords string) (roles []RoleGroup, err error) {
	if keywords == "" {
		err = global.Db.Model(&RoleGroup{}).Find(&roles).Error
	} else {
		err = global.Db.Model(&RoleGroup{}).Where("role_name LIKE ?", "%"+keywords+"%").Find(&roles).Error
	}

	return
}

func updateRole(roleId int64, roleName string) (err error) {
	return global.Db.Model(&RoleGroup{}).Where(&RoleGroup{
		RoleId: roleId,
	}).Update("role_name", roleName).Error
}

func insertRole(roleName string) (err error) {
	return global.Db.Model(&RoleGroup{}).Create(&RoleGroup{
		RoleName: roleName,
	}).Error
}

func deleteRole(roleId int64) (err error) {
	return global.Db.Model(&RoleGroup{}).Where(&RoleGroup{
		RoleId: roleId,
	}).Delete(&RoleGroup{}).Error
}

func deleteRoles(ids []int64) (err error) {
	return global.Db.Model(&RoleGroup{}).Delete(&RoleGroup{}, ids).Error
}

func getCasbinRuleByName(name string) (rules []utils.CasbinRule, err error) {
	err = global.Db.Debug().Model(&utils.CasbinRule{}).Where(&utils.CasbinRule{
		RoleKey: name,
	}).Or(&utils.CasbinRule{
		RoleKey: "guest",
	}).Find(&rules).Error
	return
}

func updateRoleRules(roleName string, rules []utils.CasbinRule) (err error) {
	if roleName == "admin" || roleName == "guest" {
		return errors.New("不允许对admin、guest权限组进行修改")
	} else {
		err = global.Db.Model(&utils.CasbinRule{}).Where(&utils.CasbinRule{
			RoleKey: roleName,
		}).Delete(&utils.CasbinRule{}).Error

		if len(rules) == 0 {
			return
		}

		for i := range rules {
			rules[i].Ptype = "p"
			rules[i].RoleKey = roleName
		}

		err = global.Db.Model(&utils.CasbinRule{}).Create(&rules).Error
		return
	}
}
