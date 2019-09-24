package option_test

import (
	"reflect"
	"testing"

	"go.lsl.digital/lardwaz/upload"
	"go.lsl.digital/lardwaz/upload/option"
	"go.lsl.digital/lardwaz/upload/types"
)

func TestEvaluateOptions(t *testing.T) {
	tests := []struct {
		name string
		opts []func(upload.Options)
		want upload.Options
	}{
		{"empty", []func(upload.Options){}, option.NewUpload()},
		{"nil", nil, option.NewUpload()},
		{"dir", []func(upload.Options){option.Dir("/abc/def")}, option.NewUpload().SetDir("/abc/def")},
		{"destination", []func(upload.Options){option.Destination("/edf/bac")}, option.NewUpload().SetDestination("/edf/bac")},
		{"media_prefix_url", []func(upload.Options){option.MediaPrefixURL("/roose/bog")}, option.NewUpload().SetMediaPrefixURL("/roose/bog")},
		{"file_type", []func(upload.Options){option.FileType(types.TypeJPEG)}, option.NewUpload().AddFileType(types.TypeJPEG)},
		{"multiple file_type", []func(upload.Options){option.FileType(types.TypeJPEG), option.FileType(types.TypeMP4)}, option.NewUpload().AddFileType(types.TypeJPEG).AddFileType(types.TypeMP4)},
		{"max_size", []func(upload.Options){option.MaxSize(1000)}, option.NewUpload().SetMaxSize(1000)},
		{"convert_to", []func(upload.Options){option.ConvertTo(types.TypeMP3, types.TypeAAC)}, option.NewUpload().SetConvertTo(types.TypeMP3, types.TypeAAC)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := option.EvaluateOptions(tt.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EvaluateOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}
