package controllers

import (
	"github.com/fleezesd/xscheduler/pkg/xscheduler/controllers/migration"
	"github.com/fleezesd/xscheduler/pkg/xscheduler/framework/runtime"
)

func NewControllerRegistry() runtime.Registry {
	return runtime.Registry{
		migration.Name: migration.New,
	}
}
