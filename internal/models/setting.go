package models

type Setting struct {
	ID    uint   `gorm:"primaryKey"`
	Key   string `gorm:"uniqueIndex;not null"`
	Value string
}
