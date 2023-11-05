package resources_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/NimbusX-CMS/NimbusX/api/internal/models"
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
			name:               "DeleteSpaceSpaceId, t space does not exist",
			Url:                "/spaces/123123",
			ID:                 123123,
			ResponseModel:      &models.Error{},
			ExpectedBody:       &models.Error{Error: models.ErrorSpaceWithIdNotFound},
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
			ResponseModel:      &models.Error{},
			ExpectedBody:       &models.Error{Error: models.ErrorSpaceWithIdNotFound},
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
			ResponseModel:      &models.Error{},
			ExpectedBody:       &models.Error{Error: models.ErrorSpaceWithIdNotFound},
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
