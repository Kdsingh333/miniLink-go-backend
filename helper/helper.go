package helper

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/Kdsingh333/miniLink-go-backend/database"
	"github.com/Kdsingh333/miniLink-go-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/teris-io/shortid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var db *mongo.Collection 
var ctx = context.TODO()
var baseUrl = "https://url-shortener-service-3t2m.onrender.com/"

func init() {
	db = database.Setup()
}
func Str(c *gin.Context){
	c.JSON(http.StatusOK, gin.H{
		"msg":"Hey you are in root path",
	})
}
func Shorten(c *gin.Context) {
	c.Header("Context-Type", "application/x-www-form-urlencoded")
	// c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST")
	c.Header("Access-Control-Allow-Headers", "Content-Type")
	var body models.ShortenBody
	if err := c.BindJSON(&body); err != nil {
		fmt.Println(body)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	_, urlErr := url.ParseRequestURI(body.LongUrl)
	if urlErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": urlErr.Error(),
		})
		return
	}
	urlCode, idErr := shortid.Generate()
	if idErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error occurs during generating short id",
		})
		return
	}

	var result bson.M
	queryErr := db.FindOne(ctx, bson.D{{Key: "urlCode", Value: urlCode}}).Decode(&result)

	if queryErr != mongo.ErrNoDocuments {
	 fmt.Print(5);
		c.JSON(http.StatusInternalServerError, gin.H{"error": queryErr.Error()})
		return
	}

	if len(result)>0{
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Code in use: %s", urlCode)})
		return
	}

	var date = time.Now()
	var expires = date.AddDate(0,0,5)
	var newUrl = baseUrl + urlCode
	var docId = primitive.NewObjectID()
	newDoc := &models.UrlDoc{
		ID:        docId,
		UrlCode:   urlCode,
		LongUrl:   body.LongUrl,
		ShortUrl:  newUrl,
		CreatedAt: time.Now(),
	    ExperiesAt: expires,
	}

	_,err:= db.InsertOne(ctx,newDoc)
	if err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"error":"Error occur during insertion new doc",
		})
		return
	}
     
	c.JSON(http.StatusCreated,gin.H{
		"longUrl":body.LongUrl,
		"newUrl": newUrl,
		"expires": expires.Format("2006-01-02 15:04:05"),
		"db_id": docId,
	})
}

func Redirect(c *gin.Context) {
  code := c.Param("code")

  var result bson.M
  queryErr := db.FindOne(ctx,bson.D{{Key: "urlCode", Value: code}}).Decode(&result)

  if queryErr != nil{
	if queryErr == mongo.ErrNoDocuments{
		c.JSON(http.StatusBadRequest,gin.H{
		"error": fmt.Sprintf("No URL with code: %s", code)})
		return
	}else {
		c.JSON(http.StatusInternalServerError,gin.H{
			"error":queryErr.Error(),
		})
		return
	}
  }
  log.Print(result["longUrl"])
  var longUrl = fmt.Sprint(result["longUrl"])
  c.Redirect(http.StatusPermanentRedirect,longUrl)


}

func Custom(c *gin.Context) {
	var body models.CustomBody
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, urlErr := url.ParseRequestURI(body.LongUrl)
	if urlErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL"})
		return
	}
	var length = len(body.CustomCode)
	if length < 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Custom code should be more than 3 characters"})
		return
	}

	var result bson.M
	queryErr := db.FindOne(ctx, bson.D{{Key: "urlCode", Value: body.CustomCode}}).Decode(&result)

	if queryErr != nil {
		if queryErr != mongo.ErrNoDocuments {
			c.JSON(http.StatusInternalServerError, gin.H{"error": queryErr.Error()})
			return
		}
	}

	if len(result) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Code in use: %s", body.CustomCode)})
		return
	}

	var date = time.Now()
	var expires = date.AddDate(0, 0, 5)
	var newUrl = baseUrl + body.CustomCode
	var docId = primitive.NewObjectID()

	newDoc := &models.UrlDoc{
		ID:        docId,
		UrlCode:   body.CustomCode,
		LongUrl:   body.LongUrl,
		ShortUrl:  newUrl,
		CreatedAt: time.Now(),
		ExperiesAt: expires,
	}

	_, err := db.InsertOne(ctx, newDoc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"newUrl":  newUrl,
		"expires": expires.Format("2006-01-02 15:04:05"),
		"db_id":   docId,
	})
}
