package annotation

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

const (
	mark = "@tg"
)

type Tags map[string]string

func (ant Tags) MarshalJSON() (bytes []byte, err error) {

	if len(ant) == 0 {
		return json.Marshal(nil)
	}
	return json.Marshal(map[string]string(ant))
}

func (ant Tags) Merge(t Tags) Tags {

	if ant == nil {
		ant = make(Tags)
	}
	for k, v := range t {
		ant[k] = v
	}
	return ant
}

func ParseLines(text string) (tags Tags) {

	tags = make(Tags)
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		tags.Merge(ParseComment(line))
	}
	return
}

func ParseComment(comment string) (tags Tags) {

	tags = make(Tags)
	textLines := make(map[string][]string)
	comment = strings.TrimSpace(strings.TrimPrefix(comment, "//"))
	if strings.HasPrefix(comment, mark) {
		values, _ := scan(comment[len(mark):])
		for k, v := range values {
			if _, found := tags[k]; found {
				tags[k] += "," + v
			} else {
				tags[k] = v
			}
		}
	}
	for key, value := range tags {
		if !strings.HasPrefix(value, "#") {
			continue
		}
		for textKey, text := range textLines {
			if value == textKey {
				tags[key] = strings.Join(text, "\n")
			}
		}
	}
	return
}

func (ant Tags) IsSet(tagName string) (found bool) {

	_, found = ant[tagName]
	return
}

func (ant Tags) Contains(word string) (found bool) {

	for key := range ant {
		if strings.Contains(key, word) {
			return true
		}
	}
	return
}

func (ant Tags) ToDocs() (docs []string) {

	for key, value := range ant {
		docs = append(docs, fmt.Sprintf("// %s %s=`%v`", mark, key, value))
	}
	return
}

func (ant Tags) Sub(prefix string) (subTags Tags) {

	subTags = make(Tags)
	prefix = prefix + "."
	for key, value := range ant {
		if strings.HasPrefix(key, prefix) {
			subTags[strings.TrimPrefix(key, prefix)] = value
		}
	}
	return
}

func (ant Tags) Set(tagName string, values ...string) {
	ant[tagName] = strings.Join(values, ",")
}

func (ant Tags) Value(tagName string, defValue ...string) (value string) {

	var found bool
	if value, found = ant[tagName]; !found {
		value = strings.Join(defValue, " ")
	}
	return
}

func (ant Tags) ValueInt(tagName string, defValue ...int) (value int) {

	if len(defValue) != 0 {
		value = defValue[0]
	}
	if textValue, found := ant[tagName]; found {
		if newValue, err := strconv.Atoi(textValue); err == nil {
			return newValue
		}
	}
	return
}

func (ant Tags) ValueBool(tagName string, defValue ...bool) (value bool) {

	if len(defValue) != 0 {
		value = defValue[0]
	}
	if textValue, found := ant[tagName]; found {
		if newValue, err := strconv.ParseBool(textValue); err == nil {
			return newValue
		}
	}
	return
}

func (ant Tags) ToKeys(tagName, separator string, defValue ...string) map[string]int {
	return sliceStringToMap(strings.Split(ant.Value(tagName, defValue...), separator))
}

func (ant Tags) ToMap(tagName, separator, splitter string, defValue ...string) (m map[string]string) {

	m = make(map[string]string)
	pairs := strings.Split(ant.Value(tagName, defValue...), separator)
	for _, pair := range pairs {
		if kv := strings.Split(pair, splitter); len(kv) == 2 {
			m[kv[0]] = kv[1]
		}
	}
	return
}

func (ant Tags) contains(tagName string) (found bool) { // nolint

	_, found = ant[tagName]
	return
}

func sliceStringToMap(slice []string) (m map[string]int) {

	m = make(map[string]int)

	for i, v := range slice {
		m[v] = i
	}
	return
}
