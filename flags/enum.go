package flags

import "fmt"

type enum struct {
	target  *string
	options []string
}

func NewEnum(target *string, options ...string) *enum {
	return &enum{target: target, options: options}
}

func (f *enum) String() string {
	return *f.target
}

func (f *enum) Set(value string) error {
	for _, v := range f.options {
		if v == value {
			*f.target = value
			return nil
		}
	}

	return fmt.Errorf("expected one of the following %q", f.options)
}
