package dao

type ProductGroupMember struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	GroupID   uint `gorm:"not null;index"`
	ProductID uint `gorm:"not null;index"`

	// Relaciones
	Group   ProductGroup `gorm:"foreignKey:GroupID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
	Product Product      `gorm:"foreignKey:ProductID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
}

func (ProductGroupMember) TableName() string {
	return "product_group_members"
}
