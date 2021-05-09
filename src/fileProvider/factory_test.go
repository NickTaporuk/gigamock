package fileProvider

import (
	"reflect"
	"testing"

	"github.com/NickTaporuk/gigamock/src/fileType"
)

func TestFactory(t *testing.T) {
	type args struct {
		ext string
	}
	tests := []struct {
		name    string
		args    args
		want    FileProvider
		wantErr bool
	}{
		{
			name:    "positive scenario should be retrieved existing provider yaml",
			args:    args{ext: fileType.FileExtensionYAML},
			want:    NewYAMLProvider(),
			wantErr: false,
		},
		{
			name:    "positive scenario should be retrieved existing provider json",
			args:    args{ext: fileType.FileExtensionJSON},
			want:    NewJSONProvider(),
			wantErr: false,
		},
		{
			name:    "negative scenario should be retrieved the error ",
			args:    args{ext: ""},
			wantErr: true,
		},
		{
			name:    "negative scenario unknown extension should be retrieved the error ",
			args:    args{ext: "TEST"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Factory(tt.args.ext)
			if (err != nil) != tt.wantErr {
				t.Errorf("Factory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Factory() got = %v, want %v", got, tt.want)
			}
		})
	}
}
