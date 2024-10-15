package utils

import (
	"bytes"
	"encoding/xml"
	"golang.org/x/net/html/charset"
)

func XMLDecoder(v interface{}, xmlData []byte) error {
	decoder := xml.NewDecoder(bytes.NewReader(xmlData))
	decoder.CharsetReader = charset.NewReaderLabel
	err := decoder.Decode(v)

	return err
}
