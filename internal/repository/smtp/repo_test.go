package smtp

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	"github.com/stretchr/testify/require"
	gomail "gopkg.in/gomail.v2"
)

type mockDialer struct {
	sent bool
}

func (m *mockDialer) DialAndSend(msg *gomail.Message) error {
	m.sent = true
	return nil
}

func TestSMTPSender_SendCode(t *testing.T) {
	t.Parallel()

	originalEmail := config.Config.SMTP.Email
	config.Config.SMTP.Email = "noreply@domdeli.ru"
	defer func() { config.Config.SMTP.Email = originalEmail }()

	type args struct {
		ctx   context.Context
		email string
		code  string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				ctx:   context.Background(),
				email: "user@example.com",
				code:  "123456",
			},
			wantErr: false,
		},
		{
			name: "valid short code",
			args: args{
				ctx:   context.Background(),
				email: "test@test.ru",
				code:  "ABC123",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			sender := NewSMTPSender("smtp.gmail.com", 587, "user@gmail.com", "pass")

			err := sender.SendCode(tt.args.ctx, tt.args.email, tt.args.code)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
