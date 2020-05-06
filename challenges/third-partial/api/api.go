package api

import (
	"encoding/base64"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/CodersSquad/dc-labs/challenges/third-partial/controller"
	"github.com/gin-gonic/gin"
)

type LoginStruct struct {
	User     string
	Password string
	Token    string
}
type Login struct {
	User     string `form:"user" json:"user" xml:"user"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

var DB = make(map[string]LoginStruct)

func Start() {
	router := gin.Default()
	router.GET("/login", LoginHandler)
	router.GET("/logout", LogoutHandler)
	router.GET("/status", StatusHandler)
	router.GET("/workloads/test", WorkloadsHandler)
	router.POST("/upload", UploadHandler)

	router.GET("/status/:worker", StatusWorkerHandler)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

// Login User handler
func LoginHandler(c *gin.Context) {
	var form Login
	// This will infer what binder to use depending on the content-type header.
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username := form.User
	password := form.Password
	taken := false
	for i := range DB {
		if DB[i].User == username {
			taken = true
		}
	}

	if taken || username == "" {
		c.JSON(http.StatusOK, ErrorMessageResponse("This username is taken"))
	} else {
		auth := username + ":" + password
		tokenString := base64.StdEncoding.EncodeToString([]byte(auth))
		DB[tokenString] = LoginStruct{
			User:     username,
			Password: password,
			Token:    tokenString,
		}
		c.JSON(http.StatusOK, SuccessLoginResponse(username, tokenString))
	}

}

// Log Out handler
func LogoutHandler(c *gin.Context) {
	auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
	token := auth[1]
	_, ok := DB[token]
	if !ok {
		c.JSON(http.StatusOK, ErrorMessageResponse("This token doesn't exist"))
	} else {
		username := DB[token].User
		c.JSON(http.StatusOK, SuccessLogoutResponse(username))
		delete(DB, token)
	}

}

// Status Handler
func StatusHandler(c *gin.Context) {
	auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
	token := auth[1]
	_, ok := DB[token]
	if !ok {
		c.JSON(http.StatusOK, ErrorMessageResponse("This token doesn't exist"))
	} else {
		username := DB[token].User
		c.JSON(http.StatusOK, SuccessStatusResponse(username))
	}
}

// Upload handler
func UploadHandler(c *gin.Context) {
	auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
	token := auth[1]
	_, ok := DB[token]
	if !ok {
		c.JSON(http.StatusOK, ErrorMessageResponse("This token doesn't exist"))
	} else {
		file, err := c.FormFile("data")
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}

		filename := filepath.Base(file.Filename)
		if err := c.SaveUploadedFile(file, filename); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}

		f, err := os.Open(filename)
		if err != nil {
			c.JSON(http.StatusOK, ErrorMessageResponse("There was an error with the image"))
		}
		image, _, err := image.DecodeConfig(f)
		if err != nil {
			c.JSON(http.StatusOK, ErrorMessageResponse("There was an error opening image"))
		} else {
			c.JSON(http.StatusOK, SuccessUploadResponse(filename, image.Width, image.Height))
		}
	}

}

func StatusWorkerHandler(c *gin.Context) {
	auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
	token := auth[1]
	name := c.Param("worker")
	_, ok := DB[token]
	if !ok {
		c.JSON(http.StatusOK, ErrorMessageResponse("This token doesn't exist"))
	} else {
		c.JSON(http.StatusOK, gin.H{
			"Worker": name,
			"Tags":   controller.Nodes[name]["tags"],
			"Status": controller.Nodes[name]["status"],
			"Usage":  strconv.Itoa(999) + "%",
		})
	}
}
func WorkloadsHandler(c *gin.Context) {
	auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
	token := auth[1]
	_, ok := DB[token]
	if !ok {
		c.JSON(http.StatusOK, ErrorMessageResponse("This token doesn't exist"))
	} else {
		username := DB[token].User
		c.JSON(http.StatusOK, SuccessStatusResponse(username))
	}
}

func SuccessLoginResponse(username string, token string) gin.H {
	return gin.H{
		"message": "Hi " + username + ", welcome to the DPIP System",
		"token":   token,
	}
}
func SuccessLogoutResponse(username string) gin.H {
	return gin.H{
		"message": "Bye " + username + ", your token has been revoked",
	}
}

// ErrorMessageResponse Request response object ready for errors.
func ErrorMessageResponse(message string) gin.H {
	return gin.H{
		"status": "error",
		"data": gin.H{
			"message": message,
		},
	}
}
func SuccessStatusResponse(username string) gin.H {
	return gin.H{
		"message": "Hi " + username + ", the DPIP System is Up and Running",
		"time":    time.Now().Format("2006-01-02T15:04:05+07:00"),
	}
}
func SuccessUploadResponse(image string, width int, height int) gin.H {
	return gin.H{
		"message": "Image: " + image + " uploaded succefully",
		"time":    time.Now().Format("2006-01-02T15:04:05+07:00"),
		"size":    strconv.Itoa(width) + "x" + strconv.Itoa(height),
	}
}
