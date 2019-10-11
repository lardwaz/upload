package option_test

import (
	"reflect"
	"testing"

	"go.lsl.digital/lardwaz/upload"
	"go.lsl.digital/lardwaz/upload/option"
)

func TestEvaluateWatermarkOptions(t *testing.T) {
	tests := []struct {
		name string
		opts []func(upload.OptionsWatermark)
		want upload.OptionsWatermark
	}{
		{"empty", []func(upload.OptionsWatermark){}, option.NewWatermark()},
		{"nil", nil, option.NewWatermark()},
		{"path", []func(upload.OptionsWatermark){option.WatermarkPath("/some/path")}, option.NewWatermark().SetPath("/some/path")},
		{"horizontal", []func(upload.OptionsWatermark){option.WatermarkHorizontal(100)}, option.NewWatermark().SetHorizontal(100)},
		{"vertical", []func(upload.OptionsWatermark){option.WatermarkVertical(100)}, option.NewWatermark().SetVertical(100)},
		{"offsetX", []func(upload.OptionsWatermark){option.WatermarkOffsetX(100)}, option.NewWatermark().SetOffsetX(100)},
		{"offsetY", []func(upload.OptionsWatermark){option.WatermarkOffsetY(100)}, option.NewWatermark().SetOffsetY(100)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := option.EvaluateWatermarkOptions(tt.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EvaluateWatermarkOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}
