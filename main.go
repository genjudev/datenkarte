package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

func RaiseBadRequest(c *gin.Context, message string, err error) {
	fmt.Println(err)
	c.JSON(http.StatusBadRequest, gin.H{"error": message})
	return
}

func buildNestedMap(key string, value interface{}) map[string]interface{} {
	keys := strings.Split(key, ".")
	m := make(map[string]interface{})
	current := m
	for i, k := range keys {
		if i == len(keys)-1 {
			current[k] = value
		} else {
			nested := make(map[string]interface{})
			current[k] = nested
			current = nested
		}
	}
	return m
}

func uploadCSV(config Config, rule Rule) gin.HandlerFunc {
	return func(c *gin.Context) {

		queries := c.Request.URL.Query()
		dry := false
		if queries.Get("dry") != "" {
			dry = true
		}

		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "could not receive file."})
			return
		}

		if file.Header.Get("content-type") != "text/csv" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "content-type not allowed."})
			return
		}

		fileOpen, err := file.Open()
		if err != nil {
			RaiseBadRequest(c, "could not open csv.", err)
			return
		}
		defer fileOpen.Close()

		delimiter := rule.Delimiter

		if delimiter == "" {
			delimiter = ";"
		}

		fmt.Printf("%v\n", delimiter)

		reader := csv.NewReader(fileOpen)
		reader.Comma = rune(delimiter[0])
		records, err := reader.ReadAll()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse CSV file"})
			return
		}

		var payloads []map[string]interface{}

		headers := records[0]
		processedRows := 0
		for i, line := range records[1:] {
			if err := ValidateLine(c, line, headers, rule); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Validation failed: %v", err)})
				return
			}

			jsonPayload, err := MapLineToJSON(line, headers, rule, i)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Mapping failed: %v", err)})
				return
			}

			payloads = append(payloads, jsonPayload)
			processedRows++
		}

		var response interface{}
		if rule.Http.PayloadKey != "" {
			response = buildNestedMap(rule.Http.PayloadKey, payloads)
		} else {
			response = payloads
		}

		// start http stuff
		if dry {
			c.JSON(http.StatusOK, response)
			return
		}
		if err := SendPayload(rule, response); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("HTTP Failed: %v", err)})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "success", "processed_rows": processedRows})
	}
}

func main() {
	godotenv.Load()
	file, err := os.Open("config.yaml")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		log.Fatalf("Failed to decode YAML: %v", err)
	}

	r := gin.Default()

	authGroup := r.Group("/dk/upload")
	authGroup.Use(AuthenticationMiddleware())

	for _, rule := range config.Rules {
		authGroup.POST(rule.ID, uploadCSV(config, rule))
	}

	log.Println("Datenkarte Started.")
	r.Run()
}
