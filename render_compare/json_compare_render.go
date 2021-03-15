package render_compare

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type JsonDiff struct {
	HasDiff bool
	Result  string
}

var (
	Class1 = "different_a"
	Class2 = "different_b"
)

func JsonCompareRender(left, right map[string]interface{}, n int) (string, bool) {
	diff := &JsonDiff{HasDiff: false, Result: ""}
	jsonDiffDict(left, right, 1, diff)
	if diff.HasDiff {
		if n < 0 {
			return diff.Result, diff.HasDiff
		} else {
			return processContext(diff.Result, n), diff.HasDiff
		}
	}
	return "", diff.HasDiff
}

func marshal(j interface{}) string {
	value, _ := json.Marshal(j)
	return string(value)
}

func marshalAdd(j interface{}) interface{}{
	var res interface{}
	switch j.(type) {
	case map[string]interface{}:
		if v, ok := j.(map[string]interface{}); ok {
			for i := range v {
				v[i] = marshalAdd(v[i])
			}
			return v
		}
	case []interface{}:
		if v, ok := j.([]interface{}); ok {
			for i := range v {
				v[i] = marshalAdd(v[i])
			}
			return v
		}
	default:
		res = fmt.Sprintf(`<a class='%s'>%v</a>`, Class1, j)
	}
	return res
}

func marshalSub(j interface{}) interface{} {
	var res interface{}
	switch j.(type) {
	case map[string]interface{}:
		if v, ok := j.(map[string]interface{}); ok {
			for i := range v {
				v[i] = marshalSub(v[i])
			}
			return v
		}
	case []interface{}:
		if v, ok := j.([]interface{}); ok {
			for i := range v {
				v[i] = marshalSub(v[i])
			}
			return v
		}

	default:
		res = fmt.Sprintf(`<a class='%s'>%v</a>`, Class2, j)
	}
	return res
}

func jsonDiffDict(json1, json2 map[string]interface{}, depth int, diff *JsonDiff) {
	blank := strings.Repeat(" ", (2 * (depth - 1)))
	longBlank := strings.Repeat(" ", (2 * (depth)))
	diff.Result = diff.Result + "\n" + blank + "{"
	for key, value := range json1 {
		quotedKey := fmt.Sprintf("\"%s\"", key)
		if _, ok := json2[key]; ok {
			switch value.(type) {
			case map[string]interface{}:
				if _, ok2 := json2[key].(map[string]interface{}); !ok2 {
					diff.HasDiff = true
					diff.Result = diff.Result + "\n-" + blank + quotedKey + ": " + marshal(value) + ","
					json1[key] = marshalSub(value)
					diff.Result = diff.Result + "\n+" + blank + quotedKey + ": " + marshal(json2[key])
					json2[key] = marshalAdd(json2[key])
				} else {
					diff.Result = diff.Result + "\n" + longBlank + quotedKey + ": "
					jsonDiffDict(value.(map[string]interface{}), json2[key].(map[string]interface{}), depth+1, diff)
				}
			case []interface{}:
				diff.Result = diff.Result + "\n" + longBlank + quotedKey + ": "
				if _, ok2 := json2[key].([]interface{}); !ok2 {
					diff.HasDiff = true
					diff.Result = diff.Result + "\n-" + blank + quotedKey + ": " + marshal(value) + ","
					json1[key] = marshalSub(value)
					diff.Result = diff.Result + "\n+" + blank + quotedKey + ": " + marshal(json2[key])
					json2[key] = marshalAdd(json2[key])
				} else {
					jsonDiffList(value.([]interface{}), json2[key].([]interface{}), depth+1, diff)
				}
			default:
				if !reflect.DeepEqual(value, json2[key]) {
					diff.HasDiff = true
					diff.Result = diff.Result + "\n-" + blank + quotedKey + ": " + marshal(value) + ","
					json1[key] = marshalSub(value)
					diff.Result = diff.Result + "\n+" + blank + quotedKey + ": " + marshal(json2[key])
					json2[key] = marshalAdd(json2[key])
				} else {
					diff.Result = diff.Result + "\n" + longBlank + quotedKey + ": " + marshal(value)
				}
			}
		} else {
			diff.HasDiff = true
			diff.Result = diff.Result + "\n-" + blank + quotedKey + ": " + marshal(value)
			json1[key] = marshalSub(value)
		}
		diff.Result = diff.Result + ","
	}
	for key, value := range json2 {
		if _, ok := json1[key]; !ok {
			diff.HasDiff = true
			diff.Result = diff.Result + "\n+" + blank + "\"" + key + "\"" + ": " + marshal(value) + ","
			json2[key] = marshalAdd(value)
		}
	}
	diff.Result = diff.Result + "\n" + blank + "}"
}

func jsonDiffList(json1, json2 []interface{}, depth int, diff *JsonDiff) {
	blank := strings.Repeat(" ", (2 * (depth - 1)))
	longBlank := strings.Repeat(" ", (2 * (depth)))
	diff.Result = diff.Result + "\n" + blank + "["
	size := len(json1)
	if size > len(json2) {
		size = len(json2)
	}
	for i := 0; i < size; i++ {
		switch json1[i].(type) {
		case map[string]interface{}:
			if _, ok := json2[i].(map[string]interface{}); ok {
				jsonDiffDict(json1[i].(map[string]interface{}), json2[i].(map[string]interface{}), depth+1, diff)
			} else {
				diff.HasDiff = true
				diff.Result = diff.Result + "\n-" + blank + marshal(json1[i]) + ","
				json1[i] = marshalSub(json1[i])
				diff.Result = diff.Result + "\n+" + blank + marshal(json2[i])
				json2[i] = marshalAdd(json2[i])
			}
		case []interface{}:
			if _, ok2 := json2[i].([]interface{}); !ok2 {
				diff.HasDiff = true
				diff.Result = diff.Result + "\n-" + blank + marshal(json1[i]) + ","
				json1[i] = marshalSub(json1[i])
				diff.Result = diff.Result + "\n+" + blank + marshal(json2[i])
				json2[i] = marshalAdd(json2[i])
			} else {
				jsonDiffList(json1[i].([]interface{}), json2[i].([]interface{}), depth+1, diff)
			}
		default:
			if !reflect.DeepEqual(json1[i], json2[i]) {
				diff.HasDiff = true
				diff.Result = diff.Result + "\n-" + blank + marshal(json1[i]) + ","
				json1[i] = marshalSub(json1[i])
				diff.Result = diff.Result + "\n+" + blank + marshal(json2[i])
				json2[i] = marshalAdd(json2[i])
			} else {
				diff.Result = diff.Result + "\n" + longBlank + marshal(json1[i])
			}
		}
		diff.Result = diff.Result + ","
	}
	for i := size; i < len(json1); i++ {
		diff.HasDiff = true
		diff.Result = diff.Result + "\n-" + blank + marshal(json1[i])
		json1[i] = marshalSub(json1[i])
		diff.Result = diff.Result + ","
	}
	for i := size; i < len(json2); i++ {
		diff.HasDiff = true
		diff.Result = diff.Result + "\n+" + blank + marshal(json2[i])
		json2[i] = marshalAdd(json2[i])
		diff.Result = diff.Result + ","
	}
	diff.Result = diff.Result + "\n" + blank + "]"
}

func processContext(diff string, n int) string {
	index1 := strings.Index(diff, "\n-")
	index2 := strings.Index(diff, "\n+")
	begin := 0
	end := 0
	if index1 >= 0 && index2 >= 0 {
		if index1 <= index2 {
			begin = index1
		} else {
			begin = index2
		}
	} else if index1 >= 0 {
		begin = index1
	} else if index2 >= 0 {
		begin = index2
	}
	index1 = strings.LastIndex(diff, "\n-")
	index2 = strings.LastIndex(diff, "\n+")
	if index1 >= 0 && index2 >= 0 {
		if index1 <= index2 {
			end = index2
		} else {
			end = index1
		}
	} else if index1 >= 0 {
		end = index1
	} else if index2 >= 0 {
		end = index2
	}
	pre := diff[0:begin]
	post := diff[end:]
	i := 0
	l := begin
	for i < n && l >= 0 {
		i++
		l = strings.LastIndex(pre[0:l], "\n")
	}
	r := 0
	j := 0
	for j <= n && r >= 0 {
		j++
		t := strings.Index(post[r:], "\n")
		if t >= 0 {
			r = r + t + 1
		}
	}
	if r < 0 {
		r = len(post)
	}
	return pre[l+1:] + diff[begin:end] + post[0:r+1]
}