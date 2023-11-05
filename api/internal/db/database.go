package db

import "github.com/NimbusX-CMS/NimbusX/api/internal/models"

type DataBase interface {
	EnsureTablesCreation() error

	GetUser(userId int) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	GetUsers() ([]models.User, error)
	CreateUser(user models.User) (models.User, error)
	UpdateUser(user models.User) (models.User, error)
	DeleteUser(userId int) error

	GetSpace(spaceId int) (models.Space, error)
	GetSpaces() ([]models.Space, error)
	CreateSpace(space models.Space) (models.Space, error)
	UpdateSpace(space models.Space) (models.Space, error)
	DeleteSpace(spaceId int) error

	GetSpaceAccess(userId int, spaceId int) (models.SpaceAccess, error)
	GetSpaceAccesses(userId int) ([]models.SpaceAccess, error)
	CreateSpaceAccess(spaceAccess models.SpaceAccess) (models.SpaceAccess, error)
	UpdateSpaceAccess(spaceAccess models.SpaceAccess) (models.SpaceAccess, error)
	DeleteSpaceAccess(userId int, spaceId int) error
}
