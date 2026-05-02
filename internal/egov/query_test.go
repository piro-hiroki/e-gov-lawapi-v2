package egov

import "testing"

func TestQueryBuilder_SkipsZeroValues(t *testing.T) {
	q := newQuery()
	q.set("a", "")
	q.setInt("b", 0)
	q.setStringSlice("c", nil)
	q.setStringSlice("d", []string{})
	q.setBoolPtr("e", nil)
	if got := q.encode(); got != "" {
		t.Errorf("expected empty query string, got %q", got)
	}
}

func TestQueryBuilder_BoolPtrFalseIsSent(t *testing.T) {
	// nil-vs-false distinction: nil = "user didn't set it", *false = "explicitly false".
	// The API treats these differently for boolean flags like remain_in_force,
	// so a *bool=false MUST appear in the query string.
	f := false
	q := newQuery()
	q.setBoolPtr("flag", &f)
	if got := q.encode(); got != "flag=false" {
		t.Errorf("expected flag=false, got %q", got)
	}
}

func TestQueryBuilder_BoolPtrTrue(t *testing.T) {
	tr := true
	q := newQuery()
	q.setBoolPtr("flag", &tr)
	if got := q.encode(); got != "flag=true" {
		t.Errorf("expected flag=true, got %q", got)
	}
}

func TestQueryBuilder_StringSliceJoinsWithComma(t *testing.T) {
	q := newQuery()
	q.setStringSlice("law_type", []string{"Act", "Rule"})
	// url.Values escapes the comma to %2C.
	if got := q.encode(); got != "law_type=Act%2CRule" {
		t.Errorf("expected law_type=Act%%2CRule, got %q", got)
	}
}

func TestQueryBuilder_IntZeroSkipped(t *testing.T) {
	q := newQuery()
	q.setInt("limit", 0)
	q.setInt("offset", 10)
	if got := q.encode(); got != "offset=10" {
		t.Errorf("expected offset=10 only, got %q", got)
	}
}

func TestQueryBuilder_OverwritesOnRepeatedSet(t *testing.T) {
	// Same key set twice keeps only the last value (set semantics, not append).
	q := newQuery()
	q.set("k", "first")
	q.set("k", "second")
	if got := q.encode(); got != "k=second" {
		t.Errorf("expected k=second, got %q", got)
	}
}

func TestQueryBuilder_EncodesSpecialCharacters(t *testing.T) {
	// Japanese text must be percent-encoded in the query string.
	q := newQuery()
	q.set("law_title", "個人情報保護法")
	got := q.encode()
	want := "law_title=%E5%80%8B%E4%BA%BA%E6%83%85%E5%A0%B1%E4%BF%9D%E8%AD%B7%E6%B3%95"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
