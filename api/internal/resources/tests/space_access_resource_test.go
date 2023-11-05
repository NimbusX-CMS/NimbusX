package resources_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/NimbusX-CMS/NimbusX/api/internal/models"
)

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
			ResponseModel:      &models.Error{},
			ExpectedBody:       &models.Error{Error: models.ErrorUserWithIdNotFound},
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
			ResponseModel:      &models.Error{},
			ExpectedBody:       &models.Error{Error: models.ErrorSpaceAccessWithIdsNotFound},
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
			ResponseModel:      &models.Error{},
			ExpectedBody:       &models.Error{Error: models.ErrorUserWithIdNotFound},
			ExpectedStatusCode: http.StatusNotFound,
		},
	}

	tests.testDynamicIntUrlCases(t, server.PatchUserUserIdSpaces)
}
