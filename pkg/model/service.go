package model

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

type Service struct {
	gorm.Model
	Name                 string
	Token                string
	Failure              bool
	Recovering           bool
	Maintenance          bool
	PingAt               time.Time
	RecoveryAt           time.Time
	MaintenanceFailureAt time.Time
	FailureTimeout       time.Duration
	RecoveryInterval     time.Duration
	MaintenanceTimeout   time.Duration
	Users                []User `gorm:"many2many:subscriptions;"`
}

func (s *Service) Ping() {
	s.PingAt = time.Now()

	if !s.Recovering && !s.Maintenance {
		s.RecoveryAt = s.PingAt.Add(s.RecoveryInterval)
		s.Recovering = true
	}
}

func (s *Service) EnableMaintenance() {
	s.MaintenanceFailureAt = time.Now().Add(s.MaintenanceTimeout)
	s.Maintenance = true
	s.Recovering = false
}

func (s *Service) DisableMaintenance() {
	s.PingAt = time.Now()
	s.Maintenance = false
}

func (s *Service) Display() {
	fmt.Println("name                   = ", s.Name)
	fmt.Println("token                  = ", s.Token)
	fmt.Println("failure                = ", s.Failure)
	fmt.Println("recovering             = ", s.Recovering)
	fmt.Println("maintenance            = ", s.Maintenance)
	fmt.Println("failure timeout        = ", s.FailureTimeout)
	fmt.Println("recovery interval      = ", s.RecoveryInterval)
	fmt.Println("maintenance timeout    = ", s.MaintenanceTimeout)
	fmt.Println("ping at                = ", s.PingAt)
	fmt.Println("recovery at            = ", s.RecoveryAt)
	fmt.Println("maintenance failure at = ", s.MaintenanceFailureAt)
	fmt.Println()
}
