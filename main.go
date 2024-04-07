package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	os.Setenv("DB_NAME", "development_api")
	os.Setenv("COLLECTION_NAME", "ads")

	router := gin.Default()
	router.GET("/api/v1/ad", getAds)
	router.POST("/api/v1/ad", addAds)

	// Start the server in a separate goroutine
	go func() {
		if err := router.Run("localhost:8080"); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Block the main goroutine indefinitely
	select {}
}

// getAds responds with the list of all ads as JSON.
func getAds(c *gin.Context) {
	_, err := getClient()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}

	// Get the current time
	currentTime := time.Now()

	// Construct a basic MongoDB query based on the query parameters
	filter := bson.M{
		"startAt": bson.M{"$lte": currentTime},
		"endAt":   bson.M{"$gte": currentTime},
		"conditions": bson.M{
			"$elemMatch": bson.M{},
		},
	}

	// Extract query parameter
	ageCondition := c.Query("age")
	genderCondition := c.Query("gender")
	countryCondition := c.Query("country")
	platformCondition := c.Query("platform")

	// Construct age query
	ageFilter := bson.M{}
	if ageCondition != "" {
		age, err := strconv.Atoi(ageCondition)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid age parameter"})
			return
		}
		ageFilter = bson.M{
			"ageStart": bson.M{"$lte": age},
			"ageEnd":   bson.M{"$gte": age},
		}
	}

	// Construct gender query
	genderFilter := bson.M{}
	if genderCondition != "" {
		genderEnum := Gender(genderCondition)
		if !genderEnum.IsValid() {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid gender parameter"})
			return
		}
		genderFilter = bson.M{
			"$or": []bson.M{
				{"gender": bson.M{"$in": []Gender{genderEnum}}},
				{"gender": nil}, // Empty array
			},
		}
	}

	// Construct country query
	countryFilter := bson.M{}
	if countryCondition != "" {
		countryEnum := Country(countryCondition)
		if !countryEnum.IsValid() {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid country parameter"})
			return
		}
		countryFilter = bson.M{
			"$or": []bson.M{
				{"country": bson.M{"$in": []Country{countryEnum}}},
				{"country": nil}, // Empty array
			},
		}
	}

	// Construct platform query
	platformFilter := bson.M{}
	if platformCondition != "" {
		platformEnum := Platform(platformCondition)
		if !platformEnum.IsValid() {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid platform parameter"})
			return
		}
		platformFilter = bson.M{
			"$or": []bson.M{
				{"platform": bson.M{"$in": []Platform{platformEnum}}},
				{"platform": nil}, // Empty array
			},
		}
	}

	// Combine all condition filters into a single filter for $elemMatch
	conditionFilter := bson.M{
		"$elemMatch": bson.M{
			"$and": []bson.M{ageFilter, genderFilter, countryFilter, platformFilter},
		},
	}
	filter["conditions"] = conditionFilter

	// Apply filter in db
	cursor, err := dbCol.Find(context.Background(), filter)
	if err != nil {
		// Handle error
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	defer cursor.Close(context.Background())

	// Define http response body element
	displayAds := DisplayAds{
		Items: []AdItem{},
	}

	// Iterate over the cursor to retrieve the documents
	for cursor.Next(context.Background()) {
		var ad Advertisement
		if err := cursor.Decode(&ad); err != nil {
			// Handle error
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "error decoding documents in data type"})
			return
		}
		// Append an ad to displayAds.Items
		displayAds.Items = append(displayAds.Items, AdItem{Title: ad.Title, EndAt: ad.EndAt})
	}
	if err := cursor.Err(); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Handle condition where no ads is found for specified query
	if len(displayAds.Items) == 0 {
		c.IndentedJSON(http.StatusOK, displayAds)
		return
	}

	// apply pagination
	sort.Slice(displayAds.Items, func(i, j int) bool {
		return displayAds.Items[i].EndAt.Before(displayAds.Items[j].EndAt)
	})

	offset, _ := strconv.Atoi(c.Query("offset"))
	if offset >= len(displayAds.Items) {
		c.IndentedJSON(http.StatusOK, DisplayAds{
			Items: []AdItem{},
		})
		return
	}

	limit, err_l := strconv.Atoi(c.Query("limit"))
	endIndex := offset
	if err_l == nil {
		endIndex += limit
		if endIndex > len(displayAds.Items) {
			endIndex = len(displayAds.Items)
		}
	} else {
		endIndex = len(displayAds.Items)
	}
	displayAds.Items = displayAds.Items[offset:endIndex]
	//println(displayAds.Items)

	c.IndentedJSON(http.StatusOK, displayAds)
}

// addAds adds an ad from JSON received in the request body.
func addAds(c *gin.Context) {
	// connect to database
	_, err := getClient()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}

	var newAd Advertisement

	// Call BindJSON to bind the received JSON to newAd.
	if err := c.ShouldBindJSON(&newAd); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Failed binding data"})
		return
	} else {
		// Check if required fields are present
		if newAd.Title == "" || newAd.StartAt.IsZero() || newAd.EndAt.IsZero() {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
			return
		}
	}

	// Add the new ad to the db.
	_, err = dbCol.InsertOne(context.Background(), newAd)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, "Failed adding data to database")
		return
	}
	c.IndentedJSON(http.StatusCreated, newAd)
}
