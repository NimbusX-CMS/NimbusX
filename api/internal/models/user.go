package models

type User struct {
	ID            int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	InstanceAdmin bool   `json:"instanceAdmin"`
	Origin        bool   `json:"origin"`
}
