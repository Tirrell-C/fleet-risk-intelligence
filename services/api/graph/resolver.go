package graph

import (
	"gorm.io/gorm"
	"github.com/Tirrell-C/fleet-risk-intelligence/pkg/config"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{
	DB     *gorm.DB
	Config *config.Config
}