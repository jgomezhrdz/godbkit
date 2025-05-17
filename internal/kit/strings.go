package kit

import (
	"bytes"
	"encoding/json"
	"unicode"
)

func JsonToStr(jsonData map[string]interface{}) string {
	var data string
	marshalledBytes, err := json.Marshal(jsonData)
	if err != nil {
		data = "{}"
	} else {
		data = string(marshalledBytes)
	}

	return data
}

func StringToJson(str string) map[string]interface{} {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(str), &data)
	if err != nil {
		data = map[string]interface{}{}
	}
	return data
}

func CamelToUnderscore(input string) string {
	var output bytes.Buffer
	for i, r := range input {
		if unicode.IsUpper(r) {
			if i > 0 {
				output.WriteByte('_')
			}
			output.WriteRune(unicode.ToLower(r))
		} else {
			output.WriteRune(r)
		}
	}
	return output.String()
}
