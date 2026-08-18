#!/usr/bin/env python3
"""Run a TUI binary in a virtual terminal (via pyte) and print its rendered
screen as plain text, for visual verification of bubbletea/lipgloss output.

Dev tooling only, provided by the nix-shell. Not part of the Go module.

Usage:
  tui-snapshot.py [--width W] [--height H] [--wait SECONDS] [--keys 'k:delay,...'] -- <cmd> [args...]

Example:
  tui-snapshot.py --keys 'l:0.2,l:0.2,q:0' -- go run ./examples/tabs
"""

import argparse
import fcntl
import os
import pty
import select
import struct
import sys
import termios
import time

import pyte

class RepeatAwareScreen(pyte.Screen):
    """pyte has no support for the ECMA-48 REP sequence (``CSI Ps b``,
    "repeat the preceding graphic character Ps times"), which bubbletea's
    renderer uses to compress runs of identical styled cells (e.g. long
    border/padding runs). Left unhandled, pyte silently drops it, corrupting
    the frame. This adds the missing handler."""

    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self._last_drawn = " "

    def draw(self, data):
        super().draw(data)
        if data:
            self._last_drawn = data[-1]

    def repeat_last_character(self, count=1):
        self.draw(self._last_drawn * (count or 1))


class RepeatAwareStream(pyte.Stream):
    """Also works around a pyte bug: it maps ECMA-48 HPA (Character Position
    Absolute, "move cursor to column Ps") to "'" (apostrophe, 0x27) instead
    of the actual standard character "`" (backtick, 0x60) - see pyte's
    escape.py. Real terminal apps send the correct backtick, which pyte then
    silently no-ops on since it's not in its dispatch table, leaving the
    cursor stuck instead of moving it. bubbletea's renderer uses HPA (mixed
    with CUF) to jump the cursor to a border's column after erasing padding
    with ECH, so without this the trailing border character gets drawn right
    after the leading one instead of at the far edge."""

    csi = {
        **pyte.Stream.csi,
        "b": "repeat_last_character",
        "`": "cursor_to_column",
    }


KEYMAP = {
    "enter": "\r",
    "esc": "\x1b",
    "tab": "\t",
    "up": "\x1b[A",
    "down": "\x1b[B",
    "right": "\x1b[C",
    "left": "\x1b[D",
    "space": " ",
}


def encode_key(key):
    """Resolve a --keys token to the bytes to write to the pty. Falls back to
    the literal string for a plain character (e.g. 'l'), but without this,
    something like 'ctrl+l' would be sent as the 6 literal characters
    'c','t','r','l','+','l' instead of the single 0x0C control byte a real
    terminal would send for Ctrl-L - silently testing the wrong thing."""
    if key in KEYMAP:
        return KEYMAP[key]
    if key.startswith("ctrl+") and len(key) == 6 and key[5].isalpha():
        return chr(ord(key[5].lower()) - ord("a") + 1)
    return key


def set_size(fd, rows, cols):
    fcntl.ioctl(fd, termios.TIOCSWINSZ, struct.pack("HHHH", rows, cols, 0, 0))


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--width", type=int, default=80)
    parser.add_argument("--height", type=int, default=24)
    parser.add_argument(
        "--wait", type=float, default=1.0, help="seconds to let the app render before the snapshot"
    )
    parser.add_argument(
        "--keys",
        default="",
        help="comma-separated key:delay pairs sent before the final snapshot, e.g. 'l:0.2,l:0.2,q:0'",
    )
    parser.add_argument(
        "--colors",
        action="store_true",
        help="also print, per row, the foreground color of each non-blank cell (for border/color bugs "
        "that don't show up in the plain-text dump)",
    )
    parser.add_argument("cmd", nargs=argparse.REMAINDER)
    args = parser.parse_args()

    cmd = args.cmd
    if cmd and cmd[0] == "--":
        cmd = cmd[1:]
    if not cmd:
        parser.error("missing command to run")

    screen = RepeatAwareScreen(args.width, args.height)
    stream = RepeatAwareStream(screen)

    # Open the pty and set its size *before* forking, so the child never
    # observes a stale/default size on its first render (pyte has no way to
    # recover from a corrupted initial frame drawn at the wrong width).
    master_fd, slave_fd = pty.openpty()
    set_size(slave_fd, args.height, args.width)

    pid = os.fork()
    if pid == 0:
        os.close(master_fd)
        os.setsid()
        fcntl.ioctl(slave_fd, termios.TIOCSCTTY, 0)
        os.dup2(slave_fd, 0)
        os.dup2(slave_fd, 1)
        os.dup2(slave_fd, 2)
        os.close(slave_fd)
        os.execvp(cmd[0], cmd)
        os._exit(1)

    os.close(slave_fd)
    fd = master_fd

    def pump(duration):
        end = time.time() + duration
        while True:
            remaining = end - time.time()
            if remaining <= 0:
                break
            r, _, _ = select.select([fd], [], [], remaining)
            if fd not in r:
                continue
            try:
                data = os.read(fd, 65536)
            except OSError:
                break
            if not data:
                break
            stream.feed(data.decode(errors="ignore"))

    pump(args.wait)

    for pair in filter(None, args.keys.split(",")):
        key, _, delay = pair.partition(":")
        try:
            os.write(fd, encode_key(key).encode())
        except OSError:
            break
        pump(float(delay) if delay else 0.3)

    try:
        os.kill(pid, 15)
        os.waitpid(pid, 0)
    except (ProcessLookupError, ChildProcessError):
        pass

    for y, line in enumerate(screen.display):
        print(line.rstrip())
        if args.colors:
            row = screen.buffer[y]
            cells = []
            for x in sorted(row):
                ch = row[x]
                if ch.data.strip():
                    fg = ch.fg if ch.fg != "default" else "-"
                    cells.append(f"{x}:{ch.data!r}:{fg}")
            if cells:
                print("  " + " ".join(cells))


if __name__ == "__main__":
    main()
