package configs

import (
	"errors"
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/018bf/example/internal/domain/errs"
)

func TestParseConfig(t *testing.T) {
	file := `
bind_addr = ":8005"
log_level = "debug"

	`
	badFile := `
bind_addr = ":8003"
log_level = 2
	`
	configPath := path.Join(os.TempDir(), "config.toml")
	badConfigPath := path.Join(os.TempDir(), "bad-config.toml")
	if err := os.WriteFile(configPath, []byte(file), 0600); err != nil {
		t.Fatal(err)
		return
	}
	if err := os.WriteFile(badConfigPath, []byte(badFile), 0600); err != nil {
		t.Fatal(err)
		return
	}
	type args struct {
		configPath string
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr error
	}{
		{
			name: "ok",
			args: args{
				configPath: configPath,
			},
			want: &Config{
				BindAddr: ":8005",
				LogLevel: "debug",
				Database: database{
					URI:                "",
					MaxOpenConnections: 50,
					MaxIDLEConnections: 10,
				},
				Auth: auth{
					PublicKey:  "",
					PrivateKey: "",
					RefreshTTL: 172800,
					AccessTTL:  86400,
				},
			},
			wantErr: nil,
		},
		{
			name: "ok from env",
			args: args{
				configPath: "",
			},
			want: &Config{
				BindAddr: ":8000",
				LogLevel: "debug",
				Database: database{
					URI:                "",
					MaxOpenConnections: 50,
					MaxIDLEConnections: 10,
				},
				Auth: auth{
					PublicKey:  "",
					PrivateKey: "",
					RefreshTTL: 172800,
					AccessTTL:  86400,
				},
			},
			wantErr: nil,
		},
		{
			name: "bad config",
			args: args{
				configPath: badConfigPath,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("config file parsing error: toml: cannot load TOML value of type int64 into a Go string"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseConfig(tt.args.configPath)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ParseConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
