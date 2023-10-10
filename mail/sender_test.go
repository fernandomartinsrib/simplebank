package mail

import (
	"testing"

	"github.com/fernandomartinsrib/simplebank/utils"
	"github.com/stretchr/testify/require"
)

func TestSendEmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	config, err := utils.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "A test email - paygo"
	content := `
	<h1>Hello World!</h1>
	<p> This is a test message from Fernando payGO</p>
	`

	to := []string{"fernandomartinsrib@gmail.com"}
	attachFiles := []string{"../README.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}
