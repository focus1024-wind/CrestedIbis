package user

import (
	"CrestedIbis/src/global"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

// Login 用户登陆
// return: nil 登陆成功
func (SysUser) Login(userLogin SysUserLogin) (roles []string, err error) {
	var user SysUser
	user.Username = userLogin.Username
	err = global.Db.Where(&user).Preload("RoleGroups").First(&user).Error

	if err != nil {
		global.Logger.Errorf("查询用户失败, err: %s", err.Error())
		return roles, errors.New("查询用户失败")
	} else {
		// 密码校验
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userLogin.Password))
		if err != nil {
			return roles, errors.New("用户名密码错误")
		} else {
			for _, role := range user.RoleGroups {
				roles = append(roles, role.RoleName)
			}
			return roles, nil
		}
	}
}

// Insert 新增用户
func (SysUser) Insert(user SysUser) (err error) {
	var count int64
	global.Db.Model(&SysUser{}).Where(&SysUser{
		SysUserLogin: SysUserLogin{
			Username: user.Username,
		},
	}).Count(&count)

	if count > 0 {
		return errors.New("用户已存在")
	} else {
		// 密码加盐加密
		password, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		user.Password = string(password)
		err = global.Db.Create(&user).Error
		return err
	}
}

func selectUsersByPages(page int64, pageSize int64) (total int64, users []SysUser, err error) {
	db := global.Db.Model(SysUser{})

	if err = db.Count(&total).Error; err != nil {
		return
	}

	offset := (page - 1) * pageSize
	if err = db.Preload("RoleGroups").Offset(int(offset)).Limit(int(pageSize)).Find(&users).Error; err != nil {
		return
	}
	return
}

func updateUserPassword(username string, password string) (err error) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	err = global.Db.Model(&SysUser{}).Where("username = ?", username).Update("password", string(hashedPassword)).Error
	return
}

func deleteUser(username string) (err error) {
	if username == "admin" {
		return errors.New("admin 用户不允许删除")
	}
	return global.Db.Where("username = ?", username).Delete(&SysUser{}).Error
}

func selectAllRoles() (roles []RoleGroup, err error) {
	err = global.Db.Model(&RoleGroup{}).Find(&roles).Error
	return
}
