package render_compare

import (
	"fmt"
	"testing"
)

func TestJsonCompareRender(t *testing.T) {
	json1 := map[string]interface{}{"test": "123", "rrr": 0, "rrreew": map[string]interface{}{"432": "234", "321": "432"}}
	json2 := map[string]interface{}{"test": "321", "rrreew": map[string]interface{}{"432": "234"}, "dfass": []string{"ewq"}}
	JsonCompareRender(json1, json2)
	fmt.Println(json1)
	fmt.Println(json2)
}
