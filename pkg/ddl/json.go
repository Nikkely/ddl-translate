package ddl

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/nikkely/ddl-translate/pkg/translate"
)

type JSONObj struct {
	data map[string]interface{}
}

func NewJSONObj(data []byte) (*JSONObj, error) {
	var d map[string]interface{}
	if err := json.Unmarshal(data, &d); err != nil {
		return nil, err
	}
	return &JSONObj{data: d}, nil
}

func NewJSONObjFromFile(path string) (*JSONObj, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var d map[string]interface{}
	if err = json.Unmarshal(raw, &d); err != nil {
		return nil, err
	}
	return &JSONObj{data: d}, nil
}

func (j JSONObj) ToString() (string, error) {
	data, err := json.Marshal(j.data)
	if err != nil {
		return "", err
	}
	str := string(data)
	return str, nil
}

// Translate translate values detected by keys
// keyQuery spec is ".xxx.yyy"
func (j JSONObj) Translate(keys []string, translater translate.Translater) error {
	for _, key := range keys {
		err := j.applyWithKey(key, translater)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: key:%s cant parsed beacause of %v", key, err.Error())
			continue
		}
	}
	return nil
}

func (j JSONObj) applyWithKey(keyQuery string, translater translate.Translater) error {
	applier := func(value interface{}) interface{} {
		switch v := value.(type) {
		case string:
			result, err := translater.Run(v)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
			}
			return result
		}
		return nil
	}
	return applyWithKeyRecursively(j.data, strings.Split(keyQuery, "."), applier)
}

func applyWithKeyRecursively(jsonData map[string]interface{}, keys []string, f func(value interface{}) interface{}) error {
	if len(keys) == 0 {
		return fmt.Errorf("invalid key")
	}

	var value, ok = jsonData[keys[0]]
	if !ok {
		return fmt.Errorf("no such key: %s", keys[0])
	}

	if len(keys) == 1 {
		switch v := value.(type) {
		case map[string]interface{}:
			applyRecursively(v, f)
		case []interface{}:
			var newValue []interface{}
			for _, x := range v {
				newValue = append(newValue, f(x))
			}
			jsonData[keys[0]] = newValue
		default:
			jsonData[keys[0]] = f(value)
		}
		return nil
	}

	switch data := value.(type) {
	case map[string]interface{}:
		// object
		return applyWithKeyRecursively(data, keys[1:], f)
	case []interface{}:
		// array
		var errMsg string
		for _, value := range data {
			switch v := value.(type) {
			case map[string]interface{}:
				if err := applyWithKeyRecursively(v, keys[1:], f); err != nil {
					errMsg += err.Error()
				}
			default:
				errMsg += fmt.Sprintf("unexpected value; keys:%v", keys[1:])
			}
		}
		if errMsg != "" {
			return fmt.Errorf(errMsg)
		}
	default:
		return fmt.Errorf("no object")
	}
	return nil
}

func applyRecursively(obj map[string]interface{}, f func(v any) any) error {
	for key, value := range obj {
		switch typedValue := value.(type) {
		case map[string]interface{}:
			applyRecursively(typedValue, f)
		case []interface{}:
			var newValue []interface{}
			for _, x := range typedValue {
				newValue = append(newValue, f(x))
			}
			obj[key] = newValue
		default:
			obj[key] = f(value)
		}
	}
	return nil
}
