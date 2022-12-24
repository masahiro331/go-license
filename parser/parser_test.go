package parser

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		lex *lexer.Lexer
	}
	tests := []struct {
		name string
		args args
		want *LicenseExpression
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Parse(tt.args.lex); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
