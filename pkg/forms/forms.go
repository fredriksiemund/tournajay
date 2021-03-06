package forms

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

// Create a custom Form struct, which anonymously embeds a url.Values object
// (to hold the form data) and an Errors field to hold any validation errors
// for the form data.
type Form struct {
	url.Values
	Errors errors
}

func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("This field is too long (maximum is %d characters)", d))
	}
}

func (f *Form) PermittedValues(field string, opts ...string) {
	value := f.Get(field)
	if value == "" {
		return
	}
	for _, opt := range opts {
		if value == opt {
			return
		}
	}
	f.Errors.Add(field, fmt.Sprintf("This field has to be one of: %v", opts))
}

func (f *Form) MinLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) < d {
		f.Errors.Add(field, "This field is invalid")
	}
}

func (f *Form) MatchesPattern(field string, pattern *regexp.Regexp) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if !pattern.MatchString(value) {
		f.Errors.Add(field, "This field is invalid")
	}
}

func (f *Form) ValidDate(field, layout string) {
	value := f.Get(field)
	if value == "" {
		return
	}
	d, err := time.Parse(layout, value)
	if err != nil {
		f.Errors.Add(field, "This field is invalid")
	}
	if d.Before(time.Now()) {
		f.Errors.Add(field, "This field should be a date in the future")
	}
	f.Set(field, d.Format("2006-01-02 15:04"))
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
