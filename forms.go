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

// Parse uses the given Schema to parse the HTTP form values in the given HTTP
// Request. If the values of the form do not match the schema, or required values
// are missing, an error is returned.
func Parse(r *http.Request, schema Schema) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	return ParseValues(r.Form, schema)
}

// ParseValues uses the given Schema to parse the values in the given url.Values.
// If the values do not match the schema, or required values are missing, an
// error is returned.
func ParseValues(data url.Values, schema Schema) error {
	for name, parser := range schema {
		values := data[name]
		if err := parser.Parse(values); err != nil {
			return fmt.Errorf("%s: %w", ErrParseFailure.Error(), err)
		}
	}
	return nil
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

// StringType represents any type compatible with the Go string built-in type,
// to be used as a destination for writing the value of an environment variable.
type StringType interface {
	~string
}

type stringParser[T StringType] struct {
	required    bool
	destination *T
}

func (p *stringParser[T]) Parse(values []string) error {
	switch {
	case len(values) > 1:
		return ErrMulitpleValues
	case len(values) == 0 && p.required:
		return ErrNoValue
	case len(values) == 0:
		return nil
	default:
		*p.destination = T(values[0])
	}
	return nil
}

// String is used to extract a form data value into a Go string. If the value
// is not a string or is missing then an error is returned during parsing.
func String[T StringType](s *T) Parser {
	return &stringParser[T]{
		required:    true,
		destination: s,
	}
}

// StringOr is used to extract a form data value into a Go string. If the value
// is missing, then the alt value is used instead.
func StringOr[T StringType](s *T, alt T) Parser {
	*s = alt
	return &stringParser[T]{
		required:    false,
		destination: s,
	}
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

// IntType represents any type compatible with the Go integer built-in types,
// to be used as a destination for writing the value of a form value.
type IntType interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type intParser[T IntType] struct {
	required    bool
	destination *T
}

func (p *intParser[T]) Parse(values []string) error {
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

	*p.destination = T(i)
	return nil
}

// Int is used to extract a form data value into a Go int. If the value is not
// an int or is missing then an error is returned during parsing.
func Int[T IntType](i *T) Parser {
	return &intParser[T]{
		required:    true,
		destination: i,
	}
}

// IntOr is used to extract a form data value into a Go int. If the value is
// missing, then the alt value is used instead.
func IntOr[T IntType](i *T, alt T) Parser {
	*i = alt
	return &intParser[T]{
		required:    false,
		destination: i,
	}
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
