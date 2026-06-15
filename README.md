# SwarmForge

SwarmForge starts a configured set of agent roles in separate git worktrees and
tmux sessions.

## Handoff Protocol

SwarmForge uses a daemon-backed file handoff protocol instead of letting agents
send tmux messages directly. The launcher starts `handoffd.bb`, which owns the
tmux socket access, watches each agent worktree for validated outbound handoff
files, copies them into recipient inboxes, and sends only generic wake-up
notifications through tmux.

Agents interact with handoffs through helper scripts:

- `swarm_handoff.sh <draft-file>` validates and queues outbound handoffs.
- `ready_for_next.sh` accepts work using the role's configured receive mode.
- `done_with_current.sh` completes the current task or batch using the role's
  configured receive mode.

The durable handoff files and their audit timestamps replace the old logbook and
direct tmux handoff strategy. See [swarmforge/handoff-protocol.md](swarmforge/handoff-protocol.md)
for the full protocol.
