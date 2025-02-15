package makemkv

import (
	"fmt"
	"strconv"
	"time"

	"github.com/curt-hash/mkvbot/pkg/makemkv/defs"
)

// Info is a slice of related Attributes.
type Info []*Attribute

// GetCode returns the Code of the Attribute where ID matches id or ErrNotFound
// if such an attribute does not exist.
func (info Info) GetCode(id defs.Attr) (int, error) {
	for _, attr := range info {
		if attr.ID == int(id) {
			return attr.Code, nil
		}
	}

	return 0, ErrNotFound
}

// GetCodeDefault returns the Code of the Attribute where ID matches id or
// defaultValue if such an attribute does not exist.
func (info Info) GetCodeDefault(id defs.Attr, defaultValue int) int {
	v, err := info.GetCode(id)
	if err != nil {
		return defaultValue
	}

	return v
}

// GetAttr returns the Value of the Attribute where ID matches id or
// ErrNotFound if such an attribute does not exist.
func (info Info) GetAttr(id defs.Attr) (string, error) {
	for _, attr := range info {
		if attr.ID == int(id) {
			return string(attr.Value), nil
		}
	}

	return "", ErrNotFound
}

// GetAttrDefault returns the Value of the Attribute where ID matches id or
// defaultValue if such an attribute does not exist.
func (info Info) GetAttrDefault(id defs.Attr, defaultValue string) string {
	v, err := info.GetAttr(id)
	if err != nil {
		return defaultValue
	}

	return v
}

// GetAttrInto is like GetAttr, except it also attempts to convert the Value to
// an integer.
func (info Info) GetAttrInt(id defs.Attr) (int, error) {
	v, err := info.GetAttr(id)
	if err != nil {
		return 0, err
	}

	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("parse %q: %w", v, err)
	}

	return n, nil
}

// GetAttrDuration is like GetAttr, except it also attempts to convert the
// value to a time.Duration.
func (info Info) GetAttrDuration(id defs.Attr) (time.Duration, error) {
	v, err := info.GetAttr(id)
	if err != nil {
		return 0, err
	}

	d, err := ParseDuration(v)
	if err != nil {
		return 0, fmt.Errorf("parse %q: %w", v, err)
	}

	return d, nil
}
