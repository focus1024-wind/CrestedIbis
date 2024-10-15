package utils

import (
	"CrestedIbis/src/global"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

type CasbinRule struct {
	ID      uint `gorm:"primaryKey;autoIncrement"`
	Ptype   string
	RoleKey string `gorm:"column:v0"`
	Path    string `gorm:"column:v1"`
	Method  string `gorm:"column:v2"`
}

func (CasbinRule) TableName() string {
	// gormadapter 根据`casbin_rule`表名获取Policy，所以必须手动指定表名，避免gorm生成表名不一致
	return "casbin_rule"
}

// CasbinService 获取casbin执行器进行权限检查
func CasbinService() *casbin.Enforcer {
	// https://casbin.org/zh/docs/supported-models
	// https://github.com/casbin/casbin/blob/master/examples/keymatch_model.conf
	// 参考 casbin RESTful模型 keymatch_model.conf模型文件
	casbinModel, err := model.NewModelFromString(`[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act) || r.sub == "admin"`)

	casbinAdapter, err := gormadapter.NewAdapterByDBWithCustomTable(global.Db, &CasbinRule{})
	if err != nil {
		global.Logger.Errorf("[casbin]: Init adapter error: {%s}", err.Error())
		return nil
	}

	casbinEnforcer, err := casbin.NewEnforcer(casbinModel, casbinAdapter)
	if err != nil {
		global.Logger.Errorf("[casbin]: Init enforcer error: {%s}", err.Error())
		return nil
	}

	err = casbinEnforcer.LoadPolicy()
	if err != nil {
		global.Logger.Errorf("[casbin]: casbin loadPolicy error: {%s}", err.Error())
		return nil
	}
	return casbinEnforcer
}
