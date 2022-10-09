package models

import (
	"database/sql"
	"encoding/json"
)

type NullInt32 struct {
	sql.NullInt32
}

func (v NullInt32) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int32)
	}

	return []byte(`null`), nil
}
