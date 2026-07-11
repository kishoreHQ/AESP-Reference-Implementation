package agentmode

import "github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"

// Mode controls autonomy (ADT-07).
type Mode string

const (
	Full    Mode = "full"    // execute under policy
	Assist  Mode = "assist"  // external actions → approval
	Observe Mode = "observe" // journal only, no execute
)

func Valid(m Mode) bool {
	return m == Full || m == Assist || m == Observe
}

// ExternalActionAllowed reports whether side effects may run.
func ExternalActionAllowed(m Mode) (allowed bool, requireApproval bool) {
	switch m {
	case Full:
		return true, false
	case Assist:
		return true, true
	case Observe:
		return false, false
	default:
		return false, false
	}
}

func SpecMapping() types.SpecMapping {
	return types.SpecMapping{
		Module: "pkg/agentmode", AESPSpecs: []string{"AESP-0002", "AESP-0014"},
		Invariants: []string{"INV-10"}, Status: "stubbed",
	}
}
