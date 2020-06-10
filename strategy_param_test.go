package hazana

import "testing"

func Test_strategyParameters_intParam(t *testing.T) {
	type fields struct {
		line string
	}
	type args struct {
		name   string
		absent int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{"no params",
			fields{"method"},
			args{"key", -1},
			-1,
		},
		{"one params",
			fields{"method key=5"},
			args{"key", -1},
			5,
		},
		{"neg param",
			fields{"method key=5 max-factor=-1"},
			args{"max-factor", 0},
			-1,
		},
		{"no int param",
			fields{"method key=5 chord=C#"},
			args{"chord", 0},
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := strategyParameters{
				line: tt.fields.line,
			}
			if got := c.intParam(tt.args.name, tt.args.absent); got != tt.want {
				t.Errorf("strategyParameters.intParam() = %v, want %v", got, tt.want)
			}
		})
	}
}
