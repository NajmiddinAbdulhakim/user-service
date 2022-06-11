package postgres

import (
	"reflect"
	"testing"

	pb "github.com/NajmiddinAbdulhakim/user-service/genproto"
)

func TestUserRepo_Create(t *testing.T) {
	tests := []struct {
		name string
		input pb.User
		want pb.User
		wantErr bool
	}{
		{
			name: `success case`,
			input: pb.User{
				FirstName: "test usersfdsa",
				LastName: "test",
			}, 
			want : pb.User{
				FirstName: "test user",
				LastName: "test",
			},
			wantErr: false, 
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := repo.CreateUser(&tc.input)
			if err != nil {
				t.Fatalf(`%s: expected: %v, got: %v`,tc.name, tc.wantErr, err)
			}
			got.Id = ""
			if !reflect.DeepEqual(&tc.want, got) {
				t.Fatalf(`%s: expected: %v, got: %v`,tc.name, tc.want, err)

			}
		})
	}
}