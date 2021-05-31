package render

import (
	"fmt"
	"testing"
)

func TestJsonRender(t *testing.T) {
	data := map[string]interface{}{"ss": "123321", "44": []interface{}{"1233334", "3344"}}
	h := []string{"33"}
	d := JsonRender(data, h)
	fmt.Println(d)
}