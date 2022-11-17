package sql

import (
	"database/sql"
	"database/sql/driver"
	"reflect"
	"regexp"
	"strings"
)

var DefaultTagName = "json"
var tagRE = regexp.MustCompile("[^a-zA-Z0-9-_]+")

// StructToNamedArgs will take a struct and pull all of its "json" (DefaultTagName)
// while ignoring the excludes
func StructToNamedArgs(entity any, excludes ...string) []any {
	return StructToNamedArgsTagName(entity, DefaultTagName, excludes...)
}

func StructToNamedArgsTagName(entity any, tagname string, excludes ...string) []any {
	args := []any{}
	entityType := reflect.TypeOf(entity)

	// resolve the entity type and name
	for entityType.Kind() == reflect.Ptr {
		entityType = entityType.Elem()
	}

	val := reflect.ValueOf(entity)
	for i := 0; i < val.NumField(); i++ {
		field := entityType.Field(i)
		tagValue, ok := field.Tag.Lookup(tagname)
		parts := tagRE.Split(tagValue, -1)
		if len(parts) < 1 {
			continue
		}

		tag := parts[0]
		if !ok || strings.TrimSpace(tag) == "" || SliceContains(excludes, tag) {
			continue
		}

		fieldValue := val.Field(i)
		args = append(args, sql.Named(tag, fieldValue.Interface()))
	}

	return args
}

// SliceContains utilty func that will check to see if a slice has an element
func SliceContains[T comparable](arr []T, element T) bool {
	for _, ele := range arr {
		if ele == element {
			return true
		}
	}

	return false
}

// PrepareQuery allows for queries to be multiple lines by stripping whitespace and
// replacing newlines with a whitespace
func PrepareQuery(query string) string {
	parts := strings.Split(strings.ReplaceAll(query, "\r\n", "\n"), "\n")
	final := make([]string, len(parts))

	for _, part := range parts {
		final = append(final, strings.TrimSpace(part))
	}

	return strings.Join(final, " ")
}

func NamedValuesToValues(namedValues []driver.NamedValue) []driver.Value {
	named := map[string]any{}
	for _, n := range namedValues {
		named[n.Name] = n.Value
	}

	return []driver.Value{
		named,
	}
}
