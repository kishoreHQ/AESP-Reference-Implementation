package kernel

import (
	"context"
	"fmt"
	"sync"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/eventbus"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/host"
	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/types"
)

// Kernel is the host-neutral Agent OS core.
// Zero vendor names. Zero host product names.
type Kernel struct {
	mu     sync.Mutex
	bus    eventbus.Bus
	missions map[types.WorkUnitID]*types.Mission
	trees  map[types.WorkUnitID]*host.ExecutionTree
}

// New constructs a kernel with an in-memory event bus.
func New(bus eventbus.Bus) *Kernel {
	if bus == nil {
		bus = eventbus.NewMemoryBus()
	}
	return &Kernel{
		bus:      bus,
		missions: make(map[types.WorkUnitID]*types.Mission),
		trees:    make(map[types.WorkUnitID]*host.ExecutionTree),
	}
}

func (k *Kernel) SubmitMission(ctx context.Context, m types.Mission) (types.WorkUnitID, error) {
	if m.ID == "" {
		return "", fmt.Errorf("mission id required")
	}
	if len(m.RequiredCaps) == 0 {
		return "", fmt.Errorf("requiredCapabilities must be non-empty (INV-03)")
	}
	k.mu.Lock()
	defer k.mu.Unlock()
	cp := m
	k.missions[m.ID] = &cp
	k.trees[m.ID] = &host.ExecutionTree{WorkUnitID: m.ID}
	_ = k.bus.Publish(ctx, eventbus.Event{
		Type:       "aesp.control.mission.accepted",
		WorkUnitID: m.ID,
		Data:       map[string]any{"goal": m.Goal},
	})
	return m.ID, nil
}

func (k *Kernel) CancelMission(ctx context.Context, id types.WorkUnitID, reason string) error {
	k.mu.Lock()
	defer k.mu.Unlock()
	if _, ok := k.missions[id]; !ok {
		return fmt.Errorf("unknown mission %s", id)
	}
	return k.bus.Publish(ctx, eventbus.Event{
		Type:       "aesp.control.mission.cancelled",
		WorkUnitID: id,
		Data:       map[string]any{"reason": reason},
	})
}

func (k *Kernel) SubscribeEvents(ctx context.Context, id types.WorkUnitID) (<-chan host.Event, error) {
	ch, err := k.bus.Subscribe(ctx, string(id))
	if err != nil {
		return nil, err
	}
	out := make(chan host.Event, 16)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case e, ok := <-ch:
				if !ok {
					return
				}
				out <- host.Event{Type: e.Type, ID: e.ID, WorkUnitID: e.WorkUnitID, Data: e.Data}
			}
		}
	}()
	return out, nil
}

func (k *Kernel) ResolveApproval(ctx context.Context, taskID types.HITLTaskID, decision host.ApprovalDecision) error {
	return k.bus.Publish(ctx, eventbus.Event{
		Type: "aesp.hitl.approval.resolved",
		Data: map[string]any{"taskId": string(taskID), "approved": decision.Approved, "actor": decision.Actor},
	})
}

func (k *Kernel) GetArtifact(ctx context.Context, digest types.ArtifactDigest) ([]byte, error) {
	return nil, fmt.Errorf("artifact store stub: %s", digest)
}

func (k *Kernel) GetExecutionTree(ctx context.Context, id types.WorkUnitID) (*host.ExecutionTree, error) {
	k.mu.Lock()
	defer k.mu.Unlock()
	t, ok := k.trees[id]
	if !ok {
		return nil, fmt.Errorf("unknown mission %s", id)
	}
	cp := *t
	return &cp, nil
}

func (k *Kernel) Health(ctx context.Context) error { return nil }

// Ensure Kernel implements host.Interface.
var _ host.Interface = (*Kernel)(nil)

func SpecMapping() types.SpecMapping {
	return types.SpecMapping{
		Module:    "pkg/kernel",
		AESPSpecs: []string{"AESP-0001", "AESP-0003", "AESP-0005", "AESP-0013", "AESP-0015"},
		Invariants: []string{"INV-02", "INV-08", "INV-10", "INV-11"},
		Status:    "stubbed",
	}
}
