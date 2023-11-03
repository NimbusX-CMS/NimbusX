package multi_db

import (
	"reflect"
	"strings"
	"testing"

	"github.com/NimbusX-CMS/NimbusX/api/internal/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() (*MultiDB, func()) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	multiDB := NewMultiDB(db)

	err := multiDB.EnsureTablesCreation()
	if err != nil {
		panic("Failed to migrate database")
	}

	return multiDB, func() {
		_ = db.Migrator().DropTable(&models.User{})
	}
}

func TestMultiDB_GetUser(t *testing.T) {
	multiDB, teardown := setupTestDB()
	defer teardown()

	user := models.User{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	createdUser, _ := multiDB.CreateUser(user)

	tests := []struct {
		name          string
		id            int
		expectedUser  models.User
		expectedError error
	}{
		{
			name:          "Existing User",
			id:            createdUser.ID,
			expectedUser:  models.User{ID: createdUser.ID, Name: "John Doe", Email: "john@example.com"},
			expectedError: nil,
		},
		{
			name:          "Non-Existing User",
			id:            1234,
			expectedUser:  models.User{},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resultUser, err := multiDB.GetUser(test.id)
			assert.Equal(t, test.expectedUser, resultUser)
			assert.Equal(t, test.expectedError, err)
		})
	}
}

func TestMultiDB_GetUserByEmail(t *testing.T) {
	multiDB, teardown := setupTestDB()
	defer teardown()

	user := models.User{
		Name:  "John",
		Email: "existing@example.com",
	}

	_, _ = multiDB.CreateUser(user)

	tests := []struct {
		name         string
		email        string
		expectedUser models.User
	}{
		{
			name:         "Existing User",
			email:        "existing@example.com",
			expectedUser: models.User{ID: 1, Name: "John", Email: "existing@example.com"},
		},
		{
			name:         "Non-Existing User",
			email:        "not-existing@example.com",
			expectedUser: models.User{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resultUser, err := multiDB.GetUserByEmail(test.email)
			assert.Equal(t, test.expectedUser, resultUser)
			assert.Equal(t, nil, err)
		})
	}
}

func TestMultiDB_GetUsers(t *testing.T) {
	multiDB, teardown := setupTestDB()
	defer teardown()

	users := []models.User{
		{
			Name:  "Alice",
			Email: "alice@example.com",
		},
		{
			Name:  "Bob",
			Email: "bob@example.com",
		},
		{
			Name:  "Charlie",
			Email: "charlie@example.com",
		},
	}

	for _, user := range users {
		_, _ = multiDB.CreateUser(user)
	}

	tests := []struct {
		name           string
		expectedLength int
	}{
		{
			name:           "Get All Users",
			expectedLength: len(users),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resultUsers, err := multiDB.GetUsers()
			assert.NoError(t, err)
			assert.Len(t, resultUsers, test.expectedLength)

			for i, user := range users {
				assert.Equal(t, user.Name, resultUsers[i].Name)
				assert.Equal(t, user.Email, resultUsers[i].Email)
			}
		})
	}
}

func TestMultiDB_UpdateUser(t *testing.T) {
	multiDB, teardown := setupTestDB()
	defer teardown()

	user := models.User{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	createdUser, _ := multiDB.CreateUser(user)

	tests := []struct {
		name         string
		userToUpdate models.User
		expectedUser models.User
		expectedErr  error
	}{
		{
			name:         "Update User",
			userToUpdate: models.User{ID: createdUser.ID, Name: "John Smith", Email: "john@example.com"},
			expectedUser: models.User{ID: createdUser.ID, Name: "John Smith", Email: "john@example.com"},
			expectedErr:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			updatedUser, err := multiDB.UpdateUser(test.userToUpdate)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedUser, updatedUser)
		})
	}
}

func TestMultiDB_DeleteUser(t *testing.T) {
	multiDB, teardown := setupTestDB()
	defer teardown()

	user := models.User{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	createdUser, _ := multiDB.CreateUser(user)

	tests := []struct {
		name         string
		idToDelete   int
		expectedErr  error
		expectedUser models.User
	}{
		{
			name:         "Delete User",
			idToDelete:   createdUser.ID,
			expectedErr:  nil,
			expectedUser: models.User{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := multiDB.DeleteUser(test.idToDelete)
			assert.NoError(t, err)

			deletedUser, err := multiDB.GetUser(test.idToDelete)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedUser, deletedUser)
		})
	}
}

func TestMultiDB_GetSpace(t *testing.T) {
	multiDB, teardown := setupTestDB()
	defer teardown()

	space := models.Space{
		Name:              "Test Space",
		Color1:            "#FFFFFF",
		Color2:            "#000000",
		Color3:            "#FF0000",
		Color4:            "#00FF00",
		ImageUrl:          "https://example.com/image.jpg",
		PrimaryLanguageID: 1,
		Languages: []models.Language{
			{
				Name: "de-DE",
			},
			{
				Name: "fr-FR",
			},
		},
	}

	createdSpace, _ := multiDB.CreateSpace(space)

	tests := []struct {
		name          string
		spaceID       int
		expectedSpace models.Space
		expectedErr   error
	}{
		{
			name:    "Existing Space",
			spaceID: createdSpace.ID,
			expectedSpace: models.Space{
				ID:                createdSpace.ID,
				Name:              "Test Space",
				Color1:            "#FFFFFF",
				Color2:            "#000000",
				Color3:            "#FF0000",
				Color4:            "#00FF00",
				ImageUrl:          "https://example.com/image.jpg",
				PrimaryLanguageID: createdSpace.PrimaryLanguageID,
				PrimaryLanguage: models.Language{
					ID:   1,
					Name: "de-DE",
				},
				Languages: []models.Language{
					{
						ID:   1,
						Name: "de-DE",
					},
					{
						ID:   2,
						Name: "fr-FR",
					},
				},
			},
			expectedErr: nil,
		},
		{
			name:          "Non-Existing Space",
			spaceID:       1234,
			expectedSpace: models.Space{},
			expectedErr:   nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resultSpace, err := multiDB.GetSpace(test.spaceID)
			assert.Equal(t, test.expectedSpace, resultSpace)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}

func TestMultiDB_GetSpaces(t *testing.T) {
	multiDB, teardown := setupTestDB()
	defer teardown()

	spaces := []models.Space{
		{
			Name:     "Space 1",
			Color1:   "#FFFFFF",
			Color2:   "#000000",
			Color3:   "#FF0000",
			Color4:   "#00FF00",
			ImageUrl: "https://example.com/image1.jpg",
			PrimaryLanguage: models.Language{
				Name: "en-EN",
			},
			Languages: []models.Language{
				{
					Name: "de-DE",
				},
			},
		},
		{
			Name:     "Space 2",
			Color1:   "#000000",
			Color2:   "#FFFFFF",
			Color3:   "#00FF00",
			Color4:   "#FF0000",
			ImageUrl: "https://example.com/image2.jpg",
			PrimaryLanguage: models.Language{
				Name: "es-ES",
			},
			Languages: []models.Language{
				{
					Name: "es-ES",
				},
				{
					Name: "it-IT",
				},
			},
		},
	}

	for _, space := range spaces {
		_, _ = multiDB.CreateSpace(space)
	}

	tests := []struct {
		name           string
		expectedLength int
	}{
		{
			name:           "Get All Spaces",
			expectedLength: len(spaces),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resultSpaces, err := multiDB.GetSpaces()
			assert.NoError(t, err)
			assert.Len(t, resultSpaces, test.expectedLength)

			for i, space := range spaces {
				assert.Equal(t, space.Name, resultSpaces[i].Name)
				assert.Equal(t, space.Color1, resultSpaces[i].Color1)
				assert.Equal(t, space.Color2, resultSpaces[i].Color2)
				assert.Equal(t, space.Color3, resultSpaces[i].Color3)
				assert.Equal(t, space.Color4, resultSpaces[i].Color4)
				assert.Equal(t, space.ImageUrl, resultSpaces[i].ImageUrl)
				assert.Equal(t, space.PrimaryLanguage.Name, resultSpaces[i].PrimaryLanguage.Name)

				for j, lang := range space.Languages {
					assert.Equal(t, lang.Name, resultSpaces[i].Languages[j].Name)
				}
			}
		})
	}
}

func TestMultiDB_CreateSpace(t *testing.T) {
	type testCase struct {
		name           string
		space          models.Space
		expectedFields []string
	}

	testCases := []testCase{
		{
			name: "Create space successfully",
			space: models.Space{
				Name:     "New Space",
				Color1:   "#FFFFFF",
				Color2:   "#000000",
				Color3:   "#FF0000",
				Color4:   "#00FF00",
				ImageUrl: "https://example.com/new-image.jpg",
				PrimaryLanguage: models.Language{
					Name: "en-EN",
				},
				Languages: []models.Language{
					{
						Name: "de-DE",
					},
					{
						Name: "fr-FR",
					},
				},
			},
			expectedFields: []string{"Name", "Color1", "Color2", "Color3", "Color4", "ImageUrl", "PrimaryLanguage.Name", "Languages"},
		},
	}

	multiDB, teardown := setupTestDB()
	defer teardown()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			createdSpace, err := multiDB.CreateSpace(tc.space)
			assert.NoError(t, err)
			assert.NotZero(t, createdSpace.ID)

			for _, field := range tc.expectedFields {
				assert.Equal(t, getField(tc.space, field), getField(createdSpace, field))
			}
		})
	}
}

func TestMultiDB_UpdateSpace(t *testing.T) {
	type testCase struct {
		name           string
		originalSpace  models.Space
		updatedSpace   models.Space
		expectedError  error
		expectedFields []string
	}

	testCases := []testCase{
		{
			name: "Update space successfully",
			originalSpace: models.Space{
				Name:     "Old Space",
				Color1:   "#FFFFFF",
				Color2:   "#000000",
				Color3:   "#FF0000",
				Color4:   "#00FF00",
				ImageUrl: "https://example.com/old-image.jpg",
				PrimaryLanguage: models.Language{
					Name: "en-EN",
				},
				Languages: []models.Language{
					{
						Name: "de-DE",
					},
				},
			},
			updatedSpace: models.Space{
				Name:     "Updated Space",
				Color1:   "#000000",
				Color2:   "#FFFFFF",
				Color3:   "#00FF00",
				Color4:   "#FF0000",
				ImageUrl: "https://example.com/updated-image.jpg",
				PrimaryLanguage: models.Language{
					Name: "fr-FR",
				},
				Languages: []models.Language{
					{
						Name: "es-ES",
					},
					{
						Name: "it-IT",
					},
				},
			},
			expectedError:  nil,
			expectedFields: []string{"Name", "Color1", "Color2", "Color3", "Color4", "ImageUrl", "PrimaryLanguage.Name", "Languages"},
		},
	}

	multiDB, teardown := setupTestDB()
	defer teardown()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			createdSpace, _ := multiDB.CreateSpace(tc.originalSpace)
			tc.updatedSpace.ID = createdSpace.ID
			updatedSpace, err := multiDB.UpdateSpace(tc.updatedSpace)

			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, createdSpace.ID, updatedSpace.ID)

				for _, field := range tc.expectedFields {
					assert.Equal(t, getField(tc.updatedSpace, field), getField(updatedSpace, field))
				}
			}
		})
	}
}

func TestMultiDB_DeleteSpace(t *testing.T) {
	type testCase struct {
		name           string
		space          models.Space
		expectedError  error
		expectedFields []string
	}

	testCases := []testCase{
		{
			name: "Delete space successfully",
			space: models.Space{
				Name:     "Space to Delete",
				Color1:   "#FFFFFF",
				Color2:   "#000000",
				Color3:   "#FF0000",
				Color4:   "#00FF00",
				ImageUrl: "https://example.com/delete-image.jpg",
				PrimaryLanguage: models.Language{
					Name: "en-EN",
				},
				Languages: []models.Language{
					{
						Name: "de-DE",
					},
				},
			},
			expectedError:  nil,
			expectedFields: []string{},
		},
	}

	multiDB, teardown := setupTestDB()
	defer teardown()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			createdSpace, _ := multiDB.CreateSpace(tc.space)

			err := multiDB.DeleteSpace(createdSpace.ID)
			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
			} else {
				assert.NoError(t, err)

				deletedSpace, err := multiDB.GetSpace(createdSpace.ID)
				assert.NoError(t, err)
				assert.Equal(t, models.Space{}, deletedSpace)

				for _, field := range tc.expectedFields {
					assert.Equal(t, getField(tc.space, field), getField(deletedSpace, field))
				}
			}
		})
	}
}

func getField(space models.Space, field string) interface{} {
	fields := strings.Split(field, ".")
	value := reflect.ValueOf(space)

	for _, f := range fields {
		if value.Kind() == reflect.Ptr {
			value = value.Elem()
		}
		value = value.FieldByName(f)
	}
	return value.Interface()
}

func TestMultiDB_GetSpaceAccess(t *testing.T) {
	multiDB, teardown := setupTestDB()
	defer teardown()

	// Prepare test data
	userId := 1
	spaceId := 1
	existingAccess := models.SpaceAccess{
		UserID:  userId,
		SpaceID: spaceId,
		Admin:   true,
	}
	_, _ = multiDB.CreateSpaceAccess(existingAccess)

	tests := []struct {
		name           string
		userId         int
		spaceId        int
		expectedAccess models.SpaceAccess
		expectedError  error
	}{
		{
			name:           "Existing Access",
			userId:         userId,
			spaceId:        spaceId,
			expectedAccess: existingAccess,
			expectedError:  nil,
		},
		{
			name:           "Non-Existing Access",
			userId:         2,
			spaceId:        2,
			expectedAccess: models.SpaceAccess{},
			expectedError:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			access, err := multiDB.GetSpaceAccess(test.userId, test.spaceId)
			assert.Equal(t, test.expectedAccess, access)
			assert.Equal(t, test.expectedError, err)
		})
	}
}

func TestMultiDB_GetSpaceAccesses(t *testing.T) {
	multiDB, teardown := setupTestDB()
	defer teardown()

	// Prepare test data
	userId := 1
	accesses := []models.SpaceAccess{
		{UserID: userId, SpaceID: 1, Admin: true},
		{UserID: userId, SpaceID: 2, Admin: false},
	}
	for _, access := range accesses {
		_, _ = multiDB.CreateSpaceAccess(access)
	}

	tests := []struct {
		name             string
		userId           int
		expectedAccesses []models.SpaceAccess
		expectedError    error
	}{
		{
			name:             "Existing Accesses",
			userId:           userId,
			expectedAccesses: accesses,
			expectedError:    nil,
		},
		{
			name:             "Non-Existing Accesses",
			userId:           2,
			expectedAccesses: []models.SpaceAccess{},
			expectedError:    nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resultAccesses, err := multiDB.GetSpaceAccesses(test.userId)
			assert.Equal(t, test.expectedAccesses, resultAccesses)
			assert.Equal(t, test.expectedError, err)
		})
	}
}
func TestMultiDB_CreateSpaceAccess(t *testing.T) {
	multiDB, teardown := setupTestDB()
	defer teardown()

	// Prepare test data
	userId := 1
	spaceId := 1

	tests := []struct {
		name          string
		access        models.SpaceAccess
		expectedError error
	}{
		{
			name: "Create Access",
			access: models.SpaceAccess{
				UserID:  userId,
				SpaceID: spaceId,
				Admin:   true,
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := multiDB.CreateSpaceAccess(test.access)
			assert.Equal(t, test.expectedError, err)

			// Check if the access was created successfully
			resultAccess, _ := multiDB.GetSpaceAccess(test.access.UserID, test.access.SpaceID)
			assert.Equal(t, test.access, resultAccess)
		})
	}
}

func TestMultiDB_UpdateSpaceAccess(t *testing.T) {
	multiDB, teardown := setupTestDB()
	defer teardown()

	// Prepare test data
	userId := 1
	spaceId := 1
	access := models.SpaceAccess{
		UserID:  userId,
		SpaceID: spaceId,
		Admin:   true,
	}
	_, _ = multiDB.CreateSpaceAccess(access)

	tests := []struct {
		name          string
		access        models.SpaceAccess
		expectedError error
	}{
		{
			name:          "Update Access",
			access:        models.SpaceAccess{UserID: userId, SpaceID: spaceId, Admin: false},
			expectedError: nil,
		},
		{
			name:          "Update Non-Existing Access",
			access:        models.SpaceAccess{UserID: 2, SpaceID: 2, Admin: true},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := multiDB.UpdateSpaceAccess(test.access)
			assert.Equal(t, test.expectedError, err)

			// Check if the access was updated successfully
			resultAccess, _ := multiDB.GetSpaceAccess(test.access.UserID, test.access.SpaceID)
			assert.Equal(t, test.access, resultAccess)
		})
	}
}

func TestMultiDB_DeleteSpaceAccess(t *testing.T) {
	multiDB, teardown := setupTestDB()
	defer teardown()

	// Prepare test data
	userId := 1
	spaceId := 1
	access := models.SpaceAccess{
		UserID:  userId,
		SpaceID: spaceId,
		Admin:   true,
	}
	_, _ = multiDB.CreateSpaceAccess(access)

	tests := []struct {
		name          string
		userId        int
		spaceId       int
		expectedError error
	}{
		{
			name:          "Delete Access",
			userId:        userId,
			spaceId:       spaceId,
			expectedError: nil,
		},
		{
			name:          "Delete Non-Existing Access",
			userId:        2,
			spaceId:       2,
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := multiDB.DeleteSpaceAccess(test.userId, test.spaceId)
			assert.Equal(t, test.expectedError, err)

			// Check if the access was deleted successfully
			_, err = multiDB.GetSpaceAccess(test.userId, test.spaceId)
			assert.Equal(t, test.expectedError, err)
		})
	}
}
