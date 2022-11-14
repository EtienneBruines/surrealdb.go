package sql

import (
	"database/sql"
	"reflect"
	"strings"
)

var DefaultTagName = "json"

// StructToNamedArgs will take a struct and pull all of its "json" (DefaultTagName)
// while ignoring the excludes
func StructToNamedArgs(entity any, excludes ...string) []sql.NamedArg {
	return StructToNamedArgsTagName(entity, DefaultTagName, excludes...)
}

func StructToNamedArgsTagName(entity any, tagname string, excludes ...string) []sql.NamedArg {
	args := []sql.NamedArg{}
	entityType := reflect.TypeOf(entity)

	// resolve the entity type and name
	for entityType.Kind() == reflect.Ptr {
		entityType = entityType.Elem()
	}

	val := reflect.ValueOf(entity)
	for i := 0; i < val.NumField(); i++ {
		field := entityType.Field(i)
		tag, ok := field.Tag.Lookup(tagname)
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
