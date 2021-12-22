package setting

import "testing"

func TestSetup(t *testing.T) {
	Setup()
}

func Test_getAbsolutePath(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
		{
			name: "Test",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			println(getAbsolutePath(2))
		})
	}
}
