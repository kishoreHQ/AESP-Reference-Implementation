package conformance

import "testing"

func TestCatalogNoSilentOmissions(t *testing.T) {
	cat := Catalog()
	if len(cat) < 10 {
		t.Fatal("catalog too small")
	}
	for _, m := range cat {
		switch m.Status {
		case "implemented", "stubbed", "missing", "gap-filed":
		default:
			t.Fatalf("invalid status %s for %s", m.Status, m.ID)
		}
		if m.Status == "missing" && m.Module != "" {
			// missing may have empty module
		}
	}
}

func TestModuleMappings(t *testing.T) {
	ms := ModuleMappings()
	if len(ms) < 15 {
		t.Fatalf("expected many module mappings, got %d", len(ms))
	}
}

func TestReport(t *testing.T) {
	r := Report()
	if len(r) < 100 {
		t.Fatal("report short")
	}
	t.Log(r)
}
