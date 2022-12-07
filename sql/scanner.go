package sql

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

var uuidRE = regexp.MustCompile("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}")

type StringSlice struct {
	Data []string
}

func (s *StringSlice) Scan(src any) error {
	arr, ok := src.([]interface{})
	if !ok {
		return fmt.Errorf("unsupported value")
	}

	if len(arr) == 0 {
		return nil
	}

	for _, elem := range arr {
		str, ok := elem.(string)
		if !ok {
			return fmt.Errorf("unsupported value in array")
		}

		s.Data = append(s.Data, str)
	}

	return nil
}

type IntSlice struct {
	Data []int
}

func (s *IntSlice) Scan(src any) error {
	arr, ok := src.([]interface{})
	if !ok {
		return fmt.Errorf("unsupported value")
	}

	if len(arr) == 0 {
		return nil
	}

	for _, elem := range arr {
		switch num := elem.(type) {
		case int:
			s.Data = append(s.Data, num)
		case int64:
			s.Data = append(s.Data, int(num))
		case int32:
			s.Data = append(s.Data, int(num))
		case int16:
			s.Data = append(s.Data, int(num))
		case int8:
			s.Data = append(s.Data, int(num))
		case uint:
			s.Data = append(s.Data, int(num))
		case uint64:
			s.Data = append(s.Data, int(num))
		case uint32:
			s.Data = append(s.Data, int(num))
		case uint16:
			s.Data = append(s.Data, int(num))
		case uint8:
			s.Data = append(s.Data, int(num))
		default:
			return fmt.Errorf("unsupported value in array")
		}
	}

	return nil
}

type FloatSlice struct {
	Data []float64
}

func (s *FloatSlice) Scan(src any) error {
	arr, ok := src.([]interface{})
	if !ok {
		return fmt.Errorf("unsupported value")
	}

	if len(arr) == 0 {
		return nil
	}

	for _, elem := range arr {
		switch num := elem.(type) {
		case int:
			s.Data = append(s.Data, float64(num))
		case int64:
			s.Data = append(s.Data, float64(num))
		case int32:
			s.Data = append(s.Data, float64(num))
		case int16:
			s.Data = append(s.Data, float64(num))
		case int8:
			s.Data = append(s.Data, float64(num))
		case uint:
			s.Data = append(s.Data, float64(num))
		case uint64:
			s.Data = append(s.Data, float64(num))
		case uint32:
			s.Data = append(s.Data, float64(num))
		case uint16:
			s.Data = append(s.Data, float64(num))
		case uint8:
			s.Data = append(s.Data, float64(num))
		case float32:
			s.Data = append(s.Data, float64(num))
		case float64:
			s.Data = append(s.Data, num)
		default:
			return fmt.Errorf("unsupported value in array")
		}
	}

	return nil
}

type Many[T any] []T

func (my *Many[T]) Scan(value any) error {
	if value == nil {
		*my = Many[T]{}
		return nil
	}

	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, my)
}

// Surreal returns times as strings
type SurrealTime struct {
	time.Time
}

func (s *SurrealTime) Scan(value any) error {
	// TODO: check other time formats
	val, err := time.Parse(time.RFC3339, value.(string))
	*s = SurrealTime{val}

	return err
}

func (s *SurrealTime) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// This supports extracting the UUID from a surrealdb id that could
// come back in the format of user:uuid or user:`uuid` or user<uuid>
type SurrealUUID uuid.UUID

func (s *SurrealUUID) Scan(value any) error {
	uuidSTR := uuidRE.FindString(value.(string))
	uid, err := uuid.Parse(uuidSTR)
	*s = SurrealUUID(uid)

	return err
}

func (s *SurrealUUID) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *SurrealUUID) UUID() uuid.UUID {
	return uuid.UUID(*s)
}

type SurrealAutoID string

func (s *SurrealAutoID) Scan(value any) error {
	parts := strings.Split(value.(string), ":")
	if len(parts) < 1 {
		return fmt.Errorf(`%v is not a valid SurrealAutoID`, value)
	}

	*s = SurrealAutoID(parts[1])

	return nil
}

func (s *SurrealAutoID) Value() (driver.Value, error) {
	return json.Marshal(s)
}
