package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
)

const (
	projectID  = "pie-play"
	bucketName = "example_bucket_sidup1"
)

type ClientUploader struct {
	// make a new structure ClientUploader
	cl         *storage.Client
	projectID  string
	bucketName string
	uploadPath string
}

var uploader *ClientUploader

func init() {
	// initialise things...
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/Users/siddharthupadhyay/Downloads/pie-play-4988b850a3ae.json") // set environment variable GOOGLE_APPLICATION_CREDENTIALS to the key of the service account.
	ctx := context.Background()                                                                                  // new variable called ctx which stores the context
	client, err := storage.NewClient(ctx)                                                                        // create new client using NewClient method of package storage sending ctx as a parameter
	if err != nil {
		log.Fatalf("Failed to create client: %v", err) // if there is an error, log this
	}

	uploader = &ClientUploader{
		cl:         client,
		bucketName: bucketName,
		projectID:  projectID,
		uploadPath: "test-files/",
	}

}

func main() {
	//uploader.UploadFile("notes_test/abc.txt")
	router := gin.Default()
	router.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file") // calling Formfile function of gin library on the variable c. It returns the file, which it takes from the body of the rest API. argument is the name that we want to give to the file here.
		// log.Println(file.Filename)
		if err != nil { // error handling
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		blobFile, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		err = uploader.UploadFile(blobFile, file.Filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "success",
		})
	})

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// UploadFile uploads an object.return type = error. input type = * ClientUploader. parameters of type multipart.File, string
func (c *ClientUploader) UploadFile(file multipart.File, object string) error {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := c.cl.Bucket(c.bucketName).Object(c.uploadPath + object).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	return nil
}
