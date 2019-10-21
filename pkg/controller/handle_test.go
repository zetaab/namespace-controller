package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type findTest struct {
	Name    string
	Correct string
	Path    string
	NSName  string
	Pass    bool
	Pass2   bool
}

func TestFindTeam(t *testing.T) {
	tests := []findTest{
		{
			Name:    "file not found",
			Path:    "/foo/bar",
			Correct: "",
			NSName:  "foo",
			Pass:    false,
			Pass2:   false,
		},
		{
			Name:    "team not found",
			Path:    "./testdata/data2.yaml",
			Correct: "",
			NSName:  "test",
			Pass:    true,
			Pass2:   false,
		},
		{
			Name:    "team kaas",
			Path:    "./testdata/data1.yaml",
			Correct: "kaas",
			NSName:  "foo",
			Pass:    true,
			Pass2:   true,
		},
		{
			Name:    "team prod",
			Path:    "./testdata/data1.yaml",
			Correct: "prod1",
			NSName:  "testing-prod",
			Pass:    true,
			Pass2:   true,
		},
		{
			Name:    "team dev",
			Path:    "./testdata/data1.yaml",
			Correct: "dev1",
			NSName:  "others",
			Pass:    true,
			Pass2:   true,
		},
		{
			Name:    "team dev2",
			Path:    "./testdata/data1.yaml",
			Correct: "dev1",
			NSName:  "testing-proding",
			Pass:    true,
			Pass2:   true,
		},
	}
	for _, test := range tests {
		config, err := makeConfig(test.Path)
		assert.Equal(t, test.Pass, err == nil, "Failed test '%s' expected %v got %v (%v)", test.Name, test.Pass, err == nil, err)
		if config == nil {
			continue
		}
		team, err := FindTeam(test.NSName, config)
		assert.Equal(t, test.Pass2, err == nil, "Failed test '%s' expected %v got %v (%v)", test.Name, test.Pass2, err == nil, err)
		assert.Equal(t, test.Correct, team, "Failed test '%s' expected %v got %v", test.Name, test.Correct, team)
	}
}
