package user

import "testing"

func Test_validateRegisterCredential(t *testing.T) {
	type args struct {
		login string
		pass  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				login: "test",
				pass:  "testTest",
			},
			wantErr: false,
		},
		{
			name: "login_short",
			args: args{
				login: "t",
				pass:  "testTest",
			},
			wantErr: true,
		},
		{
			name: "password_short",
			args: args{
				login: "test",
				pass:  "test",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateRegisterCredential(tt.args.login, tt.args.pass); (err != nil) != tt.wantErr {
				t.Errorf("validateRegisterCredential() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
