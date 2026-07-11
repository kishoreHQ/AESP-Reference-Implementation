package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/conformance"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/contextenv"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/eventbus"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/kernel"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/orchestrator"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/providerregistry"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/router"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/runtimeregistry"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
	mocklocal "github.com/kishoreHQ/AESP-Reference-Implementation/plugins/providers/mock-local"
	mockremote "github.com/kishoreHQ/AESP-Reference-Implementation/plugins/providers/mock-remote"
	genericloop "github.com/kishoreHQ/AESP-Reference-Implementation/plugins/runtimes/generic-loop"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "conformance" {
		fmt.Print(conformance.Report())
		return
	}

	bus := eventbus.NewMemoryBus()
	k := kernel.New(bus)
	pr := providerregistry.New()
	rr := runtimeregistry.New()
	_ = pr.Register(mocklocal.New())
	_ = pr.Register(mockremote.New())
	_ = rr.Register(genericloop.New())
	orch := orchestrator.New(bus, router.New(pr, rr))

	ctx := context.Background()
	m := types.Mission{
		ID:              "wu_demo",
		Tenant:          "local",
		Goal:            "prove host-neutral kernel loop",
		RequiredCaps:    []types.Capability{"coding", "local"},
		Budget:          types.Budget{MaxSteps: 5, MaxTokens: 1000, MaxWallSec: 30},
		SuccessCriteria: []string{"completed"},
		CreatedAt:       time.Now().UTC(),
	}
	id, err := k.SubmitMission(ctx, m)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	res, err := orch.RunMission(ctx, m, contextenv.Envelope{
		Prompt:   "execute mission",
		Security: contextenv.SecurityContext{Tenant: m.Tenant, Principal: "agent.demo"},
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("mission %s status=%s output=%s\n", id, res.Status, res.Output)
	tree, _ := k.GetExecutionTree(ctx, id)
	fmt.Printf("execution tree workUnit=%s\n", tree.WorkUnitID)
}
