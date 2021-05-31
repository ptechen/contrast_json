package render

import (
	"fmt"
	"strings"
)

func JsonRender(data interface{}, highlightWords []string) interface{} {

	switch data.(type) {
	case string:
		curData := data.(string)
		return render(curData, highlightWords)
	case map[string]interface{}:
		curData := data.(map[string]interface{})
		for key, val := range curData {
			curData[key] = JsonRender(val, highlightWords)
		}
	case []interface{}:
		curData := data.([]interface{})
		for idx, val := range curData {
			curData[idx] = JsonRender(val, highlightWords)
		}
	}
	return data
}

func render(data string, highlightWords []string) string {
	for _, word := range highlightWords {
		if strings.Contains(data, word) {
			data = strings.ReplaceAll(data, word, fmt.Sprintf("<a class='render'>%s</a>", word))
		}
	}
	return data
}
