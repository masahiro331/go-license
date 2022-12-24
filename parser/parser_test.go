package parser

import (
	"reflect"
	"testing"

	"github.com/masahiro331/go-license/lexer"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		normFunc  NormalizeFunc
		want      *LicenseExpression
		wantStr   string
		expectErr string
	}{
		{
			name:  "happy path single license",
			input: "Public Domain",
			want: &LicenseExpression{
				Node: Node{
					License: "Public Domain",
				},
			},
			wantStr: "Public Domain",
		},
		{
			name:  "happy path single license with norm func",
			input: "Public Domain",
			want: &LicenseExpression{
				Node: Node{
					License: "Public Domain",
				},
			},
			normFunc: func(n string) string {
				if n == "Public Domain" {
					return "Unlicense"
				}
				return n
			},
			wantStr: "Unlicense",
		},
		{
			name:  "happy path 2",
			input: "Public ._+-",
			want: &LicenseExpression{
				Node: Node{
					License: "Public ._+-",
				},
			},
			wantStr: "Public ._+-",
		},
		{
			name:  "happy path multi license",
			input: "Public Domain AND ( GPLv2+ or AFL ) AND LGPLv2+ with distribution exceptions",
			want: &LicenseExpression{
				Node: Node{
					License: "Public Domain",
				},
				Operator: "AND",
				Next: &LicenseExpression{
					Node: Node{
						LicenseExpression: &LicenseExpression{
							Node: Node{
								License: "GPLv2+",
							},
							Operator: "OR",
							Next: &LicenseExpression{
								Node: Node{
									License: "AFL",
								},
							},
						},
					},
					Operator: "AND",
					Next: &LicenseExpression{
						Node: Node{
							License: "LGPLv2+ with distribution exceptions",
						},
					},
				},
			},
			wantStr: "Public Domain AND ( GPLv2+ OR AFL ) AND LGPLv2+ with distribution exceptions",
		},
		{
			name:  "happy path nested license",
			input: "Public Domain AND ( GPLv2+ or AFL AND ( CC0 or LGPL1.0) )",
			want: &LicenseExpression{
				Node: Node{
					License: "Public Domain",
				},
				Operator: "AND",
				Next: &LicenseExpression{
					Node: Node{
						LicenseExpression: &LicenseExpression{
							Node: Node{
								License: "GPLv2+",
							},
							Operator: "OR",
							Next: &LicenseExpression{
								Node: Node{
									License: "AFL",
								},
								Operator: "AND",
								Next: &LicenseExpression{
									Node: Node{
										LicenseExpression: &LicenseExpression{
											Node: Node{
												License: "CC0",
											},
											Operator: "OR",
											Next: &LicenseExpression{
												Node: Node{
													License: "LGPL1.0",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			wantStr: "Public Domain AND ( GPLv2+ OR AFL AND ( CC0 OR LGPL1.0 ) )",
		},
		{
			name:  "happy path 2",
			input: "( GPLv2+ or CC0 )",
			want: &LicenseExpression{
				Node: Node{
					LicenseExpression: &LicenseExpression{
						Node: Node{
							License: "GPLv2+",
						},
						Operator: "OR",
						Next: &LicenseExpression{
							Node: Node{
								License: "CC0",
							},
						},
					},
				},
			},
			wantStr: "( GPLv2+ OR CC0 )",
		},
		{
			name:      "bad path close bracket not found",
			input:     "Public Domain AND ( GPLv2+ ",
			expectErr: "invalid expression error",
		},
		{
			name:      "bad path bad bracket",
			input:     "Public Domain AND { ( GPLv2+ } )",
			expectErr: "invalid expression error",
		},
		{
			name:      "bad path bad bracket",
			input:     "Public Domain AND  ({GPLv2+)}",
			expectErr: "invalid expression error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l).RegisterNormalizeFunc(tt.normFunc)

			got, err := p.Parse()
			if tt.expectErr == "" && err != nil {
				t.Errorf(err.Error())
			}
			if tt.expectErr != "" && tt.expectErr != err.Error() {
				t.Errorf("err: %s, want %s", err.Error(), tt.expectErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
			if p.Normalize(got) != tt.wantStr {
				t.Errorf("Normalize() = %v, want %v", p.Normalize(got), tt.wantStr)
			}
		})
	}
}
