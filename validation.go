package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func ValidateLine(c *gin.Context, line []string, headers []string, rule Rule) error {
	for _, validation := range rule.EachLine[0].Validation {
		for i, header := range headers {
			if validation.Field == header {
				value := line[i]
				switch validation.Type {
				case "number":
					if _, err := strconv.Atoi(value); err != nil {
						return fmt.Errorf("field %s must be a number, got: %s", header, value)
					}
				case "string":
					if len(value) == 0 {
						return fmt.Errorf("field %s must be a non-empty string", header)
					}
				case "email":
					if !strings.Contains(value, "@") {
						return fmt.Errorf("field %s must be a valid email, got: %s", header, value)
					}
				case "regex":
					matched, err := regexp.MatchString(validation.Pattern, value)
					if err != nil || !matched {
						return fmt.Errorf("field %s does not match pattern %s, got: %s", header, validation.Pattern, value)
					}
				default:
					return fmt.Errorf("unknown validation type: %s for field %s", validation.Type, header)
				}
			}
		}
	}
	return nil
}
