package storage_test

import (
	"fmt"
	"go-monolite/internal/module/storage"
	"go-monolite/pkg/respond"
	"go-monolite/pkg/testinit"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorageIntegration(t *testing.T) {
	store := testinit.SetupStoreTest(t)
	handler := storage.NewHandler(store)
	server := testinit.SetupTestServer(t, handler)
	defer server.Close()

	t.Cleanup(func() {
		err := testinit.TruncateAllTables(store.Db)
		require.NoError(t, err)
	})

	const storageUUID1 = "a1111111-b222-c333-d444-e55555555555"
	const storageUUID2 = "f6666666-7777-8888-9999-aaaaaaaaaaaa"
	const storageUUID3 = "f6666666-7777-8888-7777-aaaaaaaaaaaa"
	const productUUID1 = "123e4567-e89b-12d3-a456-426614174000"
	const productUUID2 = "123e4567-e89b-12d3-a456-426614174001"

	validUpsertJSON := fmt.Sprintf(`{
		"general": {
			"storages": [
				{
					"uuid": "%s",
					"name": "SPB",
					"active": "Y"
				},
				{
					"uuid": "%s",
					"name": "Backup Warehouse",
					"active": "Y"
				},
				{
					"uuid": "%s",
					"name": "Backup HH",
					"active": "Y"
				}
			]
		},
		"data": [
			{
				"product_uuid": "%s",
				"storages": [
					{
						"storage_uuid": "%s",
						"active": "Y",
						"quantity": 100
					},
					{
						"storage_uuid": "%s",
						"active": "Y",
						"quantity": 0
					},
					{
						"storage_uuid": "%s",
						"active": "Y",
						"quantity": 22
					}
				]
			},
			{
				"product_uuid": "%s",
				"storages": [
					{
						"storage_uuid": "%s",
						"active": "Y",
						"quantity": 50
					},
					{
						"storage_uuid": "%s",
						"active": "Y",
						"quantity": 23
					}
				]
			}
		]
	}`, storageUUID1, storageUUID2, storageUUID3, productUUID1, storageUUID1, storageUUID2, storageUUID3, productUUID2, storageUUID1, storageUUID2)

	t.Run("Create Storage Validation Error", func(t *testing.T) {
		invalidJSON := `{
			"general": {
				"storages": [
					{
						"name": "Склад без UUID",
						"active": "Y"
					}
				]
			},
			"data": []
		}`
		resp := testinit.SendRequest(t, server.URL+"/upsert", "POST", invalidJSON)

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

	t.Run("Upsert Storage and Product Storage", func(t *testing.T) {
		resp := testinit.SendRequest(t, server.URL+"/upsert", "POST", validUpsertJSON)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response respond.Response
		testinit.DecodeJSON(t, resp.Body, &response)

		var storageResp storage.UpsertResponse
		testinit.MarshalUnmarshal(t, response.Data, &storageResp)

		assert.NotNil(t, storageResp.Storage)
		assert.NotNil(t, storageResp.ProductStorage)
		assert.Equal(t, 3, storageResp.Storage.CountInserted)
		assert.Equal(t, 5, storageResp.ProductStorage.CountInserted)
	})
}
