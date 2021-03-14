package render_compare

import (
	"fmt"
	"reflect"
)

var (
	Class1 = "red"
	Class2 = "green"
)

type JsonDiff struct {
	HasDiff bool
	Result  string
}

func JsonCompareRender(left, right map[string]interface{}) {
	jsonDiffDict(left, right)
}

func marshalAdd(j interface{}) string {
	return fmt.Sprintf(`<a class='%s'>%v</a>`, Class1, j)
}

func marshalSub(j interface{}) string {
	return fmt.Sprintf(`<b class='%s'>%v</b>`, Class2, j)
}

func jsonDiffDict(json1, json2 map[string]interface{}) {
	for key, value := range json1 {
		if _, ok := json2[key]; ok {
			switch value.(type) {
			case map[string]interface{}:
				if _, ok2 := json2[key].(map[string]interface{}); !ok2 {
					json1[key] = marshalSub(value)
					json2[key] = marshalAdd(json2[key])
				} else {
					jsonDiffDict(value.(map[string]interface{}), json2[key].(map[string]interface{}))
				}
			case []interface{}:
				if _, ok2 := json2[key].([]interface{}); !ok2 {

					json1[key] = marshalSub(value)
					json2[key] = marshalAdd(json2[key])
				} else {
					jsonDiffList(value.([]interface{}), json2[key].([]interface{}))
				}
			default:
				if !reflect.DeepEqual(value, json2[key]) {

					json1[key] = marshalSub(value)
					json2[key] = marshalAdd(json2[key])
				}
			}
		} else {

			json1[key] = marshalSub(value)
		}
	}
	for key, value := range json2 {
		if _, ok := json1[key]; !ok {

			json2[key] = marshalAdd(value)
		}
	}
}

func jsonDiffList(json1, json2 []interface{}) {
	size := len(json1)
	if size > len(json2) {
		size = len(json2)
	}
	for i := 0; i < size; i++ {
		switch json1[i].(type) {
		case map[string]interface{}:
			if _, ok := json2[i].(map[string]interface{}); ok {
				jsonDiffDict(json1[i].(map[string]interface{}), json2[i].(map[string]interface{}))
			} else {

				json1[i] = marshalSub(json1[i])
				json2[i] = marshalAdd(json2[i])
			}
		case []interface{}:
			if _, ok2 := json2[i].([]interface{}); !ok2 {

				json1[i] = marshalSub(json1[i])
				json2[i] = marshalAdd(json2[i])
			} else {
				jsonDiffList(json1[i].([]interface{}), json2[i].([]interface{}))
			}
		default:
			if !reflect.DeepEqual(json1[i], json2[i]) {

				json1[i] = marshalSub(json1[i])
				json2[i] = marshalAdd(json2[i])
			}
		}
	}
	for i := size; i < len(json1); i++ {

		json1[i] = marshalSub(json1[i])
	}
	for i := size; i < len(json2); i++ {

		json2[i] = marshalAdd(json2[i])
	}
}
