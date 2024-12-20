package main

import (
	"datenkarte/internal/handlers"
	"datenkarte/internal/mapping"
	"datenkarte/internal/middlewares"
	"datenkarte/internal/models"
	"datenkarte/internal/networking"
	"datenkarte/internal/plugins"
	"datenkarte/internal/validation"
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

func uploadCSV(rule models.Rule, pm *plugins.PluginManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Execute ENTER_RULE hook
		ruleData := map[string]interface{}{
			"rule_id": rule.ID,
			"type":    rule.Type,
		}

		if _, err := pm.ExecuteHook(plugins.ENTER_RULE, ruleData); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Plugin execution failed: %v", err)})
			return
		}

		queries := c.Request.URL.Query()
		dry := false
		if queries.Get("dry") != "" {
			dry = true
		}

		file, err := c.FormFile("file")
		if err != nil {
			log.Printf("%v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "could not receive file."})
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

		reader := csv.NewReader(fileOpen)
		reader.Comma = []rune(delimiter)[0]
		records, err := reader.ReadAll()
		if err != nil {
			log.Printf("%v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse CSV file"})
			return
		}

		var payloads []map[string]interface{}

		headers := records[0]
		processedRows := 0
		for i, line := range records[1:] {
			if err := validation.ValidateLine(c, line, headers, rule); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Validation failed: %v", err)})
				return
			}

			jsonPayload, err := mapping.MapLineToJSON(line, headers, rule, i, pm)
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

		// Execute EXIT_RULE hook before HTTP operations
		exitData := map[string]interface{}{
			"rule_id":        rule.ID,
			"processed_rows": processedRows,
			"payloads":       payloads,
		}

		if _, err := pm.ExecuteHook(plugins.EXIT_RULE, exitData); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Plugin execution failed: %v", err)})
			return
		}

		// start http stuff
		if dry {
			c.JSON(http.StatusOK, response)
			return
		}
		if err := networking.SendPayload(rule, response); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("HTTP Failed: %v", err)})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "success", "processed_rows": processedRows})
	}
}

func main() {
	godotenv.Load()
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v \n", err)
	}

	file, err := os.Open("config.yaml")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	var config models.Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		log.Fatalf("Failed to decode YAML: %v", err)
	}

	// Initialize plugin manager
	pm := plugins.NewPluginManager()

	// Load plugins from config
	for _, pluginPath := range config.Plugins {
		if err := pm.LoadPlugin(pluginPath); err != nil {
			log.Fatalf("Failed to load plugin %s: %v", pluginPath, err)
		}
	}

	// starting persistent handlers
	for _, handler := range config.Handlers {
		if !handler.Persistent {
			continue
		}
		if _, err := handlers.NewProcess(handler.Name); err != nil {
			log.Fatalf("%v", err)
		}
	}

	r := gin.Default()

	authGroup := r.Group("/dk/upload")
	authGroup.Use(middlewares.AuthenticationMiddleware())

	for _, rule := range config.Rules {
		authGroup.POST(rule.ID, uploadCSV(rule, pm))
	}

	log.Println("Datenkarte Started.")
	r.Run()
}
