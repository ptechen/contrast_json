package compare

import (
	"fmt"
	"testing"
)

func TestJsonCompare(t *testing.T) {
	json1 := map[string]interface{}{"test": "123", "rrr": 0, "rrreew": map[string]interface{}{"432": "234", "321": "432"}}
	json2 := map[string]interface{}{"test": "321", "rrreew": map[string]interface{}{"432": "234"}, "dfass": []string{"ewq"}}
	res, ok := JsonCompare(json1, json2, -1)
	fmt.Println(ok)
	fmt.Println(res)
}
