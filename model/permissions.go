package model

type Permissions struct {
	ID    uint64 `json:"id" gorm:"primary_key;AUTO_INCREMENT;column:id"`
	Name  string `json:"name" gorm:"name"`
	Roles []Role `json:"permissions" gorm:"many2many:permissions_roles;"`
}

const PermissionsTableName = "tb_permissions"

// TableName :
func (m *Permissions) TableName() string {
	return PermissionsTableName
}
