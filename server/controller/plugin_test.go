package controller

import "testing"

func TestShouldRecordPluginAccess(t *testing.T) {
	tests := []struct {
		name   string
		status int
		want   bool
	}{
		{name: "approved", status: 0, want: true},
		{name: "rejected", status: 1, want: false},
		{name: "pending", status: 2, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldRecordPluginAccess(tt.status); got != tt.want {
				t.Errorf("shouldRecordPluginAccess(%d) = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}
