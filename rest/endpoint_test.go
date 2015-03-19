package rest

import (
	"sort"
	"testing"
)

func TestSortWithStaticUrl(t *testing.T) {
	input := []*Endpoint{
		&Endpoint{"/", "GET", nil},
		&Endpoint{"/user/sample", "GET", nil},
		&Endpoint{"/user", "GET", nil},
	}

	output := []*Endpoint{
		&Endpoint{"/user/sample", "GET", nil},
		&Endpoint{"/user", "GET", nil},
		&Endpoint{"/", "GET", nil},
	}
	sort.Sort(Endpoints(input))
	compare(input, output, t)

}

func TestSortWithStaticUrlDifferentMethods(t *testing.T) {

	input := []*Endpoint{
		&Endpoint{"/user/sample", "GET", nil},
		&Endpoint{"/user/sample", "POST", nil},
	}

	output := []*Endpoint{
		&Endpoint{"/user/sample", "POST", nil},
		&Endpoint{"/user/sample", "GET", nil},
	}

	sort.Sort(Endpoints(input))
	compare(input, output, t)
}

func TestSortWithSameStaticUrl(t *testing.T) {

	input := []*Endpoint{
		&Endpoint{"/user/sample", "GET", nil},
		&Endpoint{"/user/sample", "GET", nil},
	}

	defer func() {
		if error := recover(); error == nil {
			t.Errorf("Expected a panic")
		}
	}()
	sort.Sort(Endpoints(input))
}

// Dynamic URL will be pushed to last
func TestSortWithDynamicUrl(t *testing.T) {

	input := []*Endpoint{
		&Endpoint{"/user/{sample}/hello", "GET", nil},
		&Endpoint{"/", "GET", nil},
		&Endpoint{"/{app}", "GET", nil},
		&Endpoint{"/user/{sample}", "GET", nil},
		&Endpoint{"/user", "GET", nil},
	}

	output := []*Endpoint{
		&Endpoint{"/user", "GET", nil},
		&Endpoint{"/", "GET", nil},
		&Endpoint{"/user/{sample}/hello", "GET", nil},
		&Endpoint{"/user/{sample}", "GET", nil},
		&Endpoint{"/{app}", "GET", nil},
	}

	sort.Sort(Endpoints(input))
	compare(input, output, t)
}

func compare(a, b Endpoints, t *testing.T) bool {
	for i, v := range a {
		if v.URL != b[i].URL || v.Method != b[i].Method {
			t.Errorf("Input and output are not matching input :\n%v\noutput:\n%v", a, b)
			return false
		}
	}
	return true
}
