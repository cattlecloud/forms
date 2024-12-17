// Copyright (c) CattleCloud LLC
// SPDX-License-Identifier: BSD-3-Clause

// Package formdata provides a way to safely and conveniently extract html Form
// data using a definied schema.
package formdata

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/shoenig/go-conceal"
)

var (
	ErrNoValue         = errors.New("expected value to exist")
	ErrMulitpleValues  = errors.New("expected only one value to exist")
	ErrFieldNotPresent = errors.New("requested field does not exist")
	ErrParseFailure    = errors.New("could not parse value")
)

func Parse(data url.Values, schema Schema) error {
	for name, parser := range schema {
		values := data[name]
		if err := parser.Parse(values); err != nil {
			return fmt.Errorf("%s: %w", ErrParseFailure.Error(), err)
		}
	}
	return nil
}

func ParseForm(r *http.Request, schema Schema) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	return Parse(r.Form, schema)
}

// A Schema describes how a set of url.Values should be parsed.
// Typically these are coming from an http.Request.Form from inside an
// http.Handler responding to an inbound request.
type Schema map[string]Parser

// do we care about multi-value? we could provide parsers into slices
// automatically, for example

// A Parser implementation is capable of extracting a value from the value of
// an url.Values, which is a slice of string.
type Parser interface {
	Parse([]string) error
}

// String is used to extract a form data value into a Go string. If the value
// is not a string or is missing then an error is returned during parsing.
func String(s *string) Parser {
	return &stringParser{
		required:    true,
		destination: s,
	}
}

// StringOr is used to extract a form data value into a Go string. If the value
// is missing, then the alt value is used instead.
func StringOr(s *string, alt string) Parser {
	*s = alt
	return &stringParser{
		required:    false,
		destination: s,
	}
}

type stringParser struct {
	required    bool
	destination *string
}

func (p *stringParser) Parse(values []string) error {
	switch {
	case len(values) > 1:
		return ErrMulitpleValues
	case len(values) == 0 && p.required:
		return ErrNoValue
	case len(values) == 0:
		return nil
	default:
		*p.destination = values[0]
	}
	return nil
}

// Secret is used to extract a form data value into a Go conceal.Text. If the
// value is missing then an error is returned during parsing.
func Secret(s **conceal.Text) Parser {
	return &secretParser{
		required:    true,
		destination: s,
	}
}

type secretParser struct {
	required    bool
	destination **conceal.Text
}

func (p *secretParser) Parse(values []string) error {
	switch {
	case len(values) > 1:
		return ErrMulitpleValues
	case len(values) == 0 && p.required:
		return ErrNoValue
	case len(values) == 0:
		return nil
	default:
		text := conceal.New(values[0])
		*p.destination = text
	}
	return nil
}

type intParser struct {
	required    bool
	destination *int
}

// Int is used to extract a form data value into a Go int. If the value is not
// an int or is missing then an error is returned during parsing.
func Int(i *int) Parser {
	return &intParser{
		required:    true,
		destination: i,
	}
}

// IntOr is used to extract a form data value into a Go int. If the value is
// missing, then the alt value is used instead.
func IntOr(i *int, alt int) Parser {
	*i = alt
	return &intParser{
		required:    false,
		destination: i,
	}
}

func (p *intParser) Parse(values []string) error {
	switch {
	case len(values) > 1:
		return ErrMulitpleValues
	case len(values) == 0 && p.required:
		return ErrNoValue
	case len(values) == 0:
		return nil
	}

	i, err := strconv.Atoi(values[0])
	if err != nil {
		return err
	}

	*p.destination = i
	return nil
}

type floatParser struct {
	required    bool
	destination *float64
}

// Float is used to extract a form data value into a Go float64. If the value is
// not a float or is missing then an error is returned during parsing.
func Float(f *float64) Parser {
	return &floatParser{
		required:    true,
		destination: f,
	}
}

// FloatOr is used to extract a form data value into a Go float64. If the value
// is missing, then the alt value is used instead.
func FloatOr(f *float64, alt float64) Parser {
	*f = alt
	return &floatParser{
		required:    false,
		destination: f,
	}
}

func (p *floatParser) Parse(values []string) error {
	switch {
	case len(values) > 1:
		return ErrMulitpleValues
	case len(values) == 0 && p.required:
		return ErrNoValue
	case len(values) == 0:
		return nil
	}

	f, err := strconv.ParseFloat(values[0], 64)
	if err != nil {
		return err
	}

	*p.destination = f
	return nil
}

type boolParser struct {
	required    bool
	destination *bool
}

// Bool is used to extract a form data value into a Go bool. If the value is not
// a bool or is missing than an error is returned during parsing.
func Bool(b *bool) Parser {
	return &boolParser{
		required:    true,
		destination: b,
	}
}

// BoolOr is used to extract a form data value into a Go bool. If the value is
// missing, then the alt value is used instead.
func BoolOr(b *bool, alt bool) Parser {
	*b = alt
	return &boolParser{
		required:    false,
		destination: b,
	}
}

func (p *boolParser) Parse(values []string) error {
	switch {
	case len(values) > 1:
		return ErrMulitpleValues
	case len(values) == 0 && p.required:
		return ErrNoValue
	case len(values) == 0:
		return nil
	}

	b, err := strconv.ParseBool(values[0])
	if err != nil {
		return err
	}

	*p.destination = b
	return nil
}
