package grayt

import (
	"reflect"
	"testing"
)

func TestEndpointAllowed(t *testing.T) {
	for _, str := range []string{
		"/renders",
		"/renders/123",
		"/renders/123/status",
	} {
		if !endpointRE.MatchString(str) {
			t.Errorf("Should have been allowed: %v", str)
		}
	}
}

func TestEndpointNotAllowed(t *testing.T) {
	for _, str := range []string{
		"aoeu/renders", // now allowed prefix
		"/foo",
		"/renders/",              // trailing /
		"/renders/abc",           // uuid must be numeric
		"/renders/123/",          // trailing /
		"/renders/123/two_words", // illegal _
		"/renders//status",       // missing uuid
	} {
		if matches := endpointRE.FindStringSubmatch(str); len(matches) > 0 {
			t.Errorf("Should not have been allowed: %v", str)
		}
	}
}

func TestExtractEndpoint(t *testing.T) {
	for _, test := range []struct {
		url      string
		captured []string
	}{
		{"/renders", []string{"/renders", "", ""}},
		{"/renders/123", []string{"/renders/123", "123", ""}},
		{"/renders/123/status", []string{"/renders/123/status", "123", "status"}},
	} {
		got := endpointRE.FindStringSubmatch(test.url)
		if !reflect.DeepEqual(got, test.captured) {
			t.Errorf("%d %d Want: %v Got: %v", len(test.captured), len(got), test.captured, got)
		}
	}
}
