package utils

import (
	"testing"
	"time"
)

func TestGetDateFromString(t *testing.T) {
	test_suits := []struct {
		input  string
		output time.Time
	}{
		{input: "3 hours ago", output: time.Now().Add(-3 * time.Hour)},
		{input: "1 hour ago", output: time.Now().Add(-1 * time.Hour)},
		{input: "3 minutes ago", output: time.Now().Add(-3 * time.Minute)},
		{input: "1 minute ago", output: time.Now().Add(-1 * time.Minute)},
		{input: "3 days ago", output: time.Now().Add(-3 * 24 * time.Hour)},
		{input: "1 day ago", output: time.Now().Add(-1 * 24 * time.Hour)},
	}

	for _, test := range test_suits {
		out, err := GetDateFromString(test.input)
		if err != nil {
			t.Errorf("%s", err.Error())
		}
		// 不能超过10s的延迟
		if (out.Unix()-test.output.Unix()) > 10 || (out.Unix()-test.output.Unix()) < (-10) {
			t.Errorf("Test Error expect %s, but got %s", test.output, out)
		}
	}
}

func BenchmarkGetDateFromString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := GetDateFromString("3 hours ago")
		if err != nil {
			b.Error(err)
		}
	}
}

func TestGetMD5Hash(t *testing.T) {
	test_suits := []struct {
		input  string
		output string
	}{
		{input: "Hello this is mike", output: "2a30929c32fdf03f884ef462610bc3f5"},
		{input: "", output: "d41d8cd98f00b204e9800998ecf8427e"},
		{input: "Hello,my name is not mike", output: "2cd00961fb3ada34779c31260f3382b5"},
	}
	for _, test := range test_suits {
		out := GetMD5Hash(test.input)
		if out != test.output {
			t.Errorf("Expect %s, But got %s", test.output, out)
		}
	}
}
