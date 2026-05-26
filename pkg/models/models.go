package models

import (
	"time"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"uniqueIndex;not null"`
	Password  string    `json:"-" gorm:"not null"`
	IsAdmin   bool      `json:"is_admin" gorm:"default:false"`
	UnitID    *uint     `json:"unit_id"`
	Unit      *Unit     `json:"unit,omitempty" gorm:"foreignKey:UnitID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Unit struct {
	ID        uint        `json:"id" gorm:"primaryKey"`
	Name      string      `json:"name" gorm:"uniqueIndex;not null"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	Users     []User      `json:"-" gorm:"foreignKey:UnitID"`
	Tests     []SpeedTest `json:"-" gorm:"foreignKey:UnitID"`
}

type SpeedTest struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UnitID     uint      `json:"unit_id" gorm:"not null"`
	Unit       Unit      `json:"unit,omitempty" gorm:"foreignKey:UnitID"`
	Download   float64   `json:"download" gorm:"not null"`
	Upload     float64   `json:"upload"`
	Ping       float64   `json:"ping"`
	Jitter     float64   `json:"jitter"`
	IP         string    `json:"ip"`
	UserAgent  string    `json:"user_agent"`
	CreatedAt  time.Time `json:"created_at"`
}

type TopologyLink struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	FromUnitID uint     `json:"from_unit_id" gorm:"not null"`
	FromUnit  Unit      `json:"from_unit,omitempty" gorm:"foreignKey:FromUnitID"`
	ToUnitID   uint      `json:"to_unit_id" gorm:"not null"`
	ToUnit    Unit      `json:"to_unit,omitempty" gorm:"foreignKey:ToUnitID"`
	Bandwidth float64   `json:"bandwidth"`
	CreatedAt time.Time `json:"created_at"`
}
