package dao

import "time"

type VerificationCode struct {
	ID         uint      `gorm:"primaryKey"`     // ID único (clave primaria)
	UserID     uint      `gorm:"not null"`       // ID del usuario
	Code       string    `gorm:"not null"`       // Código de verificación
	Expiration time.Time `gorm:"not null"`       // Fecha y hora de expiración
	Claimed    bool      `gorm:"default:false"`  // Si el código fue reclamado
	CreatedAt  time.Time `gorm:"autoCreateTime"` // Fecha y hora de creación automática
	UpdatedAt  time.Time `gorm:"autoUpdateTime"` // Fecha de actualización automática
}
