// Package database contains our database connection. This file hosts all the basics for
// creating a connection, wihle other files have the query functions.
package database

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/vinovest/sqlx"
)

type DB struct {
	conn *sqlx.DB
}

func New(databaseURI string) (*DB, error) {
	dsn, err := URLtoDSN(databaseURI)
	if err != nil {
		return nil, err
	}

	conn, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return &DB{
		conn,
	}, nil
}

func URLtoDSN(databaseURL string) (string, error) {
	u, err := url.Parse(databaseURL)
	if err != nil {
		return "", err
	}
	user := u.User.Username()
	pass, _ := u.User.Password()
	host := u.Host
	dbname := strings.TrimPrefix(u.Path, "/")

	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", user, pass, host, dbname), nil
}

// BinaryUUID type, neeed for scanning in and pushing out from the database, taking into account that
// the existing Athena and ModMail databases store UUIDs as varchars, while we use BINARY(16).
type BinaryUUID uuid.UUID

// Value converts to BINARY(16) when writing to DB
func (b BinaryUUID) Value() (driver.Value, error) {
	return uuid.UUID(b).MarshalBinary()
}

// Scan converts from BINARY(16) when reading from DB
func (b *BinaryUUID) Scan(src any) error {
	switch v := src.(type) {
	case []byte:
		if len(v) == 16 {
			// Already binary
			parsed, err := uuid.FromBytes(v)
			if err != nil {
				return fmt.Errorf("BinaryUUID: %w", err)
			}
			*b = BinaryUUID(parsed)
		} else {
			// VARCHAR coming back as []byte
			parsed, err := uuid.ParseBytes(v)
			if err != nil {
				return fmt.Errorf("BinaryUUID: %w", err)
			}
			*b = BinaryUUID(parsed)
		}
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			return fmt.Errorf("BinaryUUID: %w", err)
		}
		*b = BinaryUUID(parsed)
	default:
		return fmt.Errorf("BinaryUUID: expected []byte or string, got %T", src)
	}
	return nil
}

func (b BinaryUUID) String() string {
	return uuid.UUID(b).String()
}

func (b BinaryUUID) MarshalJSON() ([]byte, error) {
	return json.Marshal(uuid.UUID(b).String())
}

func (b *BinaryUUID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	parsed, err := uuid.Parse(s)
	if err != nil {
		return fmt.Errorf("BinaryUUID: %w", err)
	}
	*b = BinaryUUID(parsed)
	return nil
}

// CommaSeparated is just a string slice that is internally represented as String,String,String
type CommaSeparated []string

func (c CommaSeparated) Value() (driver.Value, error) {
	return strings.Join(c, ","), nil
}

func (c *CommaSeparated) Scan(src any) error {
	switch v := src.(type) {
	case []byte:
		*c = strings.Split(string(v), ",")
	case string:
		*c = strings.Split(v, ",")
	case nil:
		*c = nil
	default:
		return fmt.Errorf("CommaSeparated: unsupported type %T", src)
	}
	return nil
}

func (c CommaSeparated) MarshalJSON() ([]byte, error) {
	return json.Marshal([]string(c))
}

func (c *CommaSeparated) UnmarshalJSON(data []byte) error {
	var s []string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*c = s
	return nil
}

// JSONStringArray scans a JSON-encoded TEXT column into []string
type JSONStringArray []string

func (j JSONStringArray) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal([]string(j))
}

func (j *JSONStringArray) Scan(src any) error {
	var b []byte
	switch v := src.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	case nil:
		*j = nil
		return nil
	default:
		return fmt.Errorf("JSONStringArray: unsupported type %T", src)
	}
	return json.Unmarshal(b, j)
}

func (j JSONStringArray) MarshalJSON() ([]byte, error) {
	return json.Marshal([]string(j))
}

func (j *JSONStringArray) UnmarshalJSON(data []byte) error {
	var s []string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*j = s
	return nil
}

// JSONMap scans a JSON-encoded TEXT column into map[string]any
type JSONMap map[string]any

func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(map[string]any(j))
}

func (j *JSONMap) Scan(src any) error {
	var b []byte
	switch v := src.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	case nil:
		*j = nil
		return nil
	default:
		return fmt.Errorf("JSONMap: unsupported type %T", src)
	}
	return json.Unmarshal(b, (*map[string]any)(j))
}

func (j JSONMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any(j))
}

func (j *JSONMap) UnmarshalJSON(data []byte) error {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	*j = m
	return nil
}
