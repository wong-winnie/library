package common

import (
	"github.com/sirupsen/logrus"
	"runtime"
)

//***********************打印日志************************************
func SimplePanic(args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	args = append(args, file)
	args = append(args, line)
	logrus.Fatal(args)
} //***********************打印日志************************************
