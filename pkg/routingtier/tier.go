package routingtier

import "github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"

// Tier is a cost class for capability routing (ADT-01). Not a vendor name.
type Tier string

const (
	FreeLocal  Tier = "free-local"
	FreeHosted Tier = "free-hosted"
	Budget     Tier = "budget"
	Standard   Tier = "standard"
	Premium    Tier = "premium"
)

// Order free → premium for free-first policies.
var Ascending = []Tier{FreeLocal, FreeHosted, Budget, Standard, Premium}

func Rank(t Tier) int {
	for i, x := range Ascending {
		if x == t {
			return i
		}
	}
	return len(Ascending)
}

// InferFromDescriptor maps provider flags to a default tier (no vendor strings).
func InferFromDescriptor(local bool, priority int) Tier {
	if local {
		return FreeLocal
	}
	if priority >= 100 {
		return Premium
	}
	if priority >= 50 {
		return Standard
	}
	return Budget
}

// FreeFirstCoding is an example policy name (registry/policy object, not hardcode routing).
const PolicyFreeFirstCoding = "free-first-coding"

func SpecMapping() types.SpecMapping {
	return types.SpecMapping{
		Module: "pkg/routingtier", AESPSpecs: []string{"AESP-0015"},
		Invariants: []string{"INV-03"}, Status: "stubbed",
	}
}
