package nats

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// BitLocker NATS subject constants and pattern definitions.
const (
	// BitLockerEscrowSubjectPrefix is the NATS subject prefix for BitLocker recovery key escrow reports from agents.
	BitLockerEscrowSubjectPrefix = "agent.bitlocker.escrow."
	// BitLockerEscrowSubjectPattern is the NATS subject pattern to subscribe to all BitLocker key escrow reports.
	BitLockerEscrowSubjectPattern = "agent.bitlocker.escrow.*"

	// BitLockerForceEscrowSubjectPrefix is the NATS subject prefix to request an agent to force escrow its BitLocker keys.
	BitLockerForceEscrowSubjectPrefix = "agent.bitlocker.force_escrow."
	// BitLockerForceEscrowSubjectPattern is the NATS subject pattern to subscribe to all forced BitLocker escrow requests.
	BitLockerForceEscrowSubjectPattern = "agent.bitlocker.force_escrow.*"

	// BitLockerRotateSubjectPrefix is the NATS subject prefix for BitLocker recovery key rotation requests or reports.
	BitLockerRotateSubjectPrefix = "agent.bitlocker.rotate."
	// BitLockerRotateSubjectPattern is the NATS subject pattern to subscribe to all BitLocker recovery key rotation messages.
	BitLockerRotateSubjectPattern = "agent.bitlocker.rotate.*"
)

// BitLockerEscrowSubject returns the NATS subject for BitLocker escrow for a specific agent.
func BitLockerEscrowSubject(agentID string) string {
	return BitLockerEscrowSubjectPrefix + agentID
}

// BitLockerForceEscrowSubject returns the NATS subject to trigger forced BitLocker escrow for a specific agent.
func BitLockerForceEscrowSubject(agentID string) string {
	return BitLockerForceEscrowSubjectPrefix + agentID
}

// BitLockerRotateSubject returns the NATS subject for BitLocker key rotation for a specific agent.
func BitLockerRotateSubject(agentID string) string {
	return BitLockerRotateSubjectPrefix + agentID
}

// BitLockerEscrowRequest represents a payload sent by the agent to escrow a BitLocker recovery password.
type BitLockerEscrowRequest struct {
	AgentID          string    `json:"agent_id,omitempty"`
	DriveLetter      string    `json:"drive_letter,omitempty"`
	KeyID            string    `json:"key_id,omitempty"`
	RecoveryPassword string    `json:"recovery_password,omitempty"`
	CreatedAt        time.Time `json:"created_at,omitempty"`
}

// BitLockerEscrowResponse represents a response returned by the server when an agent escrows a BitLocker recovery key.
type BitLockerEscrowResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// NewBitLockerEscrowSuccessResponse creates a successful escrow response.
func NewBitLockerEscrowSuccessResponse() BitLockerEscrowResponse {
	return BitLockerEscrowResponse{Success: true}
}

// NewBitLockerEscrowErrorResponse creates an error escrow response.
func NewBitLockerEscrowErrorResponse(err string) BitLockerEscrowResponse {
	return BitLockerEscrowResponse{Success: false, Error: err}
}

// BitLockerForceEscrowRequest represents a command sent to an agent to trigger immediate escrow of BitLocker keys.
type BitLockerForceEscrowRequest struct {
	AgentID     string `json:"agent_id,omitempty"`
	DriveLetter string `json:"drive_letter,omitempty"`
	RequestedBy string `json:"requested_by,omitempty"`
}

// BitLockerRotateRequest represents a command sent to an agent to trigger BitLocker recovery key rotation.
type BitLockerRotateRequest struct {
	AgentID     string `json:"agent_id,omitempty"`
	DriveLetter string `json:"drive_letter,omitempty"`
	RequestedBy string `json:"requested_by,omitempty"`
}

// BitLockerRotateConfirmation represents a confirmation sent by an agent after rotating a BitLocker recovery key.
type BitLockerRotateConfirmation struct {
	AgentID     string    `json:"agent_id,omitempty"`
	DriveLetter string    `json:"drive_letter,omitempty"`
	OldKeyID    string    `json:"old_key_id,omitempty"`
	NewKeyID    string    `json:"new_key_id,omitempty"`
	RotatedAt   time.Time `json:"rotated_at,omitempty"`
}

func parseFlexibleTime(raw []byte) (time.Time, error) {
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		s = strings.TrimSpace(s)
		if s == "" {
			return time.Time{}, nil
		}
		formats := []string{
			time.RFC3339Nano,
			time.RFC3339,
			"2006-01-02T15:04:05",
			"2006-01-02 15:04:05",
			"2006-01-02",
		}
		for _, layout := range formats {
			if t, parseErr := time.Parse(layout, s); parseErr == nil {
				return t, nil
			}
		}
		return time.Time{}, fmt.Errorf("cannot parse time %q", s)
	}

	var t time.Time
	if err := json.Unmarshal(raw, &t); err != nil {
		return time.Time{}, err
	}
	return t, nil
}

// UnmarshalJSON supports snake_case, camelCase, and flexible time formats.
// UnmarshalJSON supports snake_case and Pascal/camelCase for BitLockerEscrowResponse.
func (r *BitLockerEscrowResponse) UnmarshalJSON(data []byte) error {
	type Alias BitLockerEscrowResponse
	aux := struct {
		*Alias
		SuccessAlt *bool  `json:"Success"`
		ErrorAlt   string `json:"Error"`
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.SuccessAlt != nil {
		r.Success = *aux.SuccessAlt
	}
	if r.Error == "" && aux.ErrorAlt != "" {
		r.Error = aux.ErrorAlt
	}
	return nil
}

func (r *BitLockerEscrowRequest) UnmarshalJSON(data []byte) error {
	type Alias BitLockerEscrowRequest
	aux := struct {
		*Alias
		AgentIDAlt          string          `json:"agentID"`
		AgentIdAlt          string          `json:"agentId"`
		DriveLetterAlt      string          `json:"driveLetter"`
		KeyIDAlt            string          `json:"keyID"`
		KeyIdAlt            string          `json:"keyId"`
		RecoveryPasswordAlt string          `json:"recoveryPassword"`
		CreatedAtRaw        json.RawMessage `json:"created_at"`
		CreatedAtAlt        json.RawMessage `json:"createdAt"`
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if r.AgentID == "" {
		if aux.AgentIDAlt != "" {
			r.AgentID = aux.AgentIDAlt
		} else if aux.AgentIdAlt != "" {
			r.AgentID = aux.AgentIdAlt
		}
	}
	if r.DriveLetter == "" && aux.DriveLetterAlt != "" {
		r.DriveLetter = aux.DriveLetterAlt
	}
	if r.KeyID == "" {
		if aux.KeyIDAlt != "" {
			r.KeyID = aux.KeyIDAlt
		} else if aux.KeyIdAlt != "" {
			r.KeyID = aux.KeyIdAlt
		}
	}
	if r.RecoveryPassword == "" && aux.RecoveryPasswordAlt != "" {
		r.RecoveryPassword = aux.RecoveryPasswordAlt
	}

	rawTime := aux.CreatedAtRaw
	if len(rawTime) == 0 {
		rawTime = aux.CreatedAtAlt
	}
	if len(rawTime) > 0 {
		t, err := parseFlexibleTime(rawTime)
		if err != nil {
			return err
		}
		r.CreatedAt = t
	}

	return nil
}

// UnmarshalJSON supports snake_case and camelCase.
func (r *BitLockerForceEscrowRequest) UnmarshalJSON(data []byte) error {
	type Alias BitLockerForceEscrowRequest
	aux := struct {
		*Alias
		AgentIDAlt     string `json:"agentID"`
		AgentIdAlt     string `json:"agentId"`
		DriveLetterAlt string `json:"driveLetter"`
		RequestedByAlt string `json:"requestedBy"`
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if r.AgentID == "" {
		if aux.AgentIDAlt != "" {
			r.AgentID = aux.AgentIDAlt
		} else if aux.AgentIdAlt != "" {
			r.AgentID = aux.AgentIdAlt
		}
	}
	if r.DriveLetter == "" && aux.DriveLetterAlt != "" {
		r.DriveLetter = aux.DriveLetterAlt
	}
	if r.RequestedBy == "" && aux.RequestedByAlt != "" {
		r.RequestedBy = aux.RequestedByAlt
	}
	return nil
}

// UnmarshalJSON supports snake_case and camelCase.
func (r *BitLockerRotateRequest) UnmarshalJSON(data []byte) error {
	type Alias BitLockerRotateRequest
	aux := struct {
		*Alias
		AgentIDAlt     string `json:"agentID"`
		AgentIdAlt     string `json:"agentId"`
		DriveLetterAlt string `json:"driveLetter"`
		RequestedByAlt string `json:"requestedBy"`
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if r.AgentID == "" {
		if aux.AgentIDAlt != "" {
			r.AgentID = aux.AgentIDAlt
		} else if aux.AgentIdAlt != "" {
			r.AgentID = aux.AgentIdAlt
		}
	}
	if r.DriveLetter == "" && aux.DriveLetterAlt != "" {
		r.DriveLetter = aux.DriveLetterAlt
	}
	if r.RequestedBy == "" && aux.RequestedByAlt != "" {
		r.RequestedBy = aux.RequestedByAlt
	}
	return nil
}

// UnmarshalJSON supports snake_case, camelCase, and flexible time formats.
func (r *BitLockerRotateConfirmation) UnmarshalJSON(data []byte) error {
	type Alias BitLockerRotateConfirmation
	aux := struct {
		*Alias
		AgentIDAlt     string          `json:"agentID"`
		AgentIdAlt     string          `json:"agentId"`
		DriveLetterAlt string          `json:"driveLetter"`
		OldKeyIDAlt    string          `json:"oldKeyID"`
		OldKeyIdAlt    string          `json:"oldKeyId"`
		NewKeyIDAlt    string          `json:"newKeyID"`
		NewKeyIdAlt    string          `json:"newKeyId"`
		RotatedAtRaw   json.RawMessage `json:"rotated_at"`
		RotatedAtAlt   json.RawMessage `json:"rotatedAt"`
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if r.AgentID == "" {
		if aux.AgentIDAlt != "" {
			r.AgentID = aux.AgentIDAlt
		} else if aux.AgentIdAlt != "" {
			r.AgentID = aux.AgentIdAlt
		}
	}
	if r.DriveLetter == "" && aux.DriveLetterAlt != "" {
		r.DriveLetter = aux.DriveLetterAlt
	}
	if r.OldKeyID == "" {
		if aux.OldKeyIDAlt != "" {
			r.OldKeyID = aux.OldKeyIDAlt
		} else if aux.OldKeyIdAlt != "" {
			r.OldKeyID = aux.OldKeyIdAlt
		}
	}
	if r.NewKeyID == "" {
		if aux.NewKeyIDAlt != "" {
			r.NewKeyID = aux.NewKeyIDAlt
		} else if aux.NewKeyIdAlt != "" {
			r.NewKeyID = aux.NewKeyIdAlt
		}
	}

	rawTime := aux.RotatedAtRaw
	if len(rawTime) == 0 {
		rawTime = aux.RotatedAtAlt
	}
	if len(rawTime) > 0 {
		t, err := parseFlexibleTime(rawTime)
		if err != nil {
			return err
		}
		r.RotatedAt = t
	}

	return nil
}
