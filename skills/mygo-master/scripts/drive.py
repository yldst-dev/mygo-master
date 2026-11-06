#!/usr/bin/env python3
import argparse
import base64
import json
import os
import socket
import subprocess
import sys
import tempfile
import time

PLUGIN = 'package main\n\nimport _ "github.com/egoist/mygo/plugins/automation"\n'
SOCKET = "e2e.sock"


class Client:
    def __init__(self, path, wait=30):
        deadline = time.monotonic() + wait
        while True:
            try:
                self.sock = socket.socket(socket.AF_UNIX)
                self.sock.connect(path)
                break
            except OSError:
                if time.monotonic() > deadline:
                    raise
                time.sleep(0.2)
        self.stream = self.sock.makefile("rwb")
        self.next_id = 0
        while True:
            try:
                self.call("idle", timeout=5000)
                break
            except RuntimeError:
                if time.monotonic() > deadline:
                    raise
                time.sleep(0.2)

    def call(self, method, **params):
        self.next_id += 1
        message = {"id": self.next_id, "method": method}
        if params:
            message["params"] = params
        self.stream.write(json.dumps(message).encode() + b"\n")
        self.stream.flush()
        reply = json.loads(self.stream.readline())
        if "error" in reply:
            raise RuntimeError(f"{method}: {reply['error']}")
        return reply["result"]

    def settle(self, pause=0.4):
        time.sleep(pause)
        self.call("idle", timeout=5000)

    def nodes(self):
        return self.call("tree")["nodes"]

    def box(self, label):
        for node in self.nodes():
            if node.get("label") == label:
                return node["bounds"]
        raise RuntimeError(f"no element labeled {label!r}")

    def move(self, x, y):
        self.call("input", events=[{"type": "pointerMove", "x": x, "y": y}])

    def hover(self, label, fx=0.5, fy=0.5):
        b = self.box(label)
        self.move(b["x"] + b["width"] * fx, b["y"] + b["height"] * fy)

    def click(self, label, fx=0.5, fy=0.5):
        b = self.box(label)
        self.call("input", events=[{"type": "click", "x": b["x"] + b["width"] * fx, "y": b["y"] + b["height"] * fy}])

    def key(self, name):
        self.call("input", events=[{"type": "key", "key": name}])

    def text(self, value):
        self.call("input", events=[{"type": "text", "text": value}])

    def role(self, role):
        return [node for node in self.nodes() if node.get("role") == role]

    def shot(self, path):
        with open(path, "wb") as f:
            f.write(base64.b64decode(self.call("capture")["png"]))


def build(project, binary):
    project, binary = os.path.abspath(project), os.path.abspath(binary)
    with tempfile.TemporaryDirectory() as work:
        plugin = os.path.join(work, "plugin.go")
        with open(plugin, "w") as f:
            f.write(PLUGIN)
        overlay = os.path.join(work, "overlay.json")
        with open(overlay, "w") as f:
            json.dump({"Replace": {os.path.join(project, "zz_e2e.go"): plugin}}, f)
        subprocess.run(["go", "build", "-overlay", overlay, "-o", binary, "."], cwd=project, check=True)
    print(binary)


def start(binary, workdir):
    binary = os.path.abspath(binary)
    workdir = os.path.abspath(workdir)
    os.makedirs(workdir, exist_ok=True)
    environment = dict(os.environ, MYGO_AUTOMATION=SOCKET, MYGO_FRAME_STATS="all")
    log = open(os.path.join(workdir, "e2e.log"), "wb")
    process = subprocess.Popen([binary], cwd=workdir, env=environment, stdout=log, stderr=log, start_new_session=True)
    print(f"pid {process.pid}")
    print(f"socket {os.path.join(workdir, SOCKET)}")
    print(f"log {log.name}")


def tree(path):
    os.chdir(os.path.dirname(os.path.abspath(path)))
    for node in Client(os.path.basename(path)).nodes():
        b = node.get("bounds", {})
        place = f"{b.get('x', 0):.0f},{b.get('y', 0):.0f} {b.get('width', 0):.0f}x{b.get('height', 0):.0f}"
        print(node.get("role"), repr(node.get("label")), repr(node.get("value")), place, node.get("states", []))


def shot(path, out):
    out = os.path.abspath(out)
    os.chdir(os.path.dirname(os.path.abspath(path)))
    Client(os.path.basename(path)).shot(out)
    print(out)


def main():
    parser = argparse.ArgumentParser(description="Build, start and read a MyGo app through its automation socket.")
    commands = parser.add_subparsers(dest="command", required=True)
    command = commands.add_parser("build", help="build a test binary with the automation plugin; the repository stays untouched")
    command.add_argument("project")
    command.add_argument("binary")
    command = commands.add_parser("start", help="start it with a short socket path and every frame logged")
    command.add_argument("binary")
    command.add_argument("workdir")
    command = commands.add_parser("tree", help="print what assistive technology sees")
    command.add_argument("socket")
    command = commands.add_parser("shot", help="capture the window as a PNG")
    command.add_argument("socket")
    command.add_argument("out")
    args = parser.parse_args()
    if args.command == "build":
        build(args.project, args.binary)
    elif args.command == "start":
        start(args.binary, args.workdir)
    elif args.command == "tree":
        tree(args.socket)
    else:
        shot(args.socket, args.out)


if __name__ == "__main__":
    main()
