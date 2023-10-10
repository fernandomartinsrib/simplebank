package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	db "github.com/fernandomartinsrib/simplebank/db/sqlc"
	"github.com/fernandomartinsrib/simplebank/utils"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const TaskSenderVerifyEmail = "task:send_verify_email"

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

const (
	environment = "local"
)

func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(ctx context.Context, payload *PayloadSendVerifyEmail, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(TaskSenderVerifyEmail, jsonPayload, opts...)
	taskInfo, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Info().
		Type("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("queue", taskInfo.Queue).
		Int("max_retry", taskInfo.MaxRetry).
		Msg("enqueued task")

	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	user, err := processor.store.GetUser(ctx, payload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user does not exist: %w", asynq.SkipRetry)
		}

		return fmt.Errorf("failed to get user: %w", err)
	}

	verifyEmail, err := processor.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
		Username:   user.Username,
		Email:      user.Email,
		SecretCode: utils.RandomString(32),
	})

	if err != nil {
		return fmt.Errorf("failed to create verify email: %w", err)
	}

	subject := "Welcome to payGo!"

	domain := "http://paygo-bank.net"
	if environment == "local" {
		domain = "http://localhost:8080"
	}

	verifyUrl := fmt.Sprintf("%s/v1/verify_email?email_id=%d&secret_code=%s", domain, verifyEmail.ID, verifyEmail.SecretCode)

	to := []string{user.Email}

	content := fmt.Sprintf(`Hello %s,<br/>
	Thank you for registering with us!<br/>
	Please <a href="%s">click here</a> to verify your email address.<br/>
	`, user.Fullname, verifyUrl)

	err = processor.mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send verify email: %w", err)
	}

	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("email", user.Email).
		Msg("processed task")

	return nil
}
