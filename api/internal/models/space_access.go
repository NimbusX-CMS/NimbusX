package models

type SpaceAccess struct {
	UserID  int  `json:"userId" gorm:"primaryKey"`
	SpaceID int  `json:"spaceId" gorm:"primaryKey"`
	Admin   bool `json:"admin"`
}
