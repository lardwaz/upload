package option_test

import (
	"reflect"
	"testing"

	"go.lsl.digital/lardwaz/upload"
	"go.lsl.digital/lardwaz/upload/option"
)

func TestEvaluateImageOptions(t *testing.T) {
	tests := []struct {
		name string
		opts []func(upload.OptionsImage)
		want upload.OptionsImage
	}{
		{"empty", []func(upload.OptionsImage){}, option.NewImage()},
		{"nil", nil, option.NewImage()},
		{"min_width", []func(upload.OptionsImage){option.MinWidth(100)}, option.NewImage().SetMinWidth(100)},
		{"min_height", []func(upload.OptionsImage){option.MinHeight(100)}, option.NewImage().SetMinHeight(100)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := option.EvaluateImageOptions(tt.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EvaluateImageOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}
