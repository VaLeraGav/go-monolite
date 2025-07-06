package category_test

import (
	"fmt"

	"go-monolite/module/category"
	"go-monolite/pkg/respond"
	"go-monolite/pkg/testinit"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCategoryIntegration(t *testing.T) {
	store := testinit.SetupStoreTest(t)
	handler := category.NewHandler(store)
	server := testinit.SetupTestServer(t, handler)
	defer server.Close()

	t.Cleanup(func() {
		err := testinit.TruncateAllTables(store.Db)
		require.NoError(t, err)
	})

	const targetUUID = "550e8400-e29b-41d4-a711-446655440002"

	validJSON := `[
		{
			"uuid": "550e8400-e29b-41d4-a711-446655440002",
			"name": "Категория a",
			"active": "Y",
			"parent_uuid": null
		},
		{
			"uuid": "550e8400-e29b-41d4-a712-446655440002",
			"name": "Категория b",
			"active": "Y",
			"parent_uuid": "550e8400-e29b-41d4-a711-446655440002"
		}
	]`

	t.Run("Create Category Validation Error", func(t *testing.T) {
		invalidJSON := `[{"name": "Категория без UUID", "active": "Y", "parent_uuid": null}]`
		resp := testinit.SendRequest(t, server.URL+"/create", "POST", invalidJSON)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var errResp struct {
			Status  string            `json:"status"`
			Message string            `json:"message"`
			Errors  map[string]string `json:"errors"`
		}
		testinit.DecodeJSON(t, resp.Body, &errResp)

		assert.Equal(t, "error", errResp.Status)
		assert.Equal(t, "ошибка в валидации поля", errResp.Message)
		assert.Equal(t, "Поле uuid обязательно для заполнения", errResp.Errors["uuid"])
	})

	t.Run("Create Category", func(t *testing.T) {
		resp := testinit.SendRequest(t, server.URL+"/create", "POST", validJSON)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response respond.Response
		testinit.DecodeJSON(t, resp.Body, &response)

		var categories []category.CategoryResponse
		testinit.MarshalUnmarshal(t, response.Data, &categories)

		require.NotEmpty(t, categories)
		expected := makeCategoryResp(1, targetUUID, "Категория a", "Y")
		assert.Equal(t, expected, categories[0])
	})

	t.Run("Get Category by UUID", func(t *testing.T) {
		resp := testinit.SendRequest(t, server.URL+"/"+targetUUID, "GET", "")
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response respond.Response
		testinit.DecodeJSON(t, resp.Body, &response)

		var cat category.CategoryResponse
		testinit.MarshalUnmarshal(t, response.Data, &cat)

		expected := makeCategoryResp(1, targetUUID, "Категория a", "Y")
		assert.Equal(t, expected, cat)
	})

	t.Run("Update Category", func(t *testing.T) {
		updateJSON := fmt.Sprintf(`{
			"uuid": "%s",
			"name": "Категория a обновленная",
			"active": "N",
			"parent_uuid": null
		}`, targetUUID)

		resp := testinit.SendRequest(t, server.URL+"/update/"+targetUUID, "PUT", updateJSON)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response respond.Response
		testinit.DecodeJSON(t, resp.Body, &response)

		var updated category.CategoryResponse
		testinit.MarshalUnmarshal(t, response.Data, &updated)

		expected := makeCategoryResp(1, targetUUID, "Категория a обновленная", "N")
		assert.Equal(t, expected, updated)
	})

	t.Run("Delete Category", func(t *testing.T) {
		resp := testinit.SendRequest(t, server.URL+"/delete/"+targetUUID, "DELETE", "")
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		respNotFound := testinit.SendRequest(t, server.URL+"/delete/"+targetUUID, "DELETE", "")
		assert.Equal(t, http.StatusNotFound, respNotFound.StatusCode)
	})
}

func makeCategoryResp(id uint, uuidStr, name, active string) category.CategoryResponse {
	uid, err := uuid.Parse(uuidStr)
	if err != nil {
		panic(err)
	}
	return category.CategoryResponse{
		ID:     id,
		UUID:   uid,
		Name:   name,
		Active: active,
		Slug:   slug.Make(name),
	}
}
