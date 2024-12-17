// Copyright (c) CattleCloud LLC
// SPDX-License-Identifier: BSD-3-Clause

package formdata

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/shoenig/go-conceal"
	"github.com/shoenig/test/must"
)

func Test_Parse_singles(t *testing.T) {
	t.Parallel()

	data := url.Values{
		"one":   []string{"1"},
		"two":   []string{"2"},
		"three": []string{"3.1"},
		"four":  []string{"true"},
		"five":  []string{"xyz"},
	}

	var (
		one   string
		two   int
		three float64
		four  bool
		five  *conceal.Text
	)

	err := ParseValues(data, Schema{
		"one":   String(&one),
		"two":   Int(&two),
		"three": Float(&three),
		"four":  Bool(&four),
		"five":  Secret(&five),
	})
	must.NoError(t, err)
	must.Eq(t, "1", one)
	must.Eq(t, 2, two)
	must.Eq(t, 3.1, three)
	must.True(t, four)
	must.Eq(t, "xyz", five.Unveil())
}

func Test_Parse_singles_Or(t *testing.T) {
	t.Parallel()

	data := url.Values{
		"string1": []string{"hi"},
		"string2": nil,
		"int1":    []string{"1"},
		"int2":    nil,
		"float1":  []string{"2.2"},
		"float2":  nil,
		"bool1":   []string{"true"},
		"bool2":   nil,
	}

	var (
		s1, s2 string
		i1, i2 int
		f1, f2 float64
		b1, b2 bool
	)

	err := ParseValues(data, Schema{
		"string1": StringOr(&s1, "X"),
		"string2": StringOr(&s2, "X"),
		"int1":    IntOr(&i1, 3),
		"int2":    IntOr(&i2, 4),
		"float1":  FloatOr(&f1, 5.5),
		"float2":  FloatOr(&f2, 6.6),
		"bool1":   BoolOr(&b1, false),
		"bool2":   BoolOr(&b2, true),
	})

	must.NoError(t, err)
	must.Eq(t, "hi", s1)
	must.Eq(t, "X", s2)
	must.Eq(t, 1, i1)
	must.Eq(t, 4, i2)
	must.Eq(t, 2.2, f1)
	must.Eq(t, 6.6, f2)
	must.Eq(t, true, b1)
	must.Eq(t, true, b2)
}

func Test_Parse_HTMLForm(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", nil)
	must.NoError(t, err)

	request.PostForm = make(url.Values)
	request.PostForm.Set("one", "1")
	request.PostForm.Set("two", "2")
	request.PostForm.Set("three", "3.1")
	request.PostForm.Set("four", "true")

	var (
		one   string
		two   int
		three float64
		four  bool
	)

	err2 := Parse(request, Schema{
		"one":   String(&one),
		"two":   Int(&two),
		"three": Float(&three),
		"four":  Bool(&four),
	})
	must.NoError(t, err2)
	must.Eq(t, "1", one)
	must.Eq(t, 2, two)
	must.Eq(t, 3.1, three)
	must.True(t, four)
}

func Test_Parse_HTMLForm_optional(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", nil)
	must.NoError(t, err)

	request.PostForm = make(url.Values)
	request.PostForm.Set("one", "1")

	var (
		one string
		two string
	)

	err2 := Parse(request, Schema{
		"one": String(&one),
		"two": StringOr(&two, "alternate"),
	})
	must.NoError(t, err2)
	must.Eq(t, "1", one)
	must.Eq(t, "alternate", two)
}

func Test_Parse_HTMLForm_not_ready(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, "/", nil)
	must.NoError(t, err)

	var one string

	// not yet a valid form, never had the FormValues field set
	err2 := Parse(request, Schema{
		"one": String(&one),
	})
	must.Error(t, err2)
}

func Test_Parse_key_missing(t *testing.T) {
	t.Parallel()

	data := url.Values{
		"one": []string{"1"},
	}

	var two int
	err := ParseValues(data, Schema{
		"two": Int(&two),
	})
	must.Error(t, err)
}

func Test_Parse_string_value_missing(t *testing.T) {
	t.Parallel()

	data := url.Values{
		"one": []string{},
	}

	var one string
	err := ParseValues(data, Schema{
		"one": String(&one),
	})
	must.Error(t, err)
}

func Test_Parse_int_value_missing(t *testing.T) {
	t.Parallel()

	data := url.Values{
		"two": []string{},
	}

	var two int
	err := ParseValues(data, Schema{
		"two": Int(&two),
	})
	must.Error(t, err)
}

func Test_Parse_int_malformed(t *testing.T) {
	t.Parallel()

	data := url.Values{
		"two": []string{"not an int"},
	}

	var two int
	err := ParseValues(data, Schema{
		"two": Int(&two),
	})
	must.Error(t, err)
}

func Test_Parse_float_value_missing(t *testing.T) {
	t.Parallel()

	data := url.Values{
		"three": []string{},
	}

	var three float64
	err := ParseValues(data, Schema{
		"three": Float(&three),
	})
	must.Error(t, err)
}

func Test_Parse_float_malformed(t *testing.T) {
	t.Parallel()

	data := url.Values{
		"three": []string{"not a float"},
	}

	var three float64
	err := ParseValues(data, Schema{
		"three": Float(&three),
	})
	must.Error(t, err)
}

func Test_Parse_bool_value_missing(t *testing.T) {
	t.Parallel()

	data := url.Values{
		"four": []string{},
	}

	var four bool
	err := ParseValues(data, Schema{
		"four": Bool(&four),
	})
	must.Error(t, err)
}

func Test_Parse_bool_malformed(t *testing.T) {
	t.Parallel()

	data := url.Values{
		"four": []string{"not a bool"},
	}

	var four bool
	err := ParseValues(data, Schema{
		"four": Bool(&four),
	})
	must.Error(t, err)
}

func Test_Parse_StringType_String(t *testing.T) {
	t.Parallel()

	data := url.Values{
		"user": []string{"bob"},
	}

	type username string

	var user username

	err := ParseValues(data, Schema{
		"user": String(&user),
	})
	must.NoError(t, err)
	must.Eq(t, "bob", user)
}

func Test_Parse_StringType_StringOr(t *testing.T) {
	t.Parallel()

	data := url.Values{
		"foo": []string{"bar"},
	}

	type username string

	var user username
	var fallback username = "alice"

	err := ParseValues(data, Schema{
		"user": StringOr(&user, fallback),
	})
	must.NoError(t, err)
	must.Eq(t, "alice", user)
}

func Test_Parse_IntType_Int(t *testing.T) {
	t.Parallel()

	data := url.Values{
		"age": []string{"34"},
	}

	type years int

	var age years

	err := ParseValues(data, Schema{
		"age": Int(&age),
	})
	must.NoError(t, err)
	must.Eq(t, 34, age)
}

func Test_Parse_IntType_IntOr(t *testing.T) {
	t.Parallel()

	data := url.Values{
		"foo": []string{"bar"},
	}

	type years int

	var age years
	var fallback years = 100

	err := ParseValues(data, Schema{
		"age": IntOr(&age, fallback),
	})
	must.NoError(t, err)
	must.Eq(t, 100, age)
}
