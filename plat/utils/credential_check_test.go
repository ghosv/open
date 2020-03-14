package utils

import (
	"testing"
)

func TestIsEmail(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"email-1", args{"1@1.com"}, true},
		{"email-2", args{"v.com"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmail(tt.args.s); got != tt.want {
				t.Errorf("IsEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsPhone(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"phone-1", args{"15012345678"}, true},
		{"phone-2", args{"123"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPhone(tt.args.s); got != tt.want {
				t.Errorf("IsPhone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckUsername(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name      string
		args      args
		wantName  string
		wantPhone string
		wantEmail string
		wantOk    bool
	}{
		{"username-1", args{"1@1.com"}, "", "", "1@1.com", true},
		{"username-2", args{"v.com"}, "", "", "", false},
		{"username-3", args{"15012345678"}, "", "15012345678", "", true},
		{"username-4", args{"123"}, "", "", "", false},
		{"username-4", args{"a123"}, "a123", "", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotPhone, gotEmail, gotOk := CheckUsername(tt.args.s)
			if gotName != tt.wantName {
				t.Errorf("CheckUsername() gotName = %v, want %v", gotName, tt.wantName)
			}
			if gotPhone != tt.wantPhone {
				t.Errorf("CheckUsername() gotPhone = %v, want %v", gotPhone, tt.wantPhone)
			}
			if gotEmail != tt.wantEmail {
				t.Errorf("CheckUsername() gotEmail = %v, want %v", gotEmail, tt.wantEmail)
			}
			if gotOk != tt.wantOk {
				t.Errorf("CheckUsername() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
