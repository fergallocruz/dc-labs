package main

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/CodersSquad/dc-labs/challenges/final/controller"
	"github.com/CodersSquad/dc-labs/challenges/final/scheduler"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	//"google.golang.org/genproto/protobuf/api"
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
var countTests int
var jobs = make(chan scheduler.Job)

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
		//log.Println(file.Filename)

		filename := filepath.Base(file.Filename)
		if err := c.SaveUploadedFile(file, filename); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}
		//imagePath := "/" + filename
		image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
		f, err := os.Open("./" + filename)
		defer f.Close()
		fmt.Println(filename)
		if err != nil {
			log.Println(err)
		}
		img, _, err := image.DecodeConfig(f)
		if err != nil {
			c.JSON(http.StatusOK, ErrorMessageResponse("There was an error opening image"))
		} else {
			c.JSON(http.StatusOK, SuccessUploadResponse(filename, img.Width, img.Height))
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

		c.JSON(http.StatusOK, map[string]interface{}{
			"Worker": controller.Nodes[name].Name,
			"Tags":   controller.Nodes[name].Tags,
			"Status": controller.Nodes[name].Status,
			"Usage":  strconv.Itoa(controller.Nodes[name].Usage) + "%",
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
		sampleJob := scheduler.Job{Address: "localhost:50051", RPCName: "test"}
		jobs <- sampleJob
		time.Sleep(time.Second * 5)
		name := controller.GetWorker(countTests)
		//name := controller.GetWorker(countTests)
		c.JSON(http.StatusOK, map[string]interface{}{
			"Workload": "test",
			"Job ID":   countTests,
			"Status":   "Scheduling",
			"Result":   "Done by " + name,
		})
		countTests++
	}
}

func FilterHandler(c *gin.Context) {
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
		//log.Println(file.Filename)

		filename := filepath.Base(file.Filename)

		if err := c.SaveUploadedFile(file, filename); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}

		workloadId := c.Params.ByName("workload-id")
		filter := c.Params.ByName("filter")
		//here it must do the filter
		sampleJob := scheduler.Job{Address: "localhost:50051", RPCName: "test"} //rpc name filter
		jobs <- sampleJob

		time.Sleep(time.Second * 5)
		//name := controller.GetWorker(countTests)
		c.JSON(http.StatusOK, map[string]interface{}{

			"Workload ID": workloadId,
			"Filter": filter,
			"Job ID": countTests, //revisar este index
			"Status": "Scheduling",
			"Results": "http://localhost:8080/results/"+workloadId,
		})
		countTests++
	}
}

func ResultsHandler(c *gin.Context) {
	auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
	token := auth[1]
	_, ok := DB[token]
	if !ok {
		c.JSON(http.StatusOK, ErrorMessageResponse("This token doesn't exist"))
	} else {
		workloadId := c.Param("workloadId")
		sampleJob := scheduler.Job{Address: "localhost:50051", RPCName: "test"}
		jobs <- sampleJob
		time.Sleep(time.Second * 5)
		//name := controller.GetWorker(countTests)
		c.JSON(http.StatusOK, map[string]interface{}{
			"Workload": workloadId,
			"Job ID":   countTests,
			"Status":   "Scheduling",
			"Result":   "results/" + workloadId,
		})
		countTests++
	}
}
func Download(c *gin.Context) {
	auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
	token := auth[1]
	_, ok := DB[token]
	if !ok {
		c.JSON(http.StatusOK, ErrorMessageResponse("This token doesn't exist"))
	} else {
		sampleJob := scheduler.Job{Address: "localhost:50051", RPCName: "test"}
		jobs <- sampleJob
		time.Sleep(time.Second * 5)
		name := controller.GetWorker(countTests)
		//name := controller.GetWorker(countTests)
		c.JSON(http.StatusOK, map[string]interface{}{
			"Workload": "test",
			"Job ID":   countTests,
			"Status":   "Scheduling",
			"Result":   "Done by " + name,
		})
		countTests++
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
	/*workers:= ""
	for _,u := range controller.Nodes{
			workers = workers+fmt.Println(u.Name, " ", u.Status, " ", u.Usage,"%")
	}*/
	return gin.H{
		"message": "Hi " + username + ", the DPIP System is Up and Running",
		"time":    time.Now().Format("2006-01-02T15:04:05+07:00"),
		"Workers": controller.Nodes,
	}
}
func SuccessUploadResponse(image string, width int, height int) gin.H {
	return gin.H{
		"message": "Image: " + image + " uploaded succefully",
		"time":    time.Now().Format("2006-01-02T15:04:05+07:00"),
		"size":    strconv.Itoa(width) + "x" + strconv.Itoa(height),
	}
}

func main() {
	log.Println("Welcome to the Distributed and Parallel Image Processing System")
	// Start Controller
	go controller.Start()
	countTests = 0

	// Send sample jobs
	//sampleJob := scheduler.Job{Address: "localhost:50051", RPCName: "test"}

	// Start Scheduler

	go scheduler.Start(jobs)
	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("/results", true)))
	router.GET("/login", LoginHandler)
	router.GET("/logout", LogoutHandler)
	router.GET("/status", StatusHandler)
	router.GET("/workloads/test", WorkloadsHandler)
	router.POST("/workloads/filter", FilterHandler)
	router.POST("/upload", UploadHandler) // worker-token auth desde los workers
	router.GET("/status/:worker", StatusWorkerHandler)
	router.GET("/results/:workloadId", ResultsHandler) //this func not done
	router.GET("/download", Download) // worker-token auth y ver cómo acceder desde los workers
	go router.Run(":8080")
	// Send sample jobs
	sampleJob := scheduler.Job{Address: "localhost:50051", RPCName: ""}
	for {
		//sampleJob.RPCName = fmt.Sprint()
		if sampleJob.RPCName == "test" {
			jobs <- sampleJob
		}
		time.Sleep(time.Second * 2)
	}

}