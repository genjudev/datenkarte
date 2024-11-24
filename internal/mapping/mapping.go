package mapping

import (
	"datenkarte/internal/handlers"
	"datenkarte/internal/models"
	"fmt"
    "log"
	"strings"
)

func stringInSlice(search string, list []string) bool {
	for _, item := range list {
		if search == item {
			return true
		}
	}
	return false
}

func MapLineToJSON(line []string, headers []string, rule models.Rule, index int) (map[string]interface{}, error) {
	mapped := make(map[string]interface{})
	for _, mapping := range rule.EachLine[0].Map {
		found := false
		targetKey := mapping.To
		if targetKey == "" {
			targetKey = mapping.Name
		}

		if mapping.Required {
			ok := stringInSlice(mapping.Name, headers)
			if !ok {
				return nil, fmt.Errorf("required header not found: %s", mapping.Name)
			}

		}

		var value interface{}

		for i, header := range headers {
			if mapping.Name != header {
				continue
			}
			found = true
			// Handle nested mapping
			if mapping.Nested != "" {
				nestedParts := strings.Split(mapping.Nested, ".")
				current := mapped
				for j, part := range nestedParts {
					if j == len(nestedParts)-1 {
						current[part] = line[i]
					} else {
						if _, exists := current[part]; !exists {
							current[part] = make(map[string]interface{})
						}
						current = current[part].(map[string]interface{})
					}
				}
			} else {
				value = line[i]
				if mapping.InsertInto == "" {
					mapped[targetKey] = value
				}
			}
			for _, handler := range mapping.Handlers {
				response, err := handlers.SendCommand(handler, value)
				if err != nil {
					log.Printf("%v\n", err)
					continue
				}
                mapped[targetKey] = response
			}
		}
		if !found && mapping.Fill != nil {
			switch mapping.Fill.Type {
			case "string":
				val := mapping.Fill.Value
				if val == "row_number" {
					value = fmt.Sprintf("%s%d", mapping.Fill.Prefix, index)
				} else {
					value = val
				}

				mapped[targetKey] = value
			case "array":
				if val, ok := mapping.Fill.Value.([]interface{}); ok {
					value = val
					mapped[targetKey] = val
				}
			}
		}

		if mapping.InsertInto != "" && found {
			existing, exists := mapped[mapping.InsertInto]
			if !exists {
				mapped[mapping.InsertInto] = []interface{}{}
				existing = mapped[mapping.InsertInto]
			}

			existingSlice, ok := existing.([]interface{})
			if !ok {
				return nil, fmt.Errorf("target %s is not an array, cannot inser into", mapping.InsertInto)
			}

			if arrayValue, ok := value.([]interface{}); ok {
				mapped[mapping.InsertInto] = append(existingSlice, arrayValue...)
			} else {
				mapped[mapping.InsertInto] = append(existingSlice, value)
			}
		}
	}
	return mapped, nil
}
