Feature: SwarmForge CLI

  The swarmforge CLI replaces the bash startup scripts (swarmforge.sh,
  notify-agent.sh, swarmlog.sh) with a single Go binary. Running
  "swarmforge start" performs preflight checks, creates project
  directories, builds agent prompts, and launches a tmux session with
  three AI agents and a metrics dashboard.

  Scenario: Preflight rejects missing dependency
    Given the system does not have "tmux" installed
    When the user runs preflight checks
    Then an error is returned containing "tmux"

  Scenario: Preflight passes with all dependencies
    Given the system has "tmux", "claude", and "watch" installed
    When the user runs preflight checks
    Then no error is returned

  Scenario: Directory setup creates required directories
    Given a project root directory exists
    When directory setup runs for the project root
    Then the directory "features" exists under the project root
    And the directory "logs" exists under the project root
    And the directory "agent_context" exists under the project root

  Scenario: Helper scripts are generated for backward compatibility
    Given a project root directory exists
    When setup writes helper scripts to the project root
    Then "notify-agent.sh" exists under the project root
    And "swarmlog.sh" exists under the project root

  Scenario: Agent prompt includes role and constitution
    Given a constitution with content "Rule 1: TDD is mandatory"
    And the agent role is "Architect" with standard instructions
    When the prompt builder generates the prompt
    Then the prompt contains "You are the Architect agent"
    And the prompt contains "Rule 1: TDD is mandatory"
    And the prompt contains "Pane 0 = Architect"

  Scenario: Agent prompt includes coordination instructions
    Given a constitution with content "Constitution content"
    And the agent role is "Coder" with standard instructions
    When the prompt builder generates the prompt
    Then the prompt contains "notify-agent.sh"
    And the prompt contains "swarmlog.sh"
    And the prompt contains "agent_context/"

  Scenario: Start kills existing tmux session before creating new one
    Given a tmux session named "swarmforge" already exists
    When the start sequence runs
    Then the existing "swarmforge" session is killed
    And a new "swarmforge" session is created

  Scenario: Start creates tmux session with 2x2 grid layout
    Given no tmux session named "swarmforge" exists
    When the start sequence creates the tmux session
    Then a new tmux session "swarmforge" with window "swarm" is created
    And the window is split into 4 panes
    And pane borders display agent titles

  Scenario: Agents are launched with correct claude commands
    Given a tmux session "swarmforge" with 4 panes exists
    And agent prompt files have been written
    When agents are launched in their panes
    Then pane 0 receives a claude command containing "SwarmForge Architect"
    And pane 1 receives a claude command containing "SwarmForge E2E-Interpreter"
    And pane 2 receives a claude command containing "SwarmForge Coder"
    And each claude command includes "--permission-mode acceptEdits"

  Scenario: Metrics pane tails the agent log file
    Given a tmux session "swarmforge" with 4 panes exists
    When the metrics pane is initialized
    Then pane 3 receives a command containing "tail -f logs/agent_messages.log"

  Scenario: Notify subcommand logs and sends message to pane
    Given a log writer is configured
    And a tmux commander is available
    When the user runs notify for pane 0 with message "hello architect"
    Then a timestamped log entry containing "[pane 0] hello architect" is written
    And tmux send-keys is invoked for session "swarmforge" pane 0

  Scenario: Log subcommand writes timestamped entry to file and stdout
    Given a log writer and stdout writer are configured
    When the user logs a message with role "Architect" and text "task started"
    Then the log writer contains "[Architect] task started"
    And the stdout writer contains "[Architect] task started"

  Scenario: CLI dispatches subcommands correctly
    Given the CLI receives arguments "start"
    Then the start handler is invoked
    Given the CLI receives arguments "notify" "1" "hello"
    Then the notify handler is invoked
    Given the CLI receives arguments "log" "Coder" "done"
    Then the log handler is invoked
    Given the CLI receives no arguments
    Then a usage error is returned

  Scenario: Full startup banner is displayed
    Given a writer captures output
    When the startup banner is printed
    Then the output contains "SwarmForge"
    And the output contains "Disciplined agents build better software"
