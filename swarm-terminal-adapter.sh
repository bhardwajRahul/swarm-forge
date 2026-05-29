#!/usr/bin/env zsh

normalize_terminal_backend() {
  local backend="${1:l}"

  case "$backend" in
    terminal|terminal-app|terminal.app)
      echo "terminal-app"
      ;;
    ghostty)
      echo "ghostty"
      ;;
    none|current|fallback)
      echo "none"
      ;;
    *)
      echo "$backend"
      ;;
  esac
}

detect_terminal_backend() {
  if [[ -n "${SWARMFORGE_TERMINAL:-}" ]]; then
    normalize_terminal_backend "$SWARMFORGE_TERMINAL"
    return
  fi

  if has_command osascript; then
    echo "terminal-app"
    return
  fi

  echo "none"
}

terminal_backend_label() {
  case "$TERMINAL_BACKEND" in
    ghostty) echo "Ghostty" ;;
    terminal-app) echo "Terminal" ;;
    *) echo "current shell" ;;
  esac
}

terminal_backend_tracks_windows() {
  case "$TERMINAL_BACKEND" in
    ghostty|terminal-app) return 0 ;;
    *) return 1 ;;
  esac
}

terminal_window_exists() {
  local window_id="$1"
  [[ -n "$window_id" ]] || return 1

  case "$TERMINAL_BACKEND" in
    ghostty)
      local result
      result="$(osascript - "$window_id" <<'APPLESCRIPT' 2>/dev/null || true
on run argv
  set targetId to item 1 of argv
  tell application "Ghostty"
    repeat with w in windows
      repeat with t in tabs of w
        if (id of t as string) is targetId then return "yes"
      end repeat
    end repeat
  end tell
  return "no"
end run
APPLESCRIPT
)"
      [[ "$result" == "yes" ]]
      ;;
    terminal-app)
      local result
      result="$(osascript - "$window_id" <<'APPLESCRIPT' 2>/dev/null || true
on run argv
  set targetId to item 1 of argv as integer
  tell application "Terminal"
    repeat with terminalWindow in windows
      if id of terminalWindow is targetId then return "yes"
    end repeat
  end tell
  return "no"
end run
APPLESCRIPT
)"
      [[ "$result" == "yes" ]]
      ;;
    *)
      return 1
      ;;
  esac
}

terminal_open_session() {
  local session="$1"
  local title="$2"
  local sibling_id="${3:-}"

  case "$TERMINAL_BACKEND" in
    ghostty)
      osascript - "$WORKING_DIR" "$session" "$title" "$TMUX_SOCKET" "$sibling_id" <<'APPLESCRIPT'
on run argv
  set workingDir to item 1 of argv
  set tmuxSession to item 2 of argv
  set tmuxSocket to item 4 of argv
  set siblingTabId to item 5 of argv
  set initialCmd to "cd " & quoted form of workingDir & " && exec tmux -S " & quoted form of tmuxSocket & " attach-session -t " & quoted form of tmuxSession & linefeed

  tell application "Ghostty"
    set cfg to new surface configuration
    set initial working directory of cfg to workingDir
    set initial input of cfg to initialCmd

    if siblingTabId is not "" then
      set targetWin to missing value
      set siblingTab to missing value
      repeat with w in windows
        repeat with t in tabs of w
          if (id of t as string) is siblingTabId then
            set targetWin to w
            set siblingTab to t
            exit repeat
          end if
        end repeat
        if targetWin is not missing value then exit repeat
      end repeat
      if targetWin is not missing value then
        select tab siblingTab
        set newTab to new tab in targetWin with configuration cfg
        return id of newTab
      end if
    end if

    try
      set targetWin to front window
      set newTab to new tab in targetWin with configuration cfg
      return id of newTab
    end try

    set newWin to new window with configuration cfg
    return id of (first tab of newWin)
  end tell
end run
APPLESCRIPT
      ;;
    terminal-app)
      osascript - "$WORKING_DIR" "$session" "$title" "$TMUX_SOCKET" <<'APPLESCRIPT'
on run argv
  set workingDir to item 1 of argv
  set tmuxSession to item 2 of argv
  set windowTitle to item 3 of argv
  set tmuxSocket to item 4 of argv

  tell application "Terminal"
    activate
    set newTab to do script ""
    do script "cd " & quoted form of workingDir & " && exec tmux -S " & quoted form of tmuxSocket & " attach-session -t " & quoted form of tmuxSession in newTab
    set custom title of newTab to windowTitle
    return id of front window
  end tell
end run
APPLESCRIPT
      ;;
    *)
      return 1
      ;;
  esac
}

terminal_close_window() {
  local window_id="$1"
  [[ -n "$window_id" ]] || return 0

  case "$TERMINAL_BACKEND" in
    ghostty)
      osascript - "$window_id" <<'APPLESCRIPT' >/dev/null 2>&1 || true
on run argv
  set targetId to item 1 of argv
  tell application "Ghostty"
    try
      repeat with w in windows
        repeat with t in tabs of w
          if (id of t as string) is targetId then
            close tab t
            return
          end if
        end repeat
      end repeat
    end try
  end tell
end run
APPLESCRIPT
      ;;
    terminal-app)
      osascript - "$window_id" <<'APPLESCRIPT' >/dev/null 2>&1 || true
on run argv
  set targetId to item 1 of argv as integer
  tell application "Terminal"
    try
      close (first window whose id is targetId) saving no
    end try
  end tell
end run
APPLESCRIPT
      ;;
  esac
}
