package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var LOG string = os.Getenv("LOG_PATH")

func main() {

	if LOG == "" {
		LOG = "/var/lib/kubelet/pods"
	}

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/logs/:uid/:filename", func(c *gin.Context) {
		uid := c.Param("uid")
		filename := c.Param("filename")

		if uid == "" || filename == "" {
			c.String(http.StatusBadRequest, "UID & FileName are required")
			return
		}

		if strings.Contains(uid, "..") || strings.Contains(filename, "..") {
			c.String(http.StatusBadRequest, "File path doesn't allow '..' ")
			return
		}

		filePath := fmt.Sprintf("%s/%s/volumes/kubernetes.io~empty-dir/logs/%s", LOG, uid, filename)
		c.File(filePath)
	})

	router.GET("/logs/:uid", func(c *gin.Context) {
		uid := c.Param("uid")

		if uid == "" {
			c.HTML(http.StatusBadRequest, "/error.tmpl", gin.H{"error": "UID is required"})
			return
		}

		if strings.Contains(uid, "..") {
			c.HTML(http.StatusBadRequest, "/error.tmpl", gin.H{"error": "File path doesn't allow '..' "})
			return
		}

		filePath := fmt.Sprintf("%s/%s/volumes/kubernetes.io~empty-dir/logs", LOG, uid)

		files, err := ioutil.ReadDir(filePath)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "/error.tmpl", gin.H{"error": "Folder does not exists, foldername:" + uid})
			return
		}

		fileNames := []string{}

		for _, f := range files {
			fileName := f.Name()
			fileNames = append(fileNames, fileName)
		}

		c.HTML(http.StatusOK, "/index.tmpl", gin.H{"title": "List of files", "uid": uid, "data": fileNames})

	})

	router.Run(":11000")
}
