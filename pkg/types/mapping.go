package types

// SpecMapping documents which AESP requirements a module realizes.
// Used by conformance harness and stub tests.
type SpecMapping struct {
	Module       string   `json:"module"`
	AESPSpecs    []string `json:"aespSpecs"`
	RequirementIDs []string `json:"requirementIds,omitempty"`
	Invariants   []string `json:"invariants"`
	Status       string   `json:"status"` // implemented | stubbed | missing
}
