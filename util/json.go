package util

import (
	"encoding/json"
	"io"
)

// FromJSON deserializes the object from JSON string
// in an io.Reader to the given interface
func FromJSON(i interface{}, r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(i)
}

// func (p *Product) FromJSON(r io.Reader) error {
// 	e := json.NewDecoder(r)
// 	return e.Decode(p)
// }

// ToJSON serializes the given interface into a string based JSON format
func ToJSON(i interface{}, w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(i)
}

// func (ps *Products) ToJSON(w io.Writer) error {
// 	e := json.NewEncoder(w)
// 	return e.Encode(ps)
// }
