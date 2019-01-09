package imagist

import (
	"reflect"
	"testing"
)

func TestSetBackdropImage(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetBackdropImage(tt.args.path)
		})
	}
}

func TestSetWatermarkImage(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetWatermarkImage(tt.args.path)
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		chansize []int
	}
	tests := []struct {
		name string
		args args
		want *Imagist
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.chansize...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImagist_listen(t *testing.T) {
	type fields struct {
		jobs chan Job
		done chan string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Imagist{
				jobs: tt.fields.jobs,
				done: tt.fields.done,
			}
			i.listen()
		})
	}
}

func TestImagist_Add(t *testing.T) {
	type fields struct {
		jobs chan Job
		done chan string
	}
	type args struct {
		buf          []byte
		fileDiskPath string
		dimensions   *ImageDimensions
		validate     bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := Imagist{
				jobs: tt.fields.jobs,
				done: tt.fields.done,
			}
			if err := i.Add(tt.args.buf, tt.args.fileDiskPath, tt.args.dimensions, tt.args.validate); (err != nil) != tt.wantErr {
				t.Errorf("Imagist.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestImagist_execute(t *testing.T) {
	type fields struct {
		jobs chan Job
		done chan string
	}
	type args struct {
		j Job
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
			i := Imagist{
				jobs: tt.fields.jobs,
				done: tt.fields.done,
			}
			i.execute(tt.args.j)
		})
	}
}

func Test_imageProcess(t *testing.T) {
	type args struct {
		imgDiskPath string
		newWidth    int
		newHeight   int
		landscape   bool
		format      FormatDimensions
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
			if err := imageProcess(tt.args.imgDiskPath, tt.args.newWidth, tt.args.newHeight, tt.args.landscape, tt.args.format); (err != nil) != tt.wantErr {
				t.Errorf("imageProcess() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
