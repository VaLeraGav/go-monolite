package price_test

import (
	"fmt"
	"go-monolite/internal/module/price"
	"go-monolite/pkg/respond"
	"go-monolite/pkg/testinit"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPriceIntegration(t *testing.T) {
	store := testinit.SetupStoreTest(t)
	handler := price.NewHandler(store)
	server := testinit.SetupTestServer(t, handler)
	defer server.Close()

	t.Cleanup(func() {
		err := testinit.TruncateAllTables(store.Db)
		require.NoError(t, err)
	})

	const typePriceUUID1 = "a1111111-b222-c333-d444-e55555555555"
	const typePriceUUID2 = "f6666666-7777-8888-9999-aaaaaaaaaaaa"
	const typePriceUUID3 = "f6666666-7777-8888-7777-aaaaaaaaaaaa"
	const productUUID1 = "123e4567-e89b-12d3-a456-426614174000"
	const productUUID2 = "123e4567-e89b-12d3-a456-426614174001"

	validUpsertJSON := fmt.Sprintf(`{
		"general": {
			"prices": [
				{
					"uuid": "%s",
					"name": "Розничная цена",
					"active": "Y"
				},
				{
					"uuid": "%s",
					"name": "Оптовая цена",
					"active": "Y"
				},
				{
					"uuid": "%s",
					"name": "Специальная цена",
					"active": "Y"
				}
			]
		},
		"data": [
			{
				"product_uuid": "%s",
				"prices": [
					{
						"type_price_uuid": "%s",
						"active": "Y",
						"price": 1000.50
					},
					{
						"type_price_uuid": "%s",
						"active": "Y",
						"price": 800.25
					},
					{
						"type_price_uuid": "%s",
						"active": "Y",
						"price": 900.75
					}
				]
			},
			{
				"product_uuid": "%s",
				"prices": [
					{
						"type_price_uuid": "%s",
						"active": "Y",
						"price": 2000.50
					},
					{
						"type_price_uuid": "%s",
						"active": "Y",
						"price": 1500.25
					}
				]
			}
		]
	}`, typePriceUUID1, typePriceUUID2, typePriceUUID3, productUUID1, typePriceUUID1, typePriceUUID2, typePriceUUID3, productUUID2, typePriceUUID1, typePriceUUID2)

	t.Run("Create Price Validation Error", func(t *testing.T) {
		invalidJSON := `{
			"general": {
				"prices": [
					{
						"name": "Цена без UUID",
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

	t.Run("Upsert Price and Product Price", func(t *testing.T) {
		resp := testinit.SendRequest(t, server.URL+"/upsert", "POST", validUpsertJSON)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response respond.Response
		testinit.DecodeJSON(t, resp.Body, &response)

		var priceResp price.UpsertResponse
		testinit.MarshalUnmarshal(t, response.Data, &priceResp)

		assert.NotNil(t, priceResp.ProductPrice)
		assert.NotNil(t, priceResp.TypePrice)
		assert.Equal(t, 5, priceResp.ProductPrice.CountInserted)
		assert.Equal(t, 3, priceResp.TypePrice.CountInserted)
	})
}
