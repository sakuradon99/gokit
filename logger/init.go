package logger

import (
	"github.com/sakuradon99/ioc"
	"github.com/sirupsen/logrus"
)

func init() {
	ioc.Register[Config](ioc.Optional())
	ioc.Register[logrus.Logger](ioc.Constructor(newLogger))
}
