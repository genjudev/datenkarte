package mapping

import (
	"datenkarte/internal/handlers"
	"datenkarte/internal/models"
	"datenkarte/internal/plugins"
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

func MapLineToJSON(line []string, headers []string, rule models.Rule, index int, pm *plugins.PluginManager) (map[string]interface{}, error) {
	// Execute ENTER_LINE hook for all plugins
	data := map[string]interface{}{
		"line":    line,
		"headers": headers,
		"index":   index,
	}

	if result, err := pm.ExecuteHook(plugins.ENTER_LINE, data); err != nil {
		return nil, fmt.Errorf("plugin execution failed at ENTER_LINE: %w", err)
	} else {
		data = result
	}

	// If headers contain spaces, they might need to be normalized
	normalizedHeaders := make([]string, len(headers))
	for i, header := range headers {
		// Remove any extra spaces and normalize header names
		normalizedHeaders[i] = strings.TrimSpace(header)
	}
	headers = normalizedHeaders

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
			value = line[i]

			// Execute plugins for this field if any
			if len(mapping.Plugins) > 0 {
				pluginData := map[string]interface{}{
					"value":   value,
					"header":  header,
					"mapping": mapping,
					"index":   index,
				}

				for _, pluginName := range mapping.Plugins {
					if result, err := pm.ExecuteHook(plugins.ENTER_LINE, pluginData); err != nil {
						log.Printf("Plugin %s execution failed: %v", pluginName, err)
						continue
					} else {
						// Update value with plugin result
						if v, ok := result["value"]; ok {
							value = v
						}
					}
				}
			}

			// Handle nested mapping
			if mapping.Nested != "" {
				nestedParts := strings.Split(mapping.Nested, ".")
				current := mapped
				for j, part := range nestedParts {
					if j == len(nestedParts)-1 {
						current[part] = value
					} else {
						if _, exists := current[part]; !exists {
							current[part] = make(map[string]interface{})
						}
						current = current[part].(map[string]interface{})
					}
				}
			} else {
				if mapping.InsertInto == "" {
					mapped[targetKey] = value
				}
			}

			// Execute handlers after plugins
			for _, handler := range mapping.Handlers {
				response, err := handlers.SendCommand(handler, value)
				if err != nil {
					log.Printf("%v\n", err)
					continue
				}
				mapped[targetKey] = response
			}
		}

		if !found && mapping.Required {
			return nil, fmt.Errorf("required field not found: %s", mapping.Name)
		}
	}

	// Execute EXIT_LINE hook
	exitData := map[string]interface{}{
		"mapped":  mapped,
		"line":    line,
		"headers": headers,
		"index":   index,
	}

	if result, err := pm.ExecuteHook(plugins.EXIT_LINE, exitData); err != nil {
		return nil, fmt.Errorf("plugin execution failed at EXIT_LINE: %w", err)
	} else {
		if m, ok := result["mapped"].(map[string]interface{}); ok {
			mapped = m
		}
	}

	return mapped, nil
}
