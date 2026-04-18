# SwarmForge

**A disciplined tmux-based agent orchestration platform that turns swarms of AI agents into reliable, professional software engineers.**

## Intent

SwarmForge exists to solve the core problem of agentic development: **chaos**.

Left unchecked, AI agents produce code quickly but often without discipline, leading to brittle, untested, hard-to-maintain software. SwarmForge changes that by embedding **strict professional craftsmanship** directly into the platform.

It enforces four foundational clean code disciplines — plus static linting — as an unbreakable **Constitution**. Every agent in the swarm must obey these rules on every task. The result is fast, scalable, and genuinely high-quality software produced reliably at swarm speed.

SwarmForge turns raw AI coding power into **disciplined, trustworthy engineering**.

## What SwarmForge Does

SwarmForge is a lightweight, tmux-based orchestration layer that:

- Launches a **config-driven swarm** from a project-local `swarmforge/swarmforge.conf`
- Creates one tmux session and one Terminal window per configured role
- Reads behavior from project-local `swarmforge/<role>.prompt` files plus `swarmforge/constitution.prompt`
- Supports per-role backends such as `claude`, `codex`, or `none`
- Creates a project-local `swarmtools/` directory with notification helpers for the active swarm
- Creates one git worktree per configured role under `.worktrees/`
- Initializes a git repository in a new working directory and creates a first commit with `logs/` and `agent_context/` ignored
- Keeps all swarm state local to the working directory in `.swarmforge/`

## Core Features

- **Config-Driven Topology** — The swarm shape comes from `swarmforge/swarmforge.conf`, not hardcoded shell variables.
- **Project-Local Roles** — Each role is defined by `swarmforge/<role>.prompt` in the working tree being orchestrated.
- **Backend Selection Per Role** — A role can launch `claude`, `codex`, or no agent at all.
- **Observable Swarm** — Open one Terminal window per role and watch the sessions in real time.
- **Self-Hosted & Lightweight** — Runs locally in tmux and Terminal with minimal machinery.

## How It Works (High Level)

1. Create a `swarmforge/` directory in the target working directory.
2. Put `swarmforge.conf`, `constitution.prompt`, and one `<role>.prompt` file per configured role inside it.
3. Define each window as `window <role> <agent> <worktree>`.
4. Add `swarmforge.sh` to your shell `PATH` before startup.
5. Run `swarmforge.sh <working-directory>` or run it from inside that directory.
6. If the working directory is not already a git repo, startup runs `git init`, renames the initial branch to `master`, writes `.gitignore` entries for `.swarmforge/`, `.worktrees/`, `swarmtools/`, `logs/`, and `agent_context/`, and makes the first commit from the current project state.
7. Startup creates a git worktree for each window under `.worktrees/<worktree>`, unless the worktree field is `none` or `master`.
8. Startup creates `swarmtools/notify-agent.sh` for that project.
9. SwarmForge creates tmux sessions, opens Terminal windows, and launches each configured backend in its assigned worktree.
10. Roles communicate through helper commands such as `notify-agent.sh` and `swarmlog.sh`.

Example config:

```conf
window architect claude architect
window coder codex coder
window e2e codex e2e
window logger none none
```

`logger` is a utility role. When configured with `none`, it tails `logs/agent_messages.log`.

In the example above, the agents run in these worktrees:

- `architect` -> `.worktrees/architect`
- `coder` -> `.worktrees/coder`
- `e2e` -> `.worktrees/e2e`
- `logger` -> main working directory

If a window uses `master` as its worktree name, SwarmForge does not create `.worktrees/master`; that role runs in the main working directory on the `master` branch.

The launcher expects these helper scripts to exist beside `swarmforge.sh`:

- `swarmlog.sh`
- `swarm-cleanup.sh`

## Who Is SwarmForge For?

- Developers who want to harness AI agents without sacrificing code quality
- Teams exploring agentic development practices
- Anyone tired of “AI wrote it” meaning “now I have to rewrite it”
- Clean Code enthusiasts who believe discipline still matters in the age of agents

## Getting Started

```bash
git clone https://github.com/unclebob/swarmforge.git
cd swarmforge
chmod +x swarmforge.sh
export PATH="/path/to/swarmforge:$PATH"
mkdir my-project
cd my-project
mkdir swarmforge
cat > swarmforge/swarmforge.conf <<'EOF'
window architect claude architect
window coder codex coder
window e2e codex e2e
window logger none none
EOF
cat > swarmforge/constitution.prompt <<'EOF'
Read this constitution and obey it on every task.
EOF
cat > swarmforge/architect.prompt <<'EOF'
You are the architect. Read swarmforge/constitution.prompt and follow it.
EOF
cat > swarmforge/coder.prompt <<'EOF'
You are the coder. Read swarmforge/constitution.prompt and follow it.
EOF
cat > swarmforge/e2e.prompt <<'EOF'
You are the e2e role. Read swarmforge/constitution.prompt and follow it.
EOF
swarmforge.sh .

# After startup, the project-local notification helper is available at:
# ./swarmtools/notify-agent.sh
```
