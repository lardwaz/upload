package option_test

import (
	"reflect"
	"testing"

	"go.lsl.digital/lardwaz/upload"
	"go.lsl.digital/lardwaz/upload/option"
)

func TestEvaluateBackdropOptions(t *testing.T) {
	tests := []struct {
		name string
		opts []func(upload.OptionsBackdrop)
		want upload.OptionsBackdrop
	}{
		{"empty", []func(upload.OptionsBackdrop){}, option.NewBackdrop()},
		{"nil", nil, option.NewBackdrop()},
		{"normal", []func(upload.OptionsBackdrop){option.BackdropPath("/test/something")}, option.NewBackdrop().SetPath("/test/something")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := option.EvaluateBackdropOptions(tt.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EvaluateBackdropOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}
