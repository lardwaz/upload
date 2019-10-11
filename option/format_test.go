package option_test

import (
	"reflect"
	"testing"

	"go.lsl.digital/lardwaz/upload"
	"go.lsl.digital/lardwaz/upload/option"
)

func TestEvaluateFormatOptions(t *testing.T) {

	tests := []struct {
		name string
		opts []func(upload.OptionsFormat)
		want upload.OptionsFormat
	}{
		{"empty", []func(upload.OptionsFormat){}, option.NewFormat()},
		{"nil", nil, option.NewFormat()},
		{"format_name", []func(upload.OptionsFormat){option.FormatName("somename")}, option.NewFormat().SetName("somename")},
		{"format_width", []func(upload.OptionsFormat){option.FormatWidth(100)}, option.NewFormat().SetWidth(100)},
		{"format_height", []func(upload.OptionsFormat){option.FormatHeight(100)}, option.NewFormat().SetHeight(100)},
		{"format_backdrop", []func(upload.OptionsFormat){option.FormatBackdrop(option.BackdropPath("/abc/def"))}, option.NewFormat().SetBackdrop(option.BackdropPath("/abc/def"))},
		{"format_watermark", []func(upload.OptionsFormat){option.FormatWatermark(option.WatermarkPath("/abc/def"), option.WatermarkVertical(10))}, option.NewFormat().SetWatermark(option.WatermarkPath("/abc/def"), option.WatermarkVertical(10))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := option.EvaluateFormatOptions(tt.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EvaluateFormatOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}
