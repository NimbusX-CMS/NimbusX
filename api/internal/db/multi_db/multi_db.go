package multi_db

import (
	"errors"
	"github.com/NimbusX-CMS/NimbusX/api/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type MultiDB struct {
	db *gorm.DB
}

func NewMultiDB(db *gorm.DB) *MultiDB {
	return &MultiDB{
		db: db,
	}
}

func (m *MultiDB) EnsureTablesCreation() error {
	return m.db.AutoMigrate(&models.User{}, &models.Language{}, &models.Space{}, &models.SpaceAccess{})
}

func (m *MultiDB) GetUser(userId int) (models.User, error) {
	var user models.User
	err := m.db.First(&user, userId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, nil
		}
		return models.User{}, err
	}
	return user, err
}

func (m *MultiDB) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := m.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, nil
		}
		return models.User{}, err
	}
	return user, nil
}

func (m *MultiDB) GetUsers() ([]models.User, error) {
	var users []models.User
	err := m.db.Find(&users).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.User{}, nil
		}
		return []models.User{}, err
	}

	return users, err
}

func (m *MultiDB) CreateUser(user models.User) (models.User, error) {
	user.ID = 0
	err := m.db.Create(&user).Error
	return user, err
}

func (m *MultiDB) UpdateUser(user models.User) (models.User, error) {
	err := m.db.Save(&user).Error
	return user, err
}

func (m *MultiDB) DeleteUser(userId int) error {
	return m.db.Delete(&models.User{}, userId).Error
}

func ConnectToSQLite(databasePath string) (*MultiDB, error) {
	db, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return NewMultiDB(db), nil
}

func (m *MultiDB) GetSpace(spaceId int) (models.Space, error) {
	var space models.Space
	err := m.db.Preload("Languages").Preload("PrimaryLanguage").First(&space, spaceId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Space{}, nil
		}
		return models.Space{}, err
	}
	return space, nil
}

func (m *MultiDB) GetSpaces() ([]models.Space, error) {
	var spaces []models.Space
	err := m.db.Preload("Languages").Preload("PrimaryLanguage").Find(&spaces).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.Space{}, nil
		}
		return []models.Space{}, err
	}
	return spaces, nil
}

func (m *MultiDB) CreateSpace(space models.Space) (models.Space, error) {
	err := m.db.Create(&space).Error
	return space, err
}

func (m *MultiDB) UpdateSpace(space models.Space) (models.Space, error) {
	err := m.db.Save(&space).Error
	return space, err
}

func (m *MultiDB) DeleteSpace(spaceId int) error {
	return m.db.Delete(&models.Space{}, spaceId).Error
}

func (m *MultiDB) GetSpaceAccess(userId int, spaceId int) (models.SpaceAccess, error) {
	var spaceAccess models.SpaceAccess
	err := m.db.Where("user_id = ? AND space_id = ?", userId, spaceId).First(&spaceAccess).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.SpaceAccess{}, nil
		}
		return models.SpaceAccess{}, err
	}
	return spaceAccess, nil
}

func (m *MultiDB) GetSpaceAccesses(userId int) ([]models.SpaceAccess, error) {
	var spaceAccesses []models.SpaceAccess
	err := m.db.Where("user_id = ?", userId).Find(&spaceAccesses).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.SpaceAccess{}, nil
		}
		return []models.SpaceAccess{}, err
	}
	return spaceAccesses, nil
}

func (m *MultiDB) CreateSpaceAccess(spaceAccess models.SpaceAccess) (models.SpaceAccess, error) {
	err := m.db.Create(&spaceAccess).Error
	return spaceAccess, err
}

func (m *MultiDB) UpdateSpaceAccess(spaceAccess models.SpaceAccess) (models.SpaceAccess, error) {
	err := m.db.Save(&spaceAccess).Error
	return spaceAccess, err
}

func (m *MultiDB) DeleteSpaceAccess(userId int, spaceId int) error {
	return m.db.Delete(&models.SpaceAccess{}, "user_id = ? AND space_id = ?", userId, spaceId).Error
}
