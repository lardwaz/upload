package upload

import (
	"net/http"
	"testing"
)

func TestUpload(t *testing.T) {
	type args struct {
		fileName    string
		fileContent []byte
		options     *Options
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := Upload(tt.args.fileName, tt.args.fileContent, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("Upload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Upload() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Upload() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	type args struct {
		fp string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Delete(tt.args.fp); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_buildFileName(t *testing.T) {
	type args struct {
		oldFilename string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildFileName(tt.args.oldFilename); got != tt.want {
				t.Errorf("buildFileName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_changeExt(t *testing.T) {
	type args struct {
		fileDiskPath string
		filepath     string
		newExt       string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := changeExt(tt.args.fileDiskPath, tt.args.filepath, tt.args.newExt)
			if (err != nil) != tt.wantErr {
				t.Errorf("changeExt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("changeExt() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("changeExt() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_httpImageDirHandler_ServeHTTP(t *testing.T) {
	type fields struct {
		root   http.FileSystem
		prefix string
		opts   *Options
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := httpImageDirHandler{
				root:   tt.fields.root,
				prefix: tt.fields.prefix,
				opts:   tt.fields.opts,
			}
			s.ServeHTTP(tt.args.w, tt.args.r)
		})
	}
}
