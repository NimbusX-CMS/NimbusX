package resources_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/NimbusX-CMS/NimbusX/api/internal/db/multi_db"
	"github.com/NimbusX-CMS/NimbusX/api/internal/resources"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type TestCases []TestCase

type TestCase struct {
	name               string
	Url                string
	ID                 int
	ID2                int
	Str                string
	RequestBody        string
	ResponseModel      interface{}
	ExpectedBody       interface{}
	ExpectedStatusCode int
}

func (tc TestCases) testStaticUrlCases(t *testing.T, toTest func(ctx *gin.Context)) {
	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			w, c := setupGinTest()
			c.Request, _ = http.NewRequest(http.MethodPost, tt.Url, strings.NewReader(tt.RequestBody))
			c.Request.Header.Set("Content-Type", "application/json")

			toTest(c)

			AssertStatusCode(t, w, tt.ExpectedStatusCode)

			AssertBody(t, w, tt.ResponseModel, tt.ExpectedBody)
		})
	}
}

func (tc TestCases) testDynamicIntUrlCases(t *testing.T, toTest func(ctx *gin.Context, id int)) {
	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			w, c := setupGinTest()
			c.Request, _ = http.NewRequest(http.MethodPost, tt.Url, strings.NewReader(tt.RequestBody))
			c.Request.Header.Set("Content-Type", "application/json")

			toTest(c, tt.ID)

			AssertStatusCode(t, w, tt.ExpectedStatusCode)

			AssertBody(t, w, tt.ResponseModel, tt.ExpectedBody)
		})
	}
}

func (tc TestCases) testDynamic2IntUrlCases(t *testing.T, toTest func(ctx *gin.Context, id int, id2 int)) {
	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			w, c := setupGinTest()
			c.Request, _ = http.NewRequest(http.MethodPost, tt.Url, strings.NewReader(tt.RequestBody))
			c.Request.Header.Set("Content-Type", "application/json")

			toTest(c, tt.ID, tt.ID2)

			AssertStatusCode(t, w, tt.ExpectedStatusCode)

			AssertBody(t, w, tt.ResponseModel, tt.ExpectedBody)
		})
	}
}

func (tc TestCases) testDynamicStringUrlCases(t *testing.T, toTest func(ctx *gin.Context, str string)) {
	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			w, c := setupGinTest()
			c.Request, _ = http.NewRequest(http.MethodPost, tt.Url, strings.NewReader(tt.RequestBody))
			c.Request.Header.Set("Content-Type", "application/json")

			toTest(c, tt.Str)

			AssertStatusCode(t, w, tt.ExpectedStatusCode)

			AssertBody(t, w, tt.ResponseModel, tt.ExpectedBody)
		})
	}
}

func setupServer(t *testing.T) *resources.Server {
	testDB, err := multi_db.ConnectToSQLite(":memory:")
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		t.Error("Error connecting to database:", err)
	}
	err = testDB.EnsureTablesCreation()
	if err != nil {
		fmt.Println("Error creating tables:", err)
		t.Error("Error connecting to database:", err)
	}
	return &resources.Server{
		DB: testDB,
	}
}

func setupGinTest() (*httptest.ResponseRecorder, *gin.Context) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	return w, c
}

func AssertStatusCode(t *testing.T, w *httptest.ResponseRecorder, expectedCode int) {
	assert.Equal(t, expectedCode, w.Code)
}

func AssertBody(t *testing.T, w *httptest.ResponseRecorder, responseModel any, expectedBody interface{}) {
	if responseModel == nil {
		return
	}
	var err = json.Unmarshal(w.Body.Bytes(), responseModel)
	if err != nil {
		t.Errorf("Error parsing response body: %v\n%v", err, w.Body.String())
	}

	if !reflect.DeepEqual(responseModel, expectedBody) {
		t.Errorf("Expected %+v but got %+v", expectedBody, responseModel)
	}
}
