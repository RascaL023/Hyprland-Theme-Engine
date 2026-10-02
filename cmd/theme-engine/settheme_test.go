package main

import "testing"

func TestParseSetThemeArgs(t *testing.T) {
	cases := []struct {
		name      string
		args      []string
		wantName  string
		wantType  string
		wantErr   bool
	}{
		{name: "no args", args: nil, wantErr: true},
		{name: "flag first", args: []string{"--type", "dark"}, wantErr: true},
		{name: "name only", args: []string{"harbor"}, wantName: "harbor"},
		{name: "name with dark", args: []string{"harbor", "--type", "dark"}, wantName: "harbor", wantType: "dark"},
		{name: "name with light", args: []string{"harbor", "--type", "light"}, wantName: "harbor", wantType: "light"},
		{name: "bad type", args: []string{"harbor", "--type", "dim"}, wantErr: true},
		{name: "type without value", args: []string{"harbor", "--type"}, wantErr: true},
		{name: "unexpected arg", args: []string{"harbor", "extra"}, wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotName, gotType, err := parseSetThemeArgs(tc.args)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got name=%q type=%q", gotName, gotType)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gotName != tc.wantName {
				t.Errorf("name = %q, want %q", gotName, tc.wantName)
			}
			if gotType != tc.wantType {
				t.Errorf("type = %q, want %q", gotType, tc.wantType)
			}
		})
	}
}
