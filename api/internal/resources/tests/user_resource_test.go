package resources_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/NimbusX-CMS/NimbusX/api/internal/models"
)

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
			ResponseModel:      &models.Error{},
			ExpectedBody:       &models.Error{Error: models.ErrorEmailAlreadyInUse},
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
			ResponseModel:      &models.Error{},
			ExpectedBody:       &models.Error{Error: models.ErrorUserWithIdNotFound},
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
			ResponseModel:      &models.Error{},
			ExpectedBody:       &models.Error{Error: models.ErrorUserWithIdNotFound},
			ExpectedStatusCode: http.StatusNotFound,
		},
	}

	tests.testDynamicIntUrlCases(t, server.PutUserUserId)
}
