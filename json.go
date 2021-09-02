package gorip

import (
	"errors"
	"strings"
)

func GetJSON(jsonMap map[string]interface{}, key string, args ...string) (result interface{}, error error) {
	if len(args) == 1 || len(args) > 2 {
		return nil, errors.New("args only can contain index and delimiter")
	}
	res, err := search(jsonMap, key, args...)
	if err!= nil {
		return nil, err
	}
	return res, err
}

func search(x map[string]interface{}, key string, args ...string) (result interface{}, error error) {
	s := strings.Split(key, ".")
	key = ""
	for i,v := range s[1:]{
		if i > 0 {
			key+="."
		}
		key+=v
	}
	r := x[s[0]]
	if key != "" {
		return search(r.(map[string]interface{}), key, args...)
	}
	if len(args) == 2 {
		idx := strings.Split(args[0], ".")
		hasil := ""
		for i,v:=range idx{
			if i > 0 {
				hasil += args[1]
			}
			hasil += r.(map[string]interface{})[v].(string)
		}
		return hasil, nil
	}
	return r, nil
}