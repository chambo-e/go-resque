package resque

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewJob(t *testing.T) {
	tests := []struct {
		args  []interface{}
		klass string
	}{
		{
			args:  []interface{}{},
			klass: "",
		},
		{
			args:  []interface{}{"a", 3, true},
			klass: "test",
		},
		{
			args:  []interface{}{"a", "a", "a", 3.14159265},
			klass: "pi",
		},
		{
			args:  []interface{}{"{json: true}"},
			klass: "json",
		},
	}

	for _, test := range tests {
		job := NewJob(test.klass, test.args...)

		require.Equal(t, Job{
			Args:  test.args,
			Class: test.klass,
		}, job, "should be equal")
	}
}

func TestJob_Marshal(t *testing.T) {
	tests := []struct {
		args     []interface{}
		klass    string
		expected []byte
	}{
		{
			args:     []interface{}{},
			klass:    "",
			expected: []byte("{\"class\":\"\",\"args\":[]}"),
		},
		{
			args:     []interface{}{"a", 3, true},
			klass:    "test",
			expected: []byte("{\"class\":\"test\",\"args\":[\"a\",3,true]}"),
		},
		{
			args:     []interface{}{"a", "a", "a", 3.14159265},
			klass:    "pi",
			expected: []byte("{\"class\":\"pi\",\"args\":[\"a\",\"a\",\"a\",3.14159265]}"),
		},
		{
			args:     []interface{}{"{a: true}"},
			klass:    "json",
			expected: []byte("{\"class\":\"json\",\"args\":[\"{a: true}\"]}"),
		},
	}

	for _, test := range tests {
		job := NewJob(test.klass, test.args...)

		require.Equal(t, Job{
			Args:  test.args,
			Class: test.klass,
		}, job, "should be equal")

		buffer, err := job.Marshal()
		require.Nil(t, err, "should be nil")

		require.Equal(t, buffer, test.expected, "should be equal")
	}

}
