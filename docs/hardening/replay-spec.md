# Deterministic Replay Specification

## Scope

Replay reconstructs **control decisions** and substitutes recorded compute outputs.
Live model re-sampling is opt-in and non-deterministic.

## Event types covered

All `aesp.*` categories from event-model.md:

control · runtime · tool · provider · memory · artifact · hitl · obs · rem

## Inputs required

1. Mission journal (ordered events)
2. Pinned artifact digests + blobs
3. Recorded provider completion payloads (or hashes)
4. Recorded tool results + trust labels
5. Policy version id effective at time T
6. Plan artifact revisions

## Algorithm

```
for event in journal:
  re-apply control transitions
  if event is compute (provider/tool):
    inject recorded output unless live=true
  rebuild execution tree
assert tree.digest == recorded.tree.digest  # structural
```

## Non-goals

- Bit-identical GPU sampling  
- Replaying host UI interactions beyond ResolveApproval records  
