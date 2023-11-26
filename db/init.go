package db

import (
	"github.com/sakuradon99/gokit/internal/db"
	"github.com/sakuradon99/ioc"
	"gorm.io/gorm"
)

func init() {
	ioc.Register[db.Config](ioc.Optional())
	ioc.Register[gorm.DB](ioc.Constructor(db.InitGorm), ioc.Optional())
	ioc.Register[db.ManagerImpl](ioc.Conditional("#db != nil"))
}
