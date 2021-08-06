package main

import (
	"testing"
)

type Profile struct {
	Age  int
	City string
}

func TestExtractStrings(t *testing.T) {
	cases := []struct {
		Name     string
		Input    interface{}
		Expected []string
	}{
		{
			"Array,Slice,Map",
			struct {
				Name     string
				City     string
				Age      int32
				tag      string
				Hobby    []string
				Label    [3]string
				Features map[string]interface{}
			}{
				"None",
				"Guangzhou",
				23,
				"15",
				[]string{"Hobby-A", "Hobby-B", "Hobby-C"},
				[3]string{"Label-A", "Label-B", "Label-C"},
				map[string]interface{}{"F-A": 1, "F-B": "Test", "F-C": [3]string{"FLabel-A", "FLabel-B", "FLabel-C"}},
			},
			[]string{"None", "Guangzhou", "Hobby-A", "Hobby-B", "Hobby-C", "Label-A", "Label-B", "Label-C", "Test", "FLabel-A", "FLabel-B", "FLabel-C"},
		},
		{
			"NestedStruct",
			struct {
				Welcome  string
				profiles []Profile
			}{
				"Hello,World",
				[]Profile{{15, "Gz"}, {20, "Hz"}},
			},
			[]string{"Hello,World", "Gz", "Hz"},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			var got []string
			ExtractStrings(c.Input, func(input string) {
				got = append(got, input)
			})

			for _, g := range got {
				assertContains(t, c.Expected, g)
			}
		})
	}
}

func assertContains(t *testing.T, list []string, expected string) {
	for _, item := range list {
		if item == expected {
			return
		}
	}

	t.Errorf("expected %+v to contains %s, but it didnt", list, expected)
}
