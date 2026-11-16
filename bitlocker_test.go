package nats

import (
	"encoding/json"
	"testing"
	"time"
)

func TestBitLockerSubjects(t *testing.T) {
	agentID := "e578c772-243e-434a-9c71-bdf47321033a"

	escrowSubject := BitLockerEscrowSubject(agentID)
	expectedEscrow := "agent.bitlocker.escrow.e578c772-243e-434a-9c71-bdf47321033a"
	if escrowSubject != expectedEscrow {
		t.Fatalf("expected escrow subject %q, got %q", expectedEscrow, escrowSubject)
	}
	if BitLockerEscrowSubjectPrefix != "agent.bitlocker.escrow." {
		t.Fatalf("unexpected escrow subject prefix: %q", BitLockerEscrowSubjectPrefix)
	}
	if BitLockerEscrowSubjectPattern != "agent.bitlocker.escrow.*" {
		t.Fatalf("unexpected escrow subject pattern: %q", BitLockerEscrowSubjectPattern)
	}

	forceSubject := BitLockerForceEscrowSubject(agentID)
	expectedForce := "agent.bitlocker.force_escrow.e578c772-243e-434a-9c71-bdf47321033a"
	if forceSubject != expectedForce {
		t.Fatalf("expected force escrow subject %q, got %q", expectedForce, forceSubject)
	}
	if BitLockerForceEscrowSubjectPrefix != "agent.bitlocker.force_escrow." {
		t.Fatalf("unexpected force escrow subject prefix: %q", BitLockerForceEscrowSubjectPrefix)
	}
	if BitLockerForceEscrowSubjectPattern != "agent.bitlocker.force_escrow.*" {
		t.Fatalf("unexpected force escrow subject pattern: %q", BitLockerForceEscrowSubjectPattern)
	}

	rotateSubject := BitLockerRotateSubject(agentID)
	expectedRotate := "agent.bitlocker.rotate.e578c772-243e-434a-9c71-bdf47321033a"
	if rotateSubject != expectedRotate {
		t.Fatalf("expected rotate subject %q, got %q", expectedRotate, rotateSubject)
	}
	if BitLockerRotateSubjectPrefix != "agent.bitlocker.rotate." {
		t.Fatalf("unexpected rotate subject prefix: %q", BitLockerRotateSubjectPrefix)
	}
	if BitLockerRotateSubjectPattern != "agent.bitlocker.rotate.*" {
		t.Fatalf("unexpected rotate subject pattern: %q", BitLockerRotateSubjectPattern)
	}
}

func TestBitLockerEscrowRequestJSON(t *testing.T) {
	now := time.Date(2026, 11, 16, 6, 30, 0, 0, time.UTC)
	req := BitLockerEscrowRequest{
		AgentID:          "agent-123",
		DriveLetter:      "C:",
		KeyID:            "{ABCD-1234-EFGH-5678}",
		RecoveryPassword: "111111-222222-333333-444444-555555-666666-777777-888888",
		CreatedAt:        now,
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal BitLockerEscrowRequest: %v", err)
	}

	var parsed BitLockerEscrowRequest
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to unmarshal BitLockerEscrowRequest: %v", err)
	}

	if parsed.AgentID != req.AgentID {
		t.Errorf("AgentID mismatch: expected %s, got %s", req.AgentID, parsed.AgentID)
	}
	if parsed.DriveLetter != req.DriveLetter {
		t.Errorf("DriveLetter mismatch: expected %s, got %s", req.DriveLetter, parsed.DriveLetter)
	}
	if parsed.KeyID != req.KeyID {
		t.Errorf("KeyID mismatch: expected %s, got %s", req.KeyID, parsed.KeyID)
	}
	if parsed.RecoveryPassword != req.RecoveryPassword {
		t.Errorf("RecoveryPassword mismatch: expected %s, got %s", req.RecoveryPassword, parsed.RecoveryPassword)
	}
	if !parsed.CreatedAt.Equal(now) {
		t.Errorf("CreatedAt mismatch: expected %v, got %v", now, parsed.CreatedAt)
	}

	// Test with camelCase JSON and custom string timestamp
	camelJSON := []byte(`{
		"agentID": "agent-456",
		"driveLetter": "D:",
		"keyID": "{KEY-999}",
		"recoveryPassword": "pass-999",
		"createdAt": "2026-11-16 06:30:00"
	}`)
	var camelParsed BitLockerEscrowRequest
	if err := json.Unmarshal(camelJSON, &camelParsed); err != nil {
		t.Fatalf("failed to unmarshal camelCase BitLockerEscrowRequest: %v", err)
	}
	if camelParsed.AgentID != "agent-456" {
		t.Errorf("expected agent-456, got %s", camelParsed.AgentID)
	}
	if camelParsed.DriveLetter != "D:" {
		t.Errorf("expected D:, got %s", camelParsed.DriveLetter)
	}
	if camelParsed.KeyID != "{KEY-999}" {
		t.Errorf("expected {KEY-999}, got %s", camelParsed.KeyID)
	}
	if camelParsed.RecoveryPassword != "pass-999" {
		t.Errorf("expected pass-999, got %s", camelParsed.RecoveryPassword)
	}
	if camelParsed.CreatedAt.Year() != 2026 {
		t.Errorf("expected year 2026, got %v", camelParsed.CreatedAt)
	}
}

func TestBitLockerForceEscrowRequestJSON(t *testing.T) {
	req := BitLockerForceEscrowRequest{
		AgentID:     "agent-force-1",
		DriveLetter: "C:",
		RequestedBy: "admin@openuem.internal",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal BitLockerForceEscrowRequest: %v", err)
	}

	var parsed BitLockerForceEscrowRequest
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to unmarshal BitLockerForceEscrowRequest: %v", err)
	}

	if parsed.AgentID != req.AgentID || parsed.DriveLetter != req.DriveLetter || parsed.RequestedBy != req.RequestedBy {
		t.Fatalf("data mismatch: expected %+v, got %+v", req, parsed)
	}

	// Test camelCase
	camelJSON := []byte(`{
		"agentId": "agent-force-2",
		"driveLetter": "E:",
		"requestedBy": "user2"
	}`)
	var camelParsed BitLockerForceEscrowRequest
	if err := json.Unmarshal(camelJSON, &camelParsed); err != nil {
		t.Fatalf("failed to unmarshal camelCase BitLockerForceEscrowRequest: %v", err)
	}
	if camelParsed.AgentID != "agent-force-2" || camelParsed.DriveLetter != "E:" || camelParsed.RequestedBy != "user2" {
		t.Fatalf("unexpected camelCase parsed values: %+v", camelParsed)
	}
}

func TestBitLockerRotateRequestJSON(t *testing.T) {
	req := BitLockerRotateRequest{
		AgentID:     "agent-rot-1",
		DriveLetter: "C:",
		RequestedBy: "security-auditor",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal BitLockerRotateRequest: %v", err)
	}

	var parsed BitLockerRotateRequest
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to unmarshal BitLockerRotateRequest: %v", err)
	}

	if parsed.AgentID != req.AgentID || parsed.DriveLetter != req.DriveLetter || parsed.RequestedBy != req.RequestedBy {
		t.Fatalf("data mismatch: expected %+v, got %+v", req, parsed)
	}
}

func TestBitLockerRotateConfirmationJSON(t *testing.T) {
	now := time.Date(2026, 11, 16, 7, 0, 0, 0, time.UTC)
	conf := BitLockerRotateConfirmation{
		AgentID:     "agent-rotate-confirm",
		DriveLetter: "C:",
		OldKeyID:    "{OLD-KEY-111}",
		NewKeyID:    "{NEW-KEY-222}",
		RotatedAt:   now,
	}

	data, err := json.Marshal(conf)
	if err != nil {
		t.Fatalf("failed to marshal BitLockerRotateConfirmation: %v", err)
	}

	var parsed BitLockerRotateConfirmation
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to unmarshal BitLockerRotateConfirmation: %v", err)
	}

	if parsed.AgentID != conf.AgentID || parsed.OldKeyID != conf.OldKeyID || parsed.NewKeyID != conf.NewKeyID {
		t.Fatalf("data mismatch: expected %+v, got %+v", conf, parsed)
	}
	if !parsed.RotatedAt.Equal(now) {
		t.Fatalf("RotatedAt mismatch: expected %v, got %v", now, parsed.RotatedAt)
	}

	// Test camelCase
	camelJSON := []byte(`{
		"agentID": "agent-rotate-confirm-2",
		"driveLetter": "D:",
		"oldKeyID": "{OLD-333}",
		"newKeyID": "{NEW-444}",
		"rotatedAt": "2026-11-16T07:15:00Z"
	}`)
	var camelParsed BitLockerRotateConfirmation
	if err := json.Unmarshal(camelJSON, &camelParsed); err != nil {
		t.Fatalf("failed to unmarshal camelCase BitLockerRotateConfirmation: %v", err)
	}
	if camelParsed.AgentID != "agent-rotate-confirm-2" || camelParsed.OldKeyID != "{OLD-333}" || camelParsed.NewKeyID != "{NEW-444}" {
		t.Fatalf("unexpected camelParsed: %+v", camelParsed)
	}
	if camelParsed.RotatedAt.Minute() != 15 {
		t.Fatalf("unexpected rotatedAt: %v", camelParsed.RotatedAt)
	}
}

func TestBitLockerEscrowResponseJSON(t *testing.T) {
	// Success response
	respSuccess := NewBitLockerEscrowSuccessResponse()
	data, err := json.Marshal(respSuccess)
	if err != nil {
		t.Fatalf("failed to marshal BitLockerEscrowResponse: %v", err)
	}

	var parsed BitLockerEscrowResponse
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to unmarshal BitLockerEscrowResponse: %v", err)
	}
	if !parsed.Success {
		t.Errorf("expected success true, got false")
	}
	if parsed.Error != "" {
		t.Errorf("expected empty error, got %s", parsed.Error)
	}

	// Error response
	respErr := NewBitLockerEscrowErrorResponse("database connection failed")
	dataErr, err := json.Marshal(respErr)
	if err != nil {
		t.Fatalf("failed to marshal error BitLockerEscrowResponse: %v", err)
	}

	var parsedErr BitLockerEscrowResponse
	if err := json.Unmarshal(dataErr, &parsedErr); err != nil {
		t.Fatalf("failed to unmarshal error BitLockerEscrowResponse: %v", err)
	}
	if parsedErr.Success {
		t.Errorf("expected success false, got true")
	}
	if parsedErr.Error != "database connection failed" {
		t.Errorf("expected error message mismatch, got %s", parsedErr.Error)
	}

	// Test PascalCase / alt keys
	altJSON := []byte(`{"Success": true, "Error": ""}`)
	var parsedAlt BitLockerEscrowResponse
	if err := json.Unmarshal(altJSON, &parsedAlt); err != nil {
		t.Fatalf("failed to unmarshal alt BitLockerEscrowResponse: %v", err)
	}
	if !parsedAlt.Success {
		t.Errorf("expected success true from alt json")
	}
}
