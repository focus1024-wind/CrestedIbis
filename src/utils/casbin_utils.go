package utils

import (
	"CrestedIbis/src/global"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

// CasbinRule 适配 Casbin 权限配置，参考：https://github.com/casbin/gorm-adapter
type CasbinRule struct {
	ID      uint `gorm:"primaryKey;autoIncrement" json:"id"`
	Ptype   string
	RoleKey string `gorm:"column:v0" json:"role_key"`
	Path    string `gorm:"column:v1" json:"path"`
	Method  string `gorm:"column:v2" json:"method"`
}

// TableName gorm-adapter 根据`casbin_rule`表名获取Policy，所以必须手动指定表名，避免gorm生成表名不一致
func (CasbinRule) TableName() string {
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
		global.Logger.Errorf("[casbin]: 初始化 CasbinAdapter 失败: %s", err.Error())
		return nil
	}

	casbinEnforcer, err := casbin.NewEnforcer(casbinModel, casbinAdapter)
	if err != nil {
		global.Logger.Errorf("[casbin]: 初始化 CasbinEnforcer 失败: %s", err.Error())
		return nil
	}

	err = casbinEnforcer.LoadPolicy()
	if err != nil {
		global.Logger.Errorf("[casbin]: CasbinEnforcer 加载策略失败: %s", err.Error())
		return nil
	}
	return casbinEnforcer
}
