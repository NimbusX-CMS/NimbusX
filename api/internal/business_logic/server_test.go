package business_logic

import (
	"fmt"
	"github.com/NimbusX-CMS/NimbusX/api/internal/error_msg"
	"github.com/NimbusX-CMS/NimbusX/api/internal/models"
	"net/http"
	"testing"
)

func TestServer_PostSpace(t *testing.T) {
	server := setupServer(t)

	tests := TestCases{
		{
			name:               "PostSpace",
			Url:                "/spaces",
			RequestBody:        `{"name":"Test Space"}`,
			ResponseModel:      &models.Space{},
			ExpectedBody:       &models.Space{ID: 1, Name: "Test Space"},
			ExpectedStatusCode: http.StatusCreated,
		},
	}
	tests.testStaticUrlCases(t, server.PostSpace)
}

func TestServer_DeleteSpaceSpaceId(t *testing.T) {
	server := setupServer(t)
	space, err := server.DB.CreateSpace(models.Space{Name: "Test Space"})
	if err != nil {
		t.Error("Cannot create sample space", err)
	}

	tests := TestCases{
		{
			name:               "DeleteSpaceSpaceId",
			Url:                fmt.Sprintf("/spaces/%v", space.ID),
			ID:                 space.ID,
			ExpectedStatusCode: http.StatusOK,
		},
		{
			name:               "DeleteSpaceSpaceId, but space does not exist",
			Url:                "/spaces/123123",
			ID:                 123123,
			ResponseModel:      &error_msg.Error{},
			ExpectedBody:       &error_msg.Error{Error: error_msg.ErrorSpaceWithIdNotFound},
			ExpectedStatusCode: http.StatusNotFound,
		},
	}
	tests.testDynamicIntUrlCases(t, server.DeleteSpaceSpaceId)
}

func TestServer_GetSpaceSpaceId(t *testing.T) {
	server := setupServer(t)
	space, err := server.DB.CreateSpace(models.Space{ID: 1, Name: "Test Space"})
	if err != nil {
		t.Error("Cannot create sample space", err)
	}
	tests := TestCases{
		{
			name:               "GetSpaceSpaceId",
			Url:                fmt.Sprintf("/spaces/%v", space.ID),
			ID:                 space.ID,
			ResponseModel:      &models.Space{},
			ExpectedBody:       &models.Space{ID: space.ID, Name: "Test Space", Languages: []models.Language{}},
			ExpectedStatusCode: http.StatusOK,
		},
		{
			name:               "GetSpaceSpaceId, but space does not exist",
			Url:                "/spaces/123123",
			ID:                 123123,
			ResponseModel:      &error_msg.Error{},
			ExpectedBody:       &error_msg.Error{Error: error_msg.ErrorSpaceWithIdNotFound},
			ExpectedStatusCode: http.StatusNotFound,
		},
	}

	tests.testDynamicIntUrlCases(t, server.GetSpaceSpaceId)
}

func TestServer_PutSpaceSpaceId(t *testing.T) {
	server := setupServer(t)
	space, err := server.DB.CreateSpace(models.Space{ID: 1, Name: "Test Space"})
	if err != nil {
		t.Error("Cannot create sample space", err)
	}

	tests := TestCases{
		{
			name:               "PutSpaceSpaceId",
			Url:                fmt.Sprintf("/spaces/%v", space.ID),
			ID:                 space.ID,
			RequestBody:        `{"name":"Updated Space"}`,
			ResponseModel:      &models.Space{},
			ExpectedBody:       &models.Space{ID: space.ID, Name: "Updated Space"},
			ExpectedStatusCode: http.StatusOK,
		},
		{
			name:               "PutSpaceSpaceId, but space does not exist",
			Url:                "/spaces/123123",
			ID:                 123123,
			RequestBody:        `{"name":"Updated Space"}`,
			ResponseModel:      &error_msg.Error{},
			ExpectedBody:       &error_msg.Error{Error: error_msg.ErrorSpaceWithIdNotFound},
			ExpectedStatusCode: http.StatusNotFound,
		},
	}
	tests.testDynamicIntUrlCases(t, server.PutSpaceSpaceId)
}

func TestServer_GetSpaces(t *testing.T) {
	server := setupServer(t)

	tests := TestCases{
		{
			name:               "GetSpaces",
			Url:                "/spaces",
			ResponseModel:      &[]models.Space{},
			ExpectedBody:       &[]models.Space{},
			ExpectedStatusCode: http.StatusOK,
		},
	}
	tests.testStaticUrlCases(t, server.GetSpaces)
}

func TestServer_GetUserUserIdSpaces(t *testing.T) {
	server := setupServer(t)
	_, err := server.DB.CreateUser(models.User{ID: 1, Name: "Test User"})
	if err != nil {
		t.Error("Cannot create sample user", err)
	}
	space, err := server.DB.CreateSpace(models.Space{ID: 1, Name: "Test Space"})
	if err != nil {
		t.Error("Cannot create sample space", err)
	}
	spaceAccess, err := server.DB.CreateSpaceAccess(models.SpaceAccess{UserID: 1, SpaceID: space.ID})
	if err != nil {
		t.Error("Cannot create sample space access", err)
	}

	tests := TestCases{
		{
			name:               "GetUserUserIdSpaces",
			Url:                "/users/1/spaces",
			ID:                 1,
			ResponseModel:      &[]models.SpaceAccess{},
			ExpectedBody:       &[]models.SpaceAccess{spaceAccess},
			ExpectedStatusCode: http.StatusOK,
		},
		{
			name:               "GetUserUserIdSpaces, but user does not exist",
			Url:                "/users/123123/spaces",
			ID:                 123123,
			ResponseModel:      &error_msg.Error{},
			ExpectedBody:       &error_msg.Error{Error: error_msg.ErrorUserWithIdNotFound},
			ExpectedStatusCode: http.StatusNotFound,
		},
	}
	tests.testDynamicIntUrlCases(t, server.GetUserUserIdSpaces)
}

func TestServer_DeleteUserUserIdSpaceSpaceId(t *testing.T) {
	server := setupServer(t)
	_, err := server.DB.CreateUser(models.User{ID: 1, Name: "Test User"})
	if err != nil {
		t.Error("Cannot create sample user", err)
	}
	space, err := server.DB.CreateSpace(models.Space{ID: 1, Name: "Test Space"})
	if err != nil {
		t.Error("Cannot create sample space", err)
	}
	spaceAccess, err := server.DB.CreateSpaceAccess(models.SpaceAccess{UserID: 1, SpaceID: space.ID})
	if err != nil {
		t.Error("Cannot create sample space access", err)
	}

	tests := TestCases{
		{
			name:               "Delete SpaceAccess",
			Url:                fmt.Sprintf("/user/%v/space/%v", spaceAccess.UserID, spaceAccess.SpaceID),
			ID:                 spaceAccess.UserID,
			ID2:                spaceAccess.SpaceID,
			ExpectedStatusCode: http.StatusOK,
		},
		{
			name:               "Delete SpaceAccess, but not existing",
			Url:                fmt.Sprintf("/user/%v/space/%v", 123123, 123123),
			ID:                 123123,
			ID2:                123123,
			ResponseModel:      &error_msg.Error{},
			ExpectedBody:       &error_msg.Error{Error: error_msg.ErrorSpaceAccessWithIdsNotFound},
			ExpectedStatusCode: http.StatusNotFound,
		},
	}

	tests.testDynamic2IntUrlCases(t, server.DeleteUserUserIdSpaceSpaceId)
}

func TestServer_PatchUserUserIdSpaces(t *testing.T) {
	server := setupServer(t)
	user, err := server.DB.CreateUser(models.User{ID: 1, Name: "Test User"})
	if err != nil {
		t.Error("Cannot create sample user", err)
	}
	space, err := server.DB.CreateSpace(models.Space{ID: 1, Name: "Test Space"})
	if err != nil {
		t.Error("Cannot create sample space", err)
	}

	tests := TestCases{
		{
			name:               "Patch SpaceAccess",
			Url:                fmt.Sprintf("/user/%v/space", space.ID),
			ID:                 user.ID,
			RequestBody:        fmt.Sprintf("{ \"userId\": %v, \"spaceId\": %v, \"admin\": false }", user.ID, space.ID),
			ResponseModel:      &models.SpaceAccess{},
			ExpectedBody:       &models.SpaceAccess{UserID: user.ID, SpaceID: space.ID, Admin: false},
			ExpectedStatusCode: http.StatusOK,
		},
		{
			name:               "Patch SpaceAccess, already exist",
			Url:                fmt.Sprintf("/user/%v/space", space.ID),
			ID:                 user.ID,
			RequestBody:        fmt.Sprintf("{ \"userId\": %v, \"spaceId\": %v, \"admin\": true }", user.ID, space.ID),
			ResponseModel:      &models.SpaceAccess{},
			ExpectedBody:       &models.SpaceAccess{UserID: user.ID, SpaceID: space.ID, Admin: true},
			ExpectedStatusCode: http.StatusOK,
		},
		{
			name:               "Patch SpaceAccess, not exist",
			Url:                fmt.Sprintf("/user/%v/space", 123123),
			ID:                 123123,
			RequestBody:        fmt.Sprintf("{ \"userId\": %v, \"spaceId\": %v, \"admin\": true }", 123123, 123123),
			ResponseModel:      &error_msg.Error{},
			ExpectedBody:       &error_msg.Error{Error: error_msg.ErrorUserWithIdNotFound},
			ExpectedStatusCode: http.StatusNotFound,
		},
	}

	tests.testDynamicIntUrlCases(t, server.PatchUserUserIdSpaces)
}

func TestPostUser(t *testing.T) {
	server := setupServer(t)

	tests := TestCases{
		{
			name:               "Create user",
			Url:                "/user",
			RequestBody:        `{"name": "John Doe", "email": "j@example.com"}`,
			ResponseModel:      &models.User{},
			ExpectedBody:       &models.User{ID: 1, Name: "John Doe", Email: "j@example.com"},
			ExpectedStatusCode: http.StatusCreated,
		},
		{
			name:               "Create another user",
			Url:                "/user",
			RequestBody:        `{"name": "Jane Doe", "email": "ja@example.com"}`,
			ResponseModel:      &models.User{},
			ExpectedBody:       &models.User{ID: 2, Name: "Jane Doe", Email: "ja@example.com"},
			ExpectedStatusCode: http.StatusCreated,
		},
		{
			name:               "Create user, but email already in use",
			Url:                "/user",
			RequestBody:        `{"name": "Jane Doe", "email": "ja@example.com"}`,
			ResponseModel:      &error_msg.Error{},
			ExpectedBody:       &error_msg.Error{Error: error_msg.ErrorEmailAlreadyInUse},
			ExpectedStatusCode: http.StatusBadRequest,
		},
	}

	tests.testStaticUrlCases(t, server.PostUser)
}

func TestServer_DeleteUserUserId(t *testing.T) {
	server := setupServer(t)
	user, err := server.DB.CreateUser(models.User{ID: 1, Name: "John Doe", Email: "a@example.com"})
	if err != nil {
		t.Error("Cannot create sample user", err)
	}
	tests := TestCases{
		{
			name:               "Delete user",
			Url:                fmt.Sprintf("/user/%v", user.ID),
			ID:                 user.ID,
			ExpectedStatusCode: http.StatusOK,
		},
		{
			name:               "Delete user",
			Url:                fmt.Sprintf("/user/%v", 12312),
			ID:                 12312,
			ExpectedStatusCode: http.StatusNotFound,
		},
	}

	tests.testDynamicIntUrlCases(t, server.DeleteUserUserId)
}

func TestServer_GetUserUserId(t *testing.T) {
	server := setupServer(t)
	user, err := server.DB.CreateUser(models.User{ID: 1, Name: "John Doe", Email: "a@example.com"})
	if err != nil {
		t.Error("Cannot create sample user", err)
	}
	tests := TestCases{
		{
			name:               "Existing User",
			Url:                fmt.Sprintf("/user/%v", user.ID),
			ID:                 user.ID,
			ResponseModel:      &models.User{},
			ExpectedBody:       &models.User{ID: 1, Name: "John Doe", Email: "a@example.com"},
			ExpectedStatusCode: http.StatusOK,
		},
		{
			name:               "Not existing User",
			Url:                fmt.Sprintf("/user/%v", 1234124),
			ID:                 1234124,
			ResponseModel:      &error_msg.Error{},
			ExpectedBody:       &error_msg.Error{Error: error_msg.ErrorUserWithIdNotFound},
			ExpectedStatusCode: http.StatusNotFound,
		},
	}

	tests.testDynamicIntUrlCases(t, server.GetUserUserId)
}

func TestServer_GetUsers(t *testing.T) {
	server := setupServer(t)

	tests := TestCases{
		{
			name:               "Get users",
			Url:                "/users",
			ResponseModel:      &[]models.User{},
			ExpectedBody:       &[]models.User{},
			ExpectedStatusCode: http.StatusOK,
		},
	}
	tests.testStaticUrlCases(t, server.GetUsers)

	users := &[]models.User{
		{ID: 1, Name: "John Doe", Email: "a@example.com"},
		{ID: 2, Name: "Jane Doe", Email: "b@example.com"},
	}
	for _, user := range *users {
		_, err := server.DB.CreateUser(user)
		if err != nil {
			t.Error("Cannot create sample user", err)
		}
	}

	tests = TestCases{
		{
			name:               "Get users",
			Url:                "/users",
			ResponseModel:      &[]models.User{},
			ExpectedBody:       users,
			ExpectedStatusCode: http.StatusOK,
		},
	}
	tests.testStaticUrlCases(t, server.GetUsers)
}

func TestServer_PutUserUserId(t *testing.T) {
	server := setupServer(t)
	user, err := server.DB.CreateUser(models.User{ID: 1, Name: "John Doee", Email: "a@example.com"})
	if err != nil {
		t.Error("Cannot create sample user", err)
	}
	tests := TestCases{
		{
			name:               "Put user by id",
			Url:                fmt.Sprintf("/user/%v", user.ID),
			ID:                 user.ID,
			RequestBody:        `{"name": "John Doe", "email": "a@example.com"}`,
			ResponseModel:      &models.User{},
			ExpectedBody:       &models.User{ID: 1, Name: "John Doe", Email: "a@example.com"},
			ExpectedStatusCode: http.StatusOK,
		},
		{
			name:               "Put user by id",
			Url:                fmt.Sprintf("/user/%v", user.ID),
			ID:                 user.ID,
			RequestBody:        `{"id": 2, "name": "John Doe", "email": "a@example.com"}`,
			ResponseModel:      &models.User{},
			ExpectedBody:       &models.User{ID: 1, Name: "John Doe", Email: "a@example.com"},
			ExpectedStatusCode: http.StatusOK,
		},
		{
			name:               "Put user by id",
			Url:                fmt.Sprintf("/user/%v", 12341234213),
			ID:                 12341234213,
			RequestBody:        `{"id": 2, "name": "John Doe", "email": "a@example.com"}`,
			ResponseModel:      &error_msg.Error{},
			ExpectedBody:       &error_msg.Error{Error: error_msg.ErrorUserWithIdNotFound},
			ExpectedStatusCode: http.StatusNotFound,
		},
	}

	tests.testDynamicIntUrlCases(t, server.PutUserUserId)
}
