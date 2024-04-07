package main

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMain(m *testing.M) {
	os.Setenv("DB_NAME", "test_api")
	os.Setenv("COLLECTION_NAME", "ads")

	// Set up MongoDB connection
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Create a test database and collection
	testDB := client.Database("test_api")
	testCollection := testDB.Collection("ads")

	test_time, _ := ParseTime("2023-04-01T00:00:00.000Z")
	test_end_time, _ := ParseTime("2099-05-31T00:00:00.000Z")
	// Insert test data into the collection
	testData := []interface{}{
		bson.M{"title": "Test Ad 1", "startAt": test_time, "endAt": test_time.AddDate(0, 1, 0), "conditions": []bson.M{
			bson.M{
				"ageStart": 40,
				"ageEnd":   60,
				"gender":   []Gender{"M"},
				"country":  []Country{"TW", "JP"},
				"platform": []Platform{"android", "ios"},
			},
		}},
		bson.M{"title": "Test Ad 2", "startAt": test_time, "endAt": test_end_time, "conditions": []bson.M{
			bson.M{
				"ageStart": 20,
				"ageEnd":   30,
				"gender":   nil,
				"country":  []Country{"KR", "JP"},
				"platform": []Platform{"ios"},
			},
		}},
		// Add more test data as needed
	}
	_, err = testCollection.InsertMany(context.Background(), testData)
	if err != nil {
		log.Fatal(err)
	}

	// Run the tests
	exitCode := m.Run()

	// Delete all documents from the collection
	_, err = testCollection.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	// Exit with the same exit code as the tests
	os.Exit(exitCode)
}

func TestAdminAPISuccess(t *testing.T) {
	// Create a new Gin router instance
	router := gin.Default()

	// Define the route and associate it with the getAds handler function
	router.POST("/api/v1/ad", addAds)

	payload := []byte(`{
		"title": "AD test",
		"startAt": "2023-12-10T03:00:00.000Z",
		"endAt": "2024-12-31T16:00:00.000Z",
		"conditions": [{
			"ageStart": 30,
			"ageEnd": 40,
			"country": ["TW", "JP"],
			"platform": ["android", "ios"]
		}]
	}`)

	req, err := http.NewRequest("POST", "/api/v1/ad", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	// Call the ServeHTTP method on the router with the mock request and response recorder
	router.ServeHTTP(rr, req)

	// Check if the status code is 201
	assert.Equal(t, http.StatusCreated, rr.Code)

	expected := `{
		"title": "AD test",
		"startAt": "2023-12-10T03:00:00.000Z",
		"endAt": "2024-12-31T16:00:00.000Z",
		"conditions": [{
				"ageStart": 30,
				"ageEnd": 40,
				"gender": null,
				"country": ["TW", "JP"],
				"platform": ["android", "ios"]
			}]
		}`
	// Assert that the response body matches the expected JSON
	assert.JSONEq(t, expected, rr.Body.String())
}

func TestAdminAPIMissingRequiredField(t *testing.T) {
	// Create a new Gin router instance
	router := gin.Default()

	// Define the route and associate it with the getAds handler function
	router.POST("/api/v1/ad", addAds)

	payload := []byte(`{
		"startAt": "2023-12-10T03:00:00.000Z",
		"endAt": "2024-12-31T16:00:00.000Z",
		"conditions": [{
			"ageStart": 20,
			"ageEnd": 30,
			"country": ["TW", "JP"],
			"platform": ["android", "ios"]
		}]
	}`)

	req, err := http.NewRequest("POST", "/api/v1/ad", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	// Call the ServeHTTP method on the router with the mock request and response recorder
	router.ServeHTTP(rr, req)

	// Check if the status code is 400
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	expected := `{
		"error": "Missing required fields"
	}`
	// Assert that the response body matches the expected JSON
	assert.JSONEq(t, expected, rr.Body.String())
}

func TestAdminAPIFailedBindingData(t *testing.T) {
	// Create a new Gin router instance
	router := gin.Default()

	// Define the route and associate it with the getAds handler function
	router.POST("/api/v1/ad", addAds)

	payload := []byte(`{
		"title": 123,
		"startAt": "2023-12-10T03:00:00.000Z",
		"endAt": "2024-12-31T16:00:00.000Z",
		"conditions": [{
			"ageStart": 20,
			"ageEnd": 30,
			"country": ["TW", "JP"],
			"platform": ["android", "ios"]
		}]
	}`)

	req, err := http.NewRequest("POST", "/api/v1/ad", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	// Call the ServeHTTP method on the router with the mock request and response recorder
	router.ServeHTTP(rr, req)

	// Check if the status code is 400
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	expected := `{
		"error": "Failed binding data"
	}`
	// Assert that the response body matches the expected JSON
	assert.JSONEq(t, expected, rr.Body.String())
}

func TestPublicAPISuccess(t *testing.T) {
	//testCollection := setupTestDB()

	// Create a new Gin router instance
	router := gin.Default()

	// Define the route and associate it with the getAds handler function
	router.GET("/api/v1/ad", getAds)

	// Create a mock HTTP request to test the getAds handler
	req, err := http.NewRequest("GET", "/api/v1/ad?limit=3&age=24&gender=F&country=KR&platform=ios", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Call the ServeHTTP method on the router with the mock request and response recorder
	router.ServeHTTP(rr, req)

	// Check if the status code is OK
	assert.Equal(t, http.StatusOK, rr.Code)

	expected := `{
		"items": [
			{
				"title": "Test Ad 2",
				"endAt": "2099-05-31T00:00:00.000Z"
			}
		]
	}`
	assert.JSONEq(t, expected, rr.Body.String())

	// Clean up test data after tests
	//err = testCollection.Drop(context.Background())
	//if err != nil {
	//	log.Fatal(err)
	//}
}

func TestPublicAPIInvalidAgeParameter(t *testing.T) {
	//testCollection := setupTestDB()

	// Create a new Gin router instance
	router := gin.Default()

	// Define the route and associate it with the getAds handler function
	router.GET("/api/v1/ad", getAds)

	// Create a mock HTTP request to test the getAds handler
	req, err := http.NewRequest("GET", "/api/v1/ad?limit=3&age=abc&gender=F&country=TW&platform=ios", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Call the ServeHTTP method on the router with the mock request and response recorder
	router.ServeHTTP(rr, req)

	// Check if the status code is 400
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	expected := `{
		"error": "invalid age parameter"
	}`
	assert.JSONEq(t, expected, rr.Body.String())

	// Clean up test data after tests
	//err = testCollection.Drop(context.Background())
	//if err != nil {
	//	log.Fatal(err)
	//}
}

func TestPublicAPIInvalidGenderParameter(t *testing.T) {
	//testCollection := setupTestDB()

	// Create a new Gin router instance
	router := gin.Default()

	// Define the route and associate it with the getAds handler function
	router.GET("/api/v1/ad", getAds)

	// Create a mock HTTP request to test the getAds handler
	req, err := http.NewRequest("GET", "/api/v1/ad?limit=3&age=24&gender=H&country=TW&platform=ios", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Call the ServeHTTP method on the router with the mock request and response recorder
	router.ServeHTTP(rr, req)

	// Check if the status code is 400
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	expected := `{
		"error": "invalid gender parameter"
	}`
	assert.JSONEq(t, expected, rr.Body.String())

	// Clean up test data after tests
	//err = testCollection.Drop(context.Background())
	//if err != nil {
	//	log.Fatal(err)
	//}
}

func TestPublicAPIInvalidCountryParameter(t *testing.T) {
	//testCollection := setupTestDB()

	// Create a new Gin router instance
	router := gin.Default()

	// Define the route and associate it with the getAds handler function
	router.GET("/api/v1/ad", getAds)

	// Create a mock HTTP request to test the getAds handler
	req, err := http.NewRequest("GET", "/api/v1/ad?limit=3&age=24&gender=F&country=AAA&platform=ios", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Call the ServeHTTP method on the router with the mock request and response recorder
	router.ServeHTTP(rr, req)

	// Check if the status code is 400
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	expected := `{
		"error": "invalid country parameter"
	}`
	assert.JSONEq(t, expected, rr.Body.String())

	// Clean up test data after tests
	//err = testCollection.Drop(context.Background())
	//if err != nil {
	//	log.Fatal(err)
	//}
}

func TestPublicAPIInvalidPlatformParameter(t *testing.T) {
	//testCollection := setupTestDB()

	// Create a new Gin router instance
	router := gin.Default()

	// Define the route and associate it with the getAds handler function
	router.GET("/api/v1/ad", getAds)

	// Create a mock HTTP request to test the getAds handler
	req, err := http.NewRequest("GET", "/api/v1/ad?limit=3&age=24&gender=F&country=TW&platform=ABC", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Call the ServeHTTP method on the router with the mock request and response recorder
	router.ServeHTTP(rr, req)

	// Check if the status code is OK
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	expected := `{
		"error": "invalid platform parameter"
	}`
	assert.JSONEq(t, expected, rr.Body.String())

	// Clean up test data after tests
	//err = testCollection.Drop(context.Background())
	//if err != nil {
	//	log.Fatal(err)
	//}
}

func TestPublicAPIEmptyItems(t *testing.T) {
	//testCollection := setupTestDB()

	// Create a new Gin router instance
	router := gin.Default()

	// Define the route and associate it with the getAds handler function
	router.GET("/api/v1/ad", getAds)

	// Create a mock HTTP request to test the getAds handler
	req, err := http.NewRequest("GET", "/api/v1/ad?limit=3&age=50&gender=F&country=US&platform=ios", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Call the ServeHTTP method on the router with the mock request and response recorder
	router.ServeHTTP(rr, req)

	// Check if the status code is OK
	assert.Equal(t, http.StatusOK, rr.Code)

	expected := `{
		"items": []
	}`
	assert.JSONEq(t, expected, rr.Body.String())

	// Clean up test data after tests
	//err = testCollection.Drop(context.Background())
	//if err != nil {
	//	log.Fatal(err)
	//}
}
