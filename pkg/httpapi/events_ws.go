package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/kishoreHQ/AESP-Reference-Implementation/pkg/eventbus"
	"nhooyr.io/websocket"
)

// mapKernelType converts aesp.* bus types to UI HostEvent types (docs/ui/CONTRACT.md).
func mapKernelType(t string) string {
	switch {
	case strings.Contains(t, "hitl.approval.requested"):
		return "approval.created"
	case strings.Contains(t, "hitl.approval.resolved"):
		return "approval.resolved"
	case strings.Contains(t, "artifact"):
		return "artifact.created"
	case strings.Contains(t, "memory"):
		return "memory.written"
	case strings.Contains(t, "tool") || strings.Contains(t, "provider.completed"):
		return "log.append"
	case strings.Contains(t, "runtime") || strings.Contains(t, "route"):
		return "node.updated"
	default:
		return "mission.updated"
	}
}

func (s *Server) apiEventsWS(w http.ResponseWriter, r *http.Request) {
	since, _ := strconv.ParseInt(r.URL.Query().Get("since"), 10, 64)
	mission := r.URL.Query().Get("mission")

	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
		OriginPatterns:     []string{"*"},
	})
	if err != nil {
		return
	}
	defer c.Close(websocket.StatusNormalClosure, "")

	ctx := r.Context()

	// Catch-up from global journal (UI-RT-01)
	if sinceBus, ok := s.sys.Bus.(interface {
		Since(ctx context.Context, since int64) ([]eventbus.Event, error)
	}); ok {
		evs, _ := sinceBus.Since(ctx, since)
		for _, e := range evs {
			if mission != "" && string(e.WorkUnitID) != mission {
				continue
			}
			if err := writeWSEvent(ctx, c, e); err != nil {
				return
			}
		}
	}

	// Live subscribe (empty filter = all work units)
	filter := mission
	ch, err := s.sys.Bus.Subscribe(ctx, filter)
	if err != nil {
		return
	}
	// Also subscribe to all if mission-scoped already handled; for all events use ""
	if mission != "" {
		// also get global? keep mission filter only
	} else {
		// already subscribed to ""
	}
	_ = since
	for {
		select {
		case <-ctx.Done():
			return
		case e, ok := <-ch:
			if !ok {
				return
			}
			if err := writeWSEvent(ctx, c, e); err != nil {
				return
			}
		}
	}
}

func writeWSEvent(ctx context.Context, c *websocket.Conn, e eventbus.Event) error {
	payload := map[string]any{
		"seq":       e.Seq,
		"type":      mapKernelType(e.Type),
		"ts":        e.Time.UTC().Format(time.RFC3339Nano),
		"missionId": string(e.WorkUnitID),
		"data": map[string]any{
			"rawType": e.Type,
			"payload": e.Data,
		},
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return c.Write(ctx, websocket.MessageText, b)
}
