package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type findTest struct {
	Name   string
	Labels map[string]string
	Path   string
	NSName string
	Pass   bool
	Pass2  bool
}

func TestFindTeam(t *testing.T) {
	tests := []findTest{
		{
			Name:   "file not found",
			Path:   "/foo/bar",
			Labels: map[string]string{},
			NSName: "foo",
			Pass:   false,
			Pass2:  false,
		},
		{
			Name:   "team not found",
			Path:   "./testdata/data2.yaml",
			Labels: map[string]string{},
			NSName: "test",
			Pass:   true,
			Pass2:  false,
		},
		{
			Name:   "team kaas",
			Path:   "./testdata/data1.yaml",
			Labels: map[string]string{"maintainer": "kaas", "foo": "bar"},
			NSName: "foo",
			Pass:   true,
			Pass2:  true,
		},
		{
			Name:   "team prod",
			Path:   "./testdata/data1.yaml",
			Labels: map[string]string{"maintainer": "prod1"},
			NSName: "testing-prod",
			Pass:   true,
			Pass2:  true,
		},
		{
			Name:   "team dev",
			Path:   "./testdata/data1.yaml",
			Labels: map[string]string{"maintainer": "dev1"},
			NSName: "others",
			Pass:   true,
			Pass2:  true,
		},
		{
			Name:   "team dev2",
			Path:   "./testdata/data1.yaml",
			Labels: map[string]string{"maintainer": "dev1"},
			NSName: "testing-proding",
			Pass:   true,
			Pass2:  true,
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
		if team != nil {
			assert.Equal(t, test.Labels, team.Labels, "Failed test '%s'", test.Name)
		}
	}
}

type handleTest struct {
	Name   string
	NS     *v1.Namespace
	Key    string
	Opt    string
	Result string
}

func TestHandle(t *testing.T) {
	tests := []handleTest{
		{
			Name: "add new label",
			NS: &v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name:   "test",
					Labels: map[string]string{},
				},
			},
			Key:    "foo",
			Opt:    "bar",
			Result: "bar",
		},
		{
			Name: "existing label",
			NS: &v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					Labels: map[string]string{
						"foo": "bar",
					},
				},
			},
			Key:    "foo",
			Opt:    "bar",
			Result: "",
		},
		{
			Name: "not allowed value",
			NS: &v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					Labels: map[string]string{
						"foo": "bar2",
					},
				},
			},
			Key:    "foo",
			Opt:    "bar",
			Result: "bar",
		},
		{
			Name: "AllowedValues valid",
			NS: &v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					Labels: map[string]string{
						"foo": "bar2",
					},
				},
			},
			Key:    "foo",
			Opt:    "bar,bar2",
			Result: "",
		},
		{
			Name: "AllowedValues not valid",
			NS: &v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					Labels: map[string]string{
						"foo": "bar3",
					},
				},
			},
			Key:    "foo",
			Opt:    "bar,bar2",
			Result: "bar",
		},

		{
			Name: "AllowedValues not valid",
			NS: &v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					Labels: map[string]string{
						"pod-security.kubernetes.io/enforce": "restricted",
					},
				},
			},
			Key:    "pod-security.kubernetes.io/warn",
			Opt:    "restricted,baseline",
			Result: "",
		},
	}
	for _, test := range tests {
		result := handleLabels(test.NS, test.Key, test.Opt)
		assert.Equal(t, test.Result, result, "Failed test '%s'", test.Name)
	}
}
