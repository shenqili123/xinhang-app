package email

import (
	"testing"

	"xinhang-backend/config"
)

func TestInit_NoConfig(t *testing.T) {
	Init(&config.Config{})
	if IsEnabled() {
		t.Error("expected disabled when SMTP not configured")
	}
}

func TestInit_WithConfig(t *testing.T) {
	Init(&config.Config{
		SMTPHost: "smtp.test.com",
		SMTPUser: "test@test.com",
		SMTPPassword: "pwd",
		SMTPPort: "465",
		SMTPFrom: "Test",
	})
	if !IsEnabled() {
		t.Error("expected enabled when SMTP configured")
	}
	cfg = nil
}

func TestIsEnabled_NilCfg(t *testing.T) {
	cfg = nil
	if IsEnabled() {
		t.Error("expected disabled when cfg is nil")
	}
}

func TestSendVerificationCode_Disabled(t *testing.T) {
	cfg = nil
	err := SendVerificationCode("test@test.com", "123456")
	if err == nil {
		t.Error("expected error when SMTP not configured")
	}
}
