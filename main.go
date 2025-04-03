package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

const uploadDir = "./storage"

func main() {
	// Ensure storage directory exists
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, os.ModePerm)
	}

	// Initialize Gin router
	r := gin.Default()

	// Serve static files
	r.Static("/static", "./public")

	// Serve the index.html at the root URL
	r.GET("/", func(c *gin.Context) {
		c.File("./public/index.html")
	})

	// Routes
	r.POST("/upload", uploadFile)
	r.GET("/files", listFiles)
	r.GET("/download/:filename", downloadFile)
	r.DELETE("/delete/:filename", deleteFile)

	// Start server
	port := "8080"
	fmt.Println("Server running on port", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// Upload file
func uploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
		return
	}

	filePath := filepath.Join(uploadDir, file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "filename": file.Filename})
}

// List all files
func listFiles(c *gin.Context) {
	files, err := os.ReadDir(uploadDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list files"})
		return
	}

	var fileList []string
	for _, file := range files {
		fileList = append(fileList, file.Name())
	}

	c.JSON(http.StatusOK, gin.H{"files": fileList})
}

// Download file
func downloadFile(c *gin.Context) {
	filename := c.Param("filename")
	filePath := filepath.Join(uploadDir, filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	c.File(filePath)
}

// Delete file
func deleteFile(c *gin.Context) {
	filename := c.Param("filename")
	filePath := filepath.Join(uploadDir, filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Delete file
	if err := os.Remove(filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}
