package egov

import (
	"net/url"
	"strconv"
	"strings"
)

// queryBuilder is an internal helper that omits zero values when building a
// query string. The e-Gov API treats missing parameters and explicit empty
// values differently for some fields, so the builder distinguishes:
//
//   - empty strings:    skipped
//   - zero ints:        skipped (no field uses 0 as a meaningful value)
//   - nil *bool:        skipped; *bool=false is sent as "false"
//   - empty []string:   skipped; multi-value fields are joined with ","
type queryBuilder struct {
	values url.Values
}

func newQuery() *queryBuilder {
	return &queryBuilder{values: url.Values{}}
}

func (q *queryBuilder) set(key, v string) {
	if v != "" {
		q.values.Set(key, v)
	}
}

func (q *queryBuilder) setInt(key string, v int) {
	if v != 0 {
		q.values.Set(key, strconv.Itoa(v))
	}
}

func (q *queryBuilder) setBoolPtr(key string, v *bool) {
	if v != nil {
		q.values.Set(key, strconv.FormatBool(*v))
	}
}

func (q *queryBuilder) setStringSlice(key string, v []string) {
	if len(v) == 0 {
		return
	}
	q.values.Set(key, strings.Join(v, ","))
}

func (q *queryBuilder) encode() string {
	return q.values.Encode()
}
