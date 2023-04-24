package nilType

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type NilType interface {
	IsNotNil() bool
	InterfaceValue() interface{}
	SetValue(val interface{}) (err error)

	fmt.Stringer

	sql.Scanner
	driver.Valuer

	json.Marshaler
	json.Unmarshaler
}
