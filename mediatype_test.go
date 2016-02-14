package mango

import (
	"sort"
	"testing"
)

func TestNewMediaTypeThrowsNoErrorIfInputIsValid(t *testing.T) {
	want := error(nil)

	s := "text/html"
	_, err := newMediaType(s)

	got := err
	if got != want {
		t.Errorf("Media type = %q, want %q", got, want)
	}
}

func TestNewMediaTypeThrowsErrorIfInputIsInvalid(t *testing.T) {
	want := "invalid media type: \"application\""

	s := "application"
	_, err := newMediaType(s)

	got := err.Error()
	if got != want {
		t.Errorf("Media type = %q, want %q", got, want)
	}
}

func TestNewMediaTypeParsesMainTypeAndSubType(t *testing.T) {
	want := "text/html"

	s := "text/html"
	mt, _ := newMediaType(s)

	got := mt.String()
	if got != want {
		t.Errorf("Media type = %q, want %q", got, want)
	}
}

func TestNewMediaTypeUsesStarSlashStarIfInputIsEmpty(t *testing.T) {
	want := "*/*"

	s := ""
	mt, _ := newMediaType(s)

	got := mt.String()
	if got != want {
		t.Errorf("Media type = %q, want %q", got, want)
	}
}

func TestNewMediaTypeParsesTypeAndSingleParameters(t *testing.T) {
	want := "text/plain;format=flowed"

	s := "text/plain;format=flowed"
	mt, _ := newMediaType(s)

	got := mt.String()
	if got != want {
		t.Errorf("Media type = %q, want %q", got, want)
	}
}

func TestNewMediaTypeParsesTypeAndMultipleParameters(t *testing.T) {
	want := "application/xml;format=flowed;charset=utf-8"

	s := "application/xml;format=flowed;charset=utf-8"
	mt, _ := newMediaType(s)

	got := mt.String()
	if got != want {
		t.Errorf("Media type = %q, want %q", got, want)
	}
}

func TestNewMediaTypeSetsQualityFactorToOneIfNotPresent(t *testing.T) {
	want := float32(1)

	s := "application/xml;format=flowed;charset=utf-8"
	mt, _ := newMediaType(s)

	got := mt.q
	if got != want {
		t.Errorf("Media type = %f, want %f", got, want)
	}
}

func TestNewMediaTypeSetsQualityFactorWhenPresent(t *testing.T) {
	want := float32(0.3)

	s := "application/xml;format=flowed;charset=utf-8;q=0.3"
	mt, _ := newMediaType(s)

	got := mt.q
	if got != want {
		t.Errorf("Media type = %f, want %f", got, want)
	}
}

func TestNewMediaTypeSetsQualityFactorToZeroIfMalformed(t *testing.T) {
	want := float32(0)

	s := "application/xml;format=flowed;charset=utf-8;q=awesome"
	mt, _ := newMediaType(s)

	got := mt.q
	if got != want {
		t.Errorf("Media type = %f, want %f", got, want)
	}
}

func TestNewMediaTypeSetsQualityFactorToZeroIfEmpty(t *testing.T) {
	want := float32(0)

	s := "application/xml;format=flowed;charset=utf-8;q="
	mt, _ := newMediaType(s)

	got := mt.q
	if got != want {
		t.Errorf("Media type = %f, want %f", got, want)
	}
}

func TestNewMediaTypeExcludesQualityFactorFromStringRepresentation(t *testing.T) {
	want := "application/xml;format=flowed;charset=utf-8"

	s := "application/xml;format=flowed;charset=utf-8;q=0.3"
	mt, _ := newMediaType(s)

	got := mt.String()
	if got != want {
		t.Errorf("Media type = %q, want %q", got, want)
	}
}

func TestNewMediaTypeIgnoresParametersAfterQualityFactor(t *testing.T) {
	want := "application/xml;format=flowed;charset=utf-8"

	s := "application/xml;format=flowed;charset=utf-8;q=0.3;custard=9"
	mt, _ := newMediaType(s)

	got := mt.String()
	if got != want {
		t.Errorf("Media type = %q, want %q", got, want)
	}
}

func TestMediaTypesSortsSpecificMainTypeOverAny(t *testing.T) {
	want := "*/*"

	mts := make(mediaTypes, 3)
	s := "text/xml"
	mt, _ := newMediaType(s)
	mts[0] = *mt
	s = "*/*"
	mt, _ = newMediaType(s)
	mts[1] = *mt
	s = "application/json"
	mt, _ = newMediaType(s)
	mts[2] = *mt
	sort.Sort(mts)

	got := mts[2].String()
	if got != want {
		t.Errorf("Media type = %q, want %q", got, want)
	}
}

func TestMediaTypesSortsSpecificSubtypeOverAnyWhenMainTypeIsSame(t *testing.T) {
	want := "text/*"

	mts := make(mediaTypes, 3)
	s := "text/xml"
	mt, _ := newMediaType(s)
	mts[0] = *mt
	s = "text/*"
	mt, _ = newMediaType(s)
	mts[1] = *mt
	s = "text/html"
	mt, _ = newMediaType(s)
	mts[2] = *mt
	sort.Sort(mts)

	got := mts[2].String()
	if got != want {
		t.Errorf("Media type = %q, want %q", got, want)
	}
}

func TestMediaTypesSortsOnParametersWhenMainAndSubtypeAreSame(t *testing.T) {
	want := "text/xml;charset=utf-8"

	mts := make(mediaTypes, 3)
	s := "text/html"
	mt, _ := newMediaType(s)
	mts[0] = *mt
	s = "text/xml"
	mt, _ = newMediaType(s)
	mts[1] = *mt
	s = "text/xml;charset=utf-8"
	mt, _ = newMediaType(s)
	mts[2] = *mt
	sort.Sort(mts)

	got := mts[0].String()
	if got != want {
		t.Errorf("Media type = %q, want %q", got, want)
	}
}

func TestMediaTypesSortsOnParameterLengthWhenMainAndSubtypeAreSame(t *testing.T) {
	want := "text/xml;charset=utf-16"

	mts := make(mediaTypes, 3)
	s := "text/xml;format=flowed"
	mt, _ := newMediaType(s)
	mts[0] = *mt
	s = "text/xml"
	mt, _ = newMediaType(s)
	mts[1] = *mt
	s = "text/xml;charset=utf-16"
	mt, _ = newMediaType(s)
	mts[2] = *mt
	sort.Sort(mts)

	got := mts[0].String()
	if got != want {
		t.Errorf("Media type = %q, want %q", got, want)
	}
}

func TestMediaTypesSortsOnQualityFactorsWhenPresent(t *testing.T) {
	want := "text/xml"

	mts := make(mediaTypes, 3)
	s := "text/html;q=0.7"
	mt, _ := newMediaType(s)
	mts[0] = *mt
	s = "text/xml;q=0.9"
	mt, _ = newMediaType(s)
	mts[1] = *mt
	s = "text/xml;charset=utf-8;q=0.4"
	mt, _ = newMediaType(s)
	mts[2] = *mt
	sort.Sort(mts)

	got := mts[0].String()
	if got != want {
		t.Errorf("Media type = %q, want %q", got, want)
	}
}
