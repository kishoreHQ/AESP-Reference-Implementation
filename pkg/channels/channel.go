package channels

import "github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"

// Plugin is a registry-discovered notification/HITL channel (ADT-11).
// Vendor-specific adapters live under plugins/channels/*/ — not kernel.
type Plugin interface {
	ID() string
	SendApproval(taskID, summary string) error
	SendNotification(text string) error
}

func SpecMapping() types.SpecMapping {
	return types.SpecMapping{
		Module: "pkg/channels", AESPSpecs: []string{"AESP-0014", "AESP-0015"},
		Invariants: []string{"INV-02", "INV-11"}, Status: "stubbed",
	}
}
