package models

type Permissions struct {
	Read  bool
	Write bool
	Admin bool
	Use   bool
}

// Unmarshaler ignores the field value completely.
func (*Permissions) UnmarshalJSON(data []byte) error {
	return nil
}
