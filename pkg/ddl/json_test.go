package ddl

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestApplyWithKeyRecursively(t *testing.T) {
	type args struct {
		jsonData map[string]interface{}
		keys     []string
		f        func(v interface{}) interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "normal",
			args: args{
				jsonData: map[string]interface{}{
					"foo": map[string]interface{}{
						"boo": "hoge",
					},
				},
				keys: []string{"foo", "boo"},
				f: func(value interface{}) interface{} {
					return "fuga"
				},
			},
			want: map[string]interface{}{
				"foo": map[string]interface{}{
					"boo": "fuga",
				},
			},
		},
		{
			name: "array",
			args: args{
				jsonData: map[string]interface{}{
					"foo": []interface{}{
						"one", "two", "three",
					},
				},
				keys: []string{"foo"},
				f: func(value interface{}) interface{} {
					return "_" + value.(string)
				},
			},
			want: map[string]interface{}{
				"foo": []interface{}{
					"_one", "_two", "_three",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := tt.args.jsonData
			if err := applyWithKeyRecursively(data, tt.args.keys, tt.args.f); err != nil {
				t.Error(err.Error())
				return
			}
			if diff := cmp.Diff(data, tt.want); diff != "" {
				t.Error(diff)
				return
			}
		})
	}
}

func TestApplyRecursively(t *testing.T) {
	type args struct {
		data map[string]interface{}
		f    func(v interface{}) interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "normal",
			args: args{
				data: map[string]interface{}{
					"foo": map[string]interface{}{
						"boo": "hoge",
						"moo": "ushi",
					},
				},
				f: func(value interface{}) interface{} {
					return "fuga"
				},
			},
			want: map[string]interface{}{
				"foo": map[string]interface{}{
					"boo": "fuga",
					"moo": "fuga",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := tt.args.data
			if err := applyRecursively(data, tt.args.f); err != nil {
				t.Error(err.Error())
				return
			}
			if diff := cmp.Diff(data, tt.want); diff != "" {
				t.Error(diff)
				return
			}
		})
	}
}
