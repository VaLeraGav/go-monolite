package product_test

// import (
// 	"go-monolite/internal/src/product"
// 	"go-monolite/pkg/testinit"
// 	"testing"
// )

// // 	t.Run("Create Product", ...) // как у тебя
// // t.Run("Get Product", ...)
// // t.Run("Update Product", ...)
// // t.Run("Delete Product", ...)
// //(GET → 404)

// func TestProductIntegration(t *testing.T) {
// 	store := testinit.SetupStoreTest(t)
// 	Handler := product.NewHandler(store)
// 	server := testinit.SetupTestServer(t, Handler)
// 	defer server.Close()

// 	// jsonData := `
// 	// 	{
// 	// 		"code": 12311,
// 	// 		"name": "лампа 1",
// 	// 		"weight": 1.25,
// 	// 		"width": 10.5,
// 	// 		"active": "Y",
// 	// 		"length": 20.0,
// 	// 		"height": 5.0,
// 	// 		"uuid": "123e4567-e89b-12d3-a455-426614174001",
// 	// 		"category_uuid": "550e8400-e29b-41d4-a716-446655440002",
// 	// 		"volume": 1102.5,
// 	// 		"unit": "шт",
// 	// 		"property" : {
// 	// 			"weight":1.25
// 	// 		}
// 	// 	}
// 	// `

// 	t.Run("Create Product", func(t *testing.T) {
// 		// // Create request
// 		// req, err := http.NewRequest("POST", server.URL+"/create", bytes.NewBuffer([]byte(jsonData)))
// 		// require.NoError(t, err)
// 		// req.Header.Set("Content-Type", "application/json")

// 		// // Send request
// 		// resp, err := http.DefaultClient.Do(req)
// 		// require.NoError(t, err)

// 		// bodyBytes, err := io.ReadAll(resp.Body)
// 		// bodyStr := string(bodyBytes)
// 		// t.Logf("Response body: %s", bodyStr)

// 		// require.NoError(t, err)
// 		// defer resp.Body.Close()

// 		// // Check response
// 		// assert.Equal(t, http.StatusCreated, resp.StatusCode)
// 	})

// }

// // 	t.Run("Get Product by UUID", func(t *testing.T) {
// // 		// Create request
// // 		req, err := http.NewRequest("GET", server.URL+"/"+productUUID.String(), nil)
// // 		require.NoError(t, err)

// // 		// Send request
// // 		resp, err := http.DefaultClient.Do(req)
// // 		require.NoError(t, err)
// // 		defer resp.Body.Close()

// // 		// Check response
// // 		assert.Equal(t, http.StatusOK, resp.StatusCode)

// // 		// Parse response
// // 		var response map[string]interface{}
// // 		err = json.NewDecoder(resp.Body).Decode(&response)
// // 		require.NoError(t, err)

// // 		// Verify product data
// // 		data := response["data"].(map[string]interface{})
// // 		assert.Equal(t, productUUID.String(), data["uuid"])
// // 		assert.Equal(t, productDto.Name, data["name"])
// // 	})

// // 	t.Run("Update Product", func(t *testing.T) {
// // 		// Update product data
// // 		updatedProduct := productDto
// // 		updatedProduct.Name = "Updated Product"
// // 		updatedProduct.Price = 200.75

// // 		// Convert to JSON
// // 		jsonData, err := json.Marshal(updatedProduct)
// // 		require.NoError(t, err)

// // 		// Create request
// // 		req, err := http.NewRequest("PUT", server.URL+"/update/"+productUUID.String(), bytes.NewBuffer(jsonData))
// // 		require.NoError(t, err)
// // 		req.Header.Set("Content-Type", "application/json")

// // 		// Send request
// // 		resp, err := http.DefaultClient.Do(req)
// // 		require.NoError(t, err)
// // 		defer resp.Body.Close()

// // 		// Check response
// // 		assert.Equal(t, http.StatusOK, resp.StatusCode)
// // 	})

// // 	t.Run("Get Product List", func(t *testing.T) {
// // 		// Create request
// // 		req, err := http.NewRequest("GET", server.URL+"/", nil)
// // 		require.NoError(t, err)

// // 		// Send request
// // 		resp, err := http.DefaultClient.Do(req)
// // 		require.NoError(t, err)
// // 		defer resp.Body.Close()

// // 		// Check response
// // 		assert.Equal(t, http.StatusOK, resp.StatusCode)

// // 		// Parse response
// // 		var response map[string]interface{}
// // 		err = json.NewDecoder(resp.Body).Decode(&response)
// // 		require.NoError(t, err)

// // 		// Verify list is not empty
// // 		data := response["data"].([]interface{})
// // 		assert.Greater(t, len(data), 0)
// // 	})

// // 	t.Run("Delete Product", func(t *testing.T) {
// // 		// Create request
// // 		req, err := http.NewRequest("DELETE", server.URL+"/delete/"+productUUID.String(), nil)
// // 		require.NoError(t, err)

// // 		// Send request
// // 		resp, err := http.DefaultClient.Do(req)
// // 		require.NoError(t, err)
// // 		defer resp.Body.Close()

// // 		// Check response
// // 		assert.Equal(t, http.StatusOK, resp.StatusCode)

// // 		// Verify product is deleted
// // 		getReq, err := http.NewRequest("GET", server.URL+"/"+productUUID.String(), nil)
// // 		require.NoError(t, err)

// // 		getResp, err := http.DefaultClient.Do(getReq)
// // 		require.NoError(t, err)
// // 		defer getResp.Body.Close()

// // 		assert.Equal(t, http.StatusNotFound, getResp.StatusCode)
// // 	})
// // }
