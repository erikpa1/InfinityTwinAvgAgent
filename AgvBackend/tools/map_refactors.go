package tools

func ReplaceStringsInInterfaceMap(data map[string]interface{}, replacements map[string]string) {
	for key, value := range data {
		switch v := value.(type) {
		case string:
			if newVal, exists := replacements[v]; exists {
				data[key] = newVal
			}
		case map[string]interface{}:
			ReplaceStringsInInterfaceMap(v, replacements)
		case []interface{}:
			for i, item := range v {
				switch itemVal := item.(type) {
				case string:
					if newVal, exists := replacements[itemVal]; exists {
						v[i] = newVal
					}
				case map[string]interface{}:
					ReplaceStringsInInterfaceMap(itemVal, replacements)
				}
			}
		}
	}
}
