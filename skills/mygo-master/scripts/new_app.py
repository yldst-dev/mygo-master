#!/usr/bin/env python3
import argparse
import os
import shutil
import subprocess
import sys

SKILL = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
STARTER = os.path.join(SKILL, "assets", "starter")
CLI = "github.com/egoist/mygo/cmd/mygo@latest"


def run(command, directory):
    print("$", " ".join(command))
    done = subprocess.run(command, cwd=directory, capture_output=True, text=True)
    if done.stdout.strip():
        print(done.stdout.strip())
    if done.returncode != 0:
        print(done.stderr.strip())
        sys.exit(f"failed: {' '.join(command)}")
    return done.stdout


def module_of(directory):
    with open(os.path.join(directory, "go.mod")) as f:
        for line in f:
            if line.startswith("module "):
                return line.split()[1]
    sys.exit("go.mod names no module")


def copy_starter(directory, module):
    for root, _, names in os.walk(STARTER):
        for name in names:
            source = os.path.join(root, name)
            target = os.path.join(directory, os.path.relpath(source, STARTER))
            os.makedirs(os.path.dirname(target), exist_ok=True)
            with open(source) as f:
                text = f.read()
            with open(target, "w") as f:
                f.write(text.replace('"starter/', f'"{module}/'))
    with open(os.path.join(SKILL, "assets", "motion.go")) as f:
        motion = f.read().replace("package main", "package view", 1)
    with open(os.path.join(directory, "internal", "view", "motion.go"), "w") as f:
        f.write(motion)


def main():
    parser = argparse.ArgumentParser(description="Start a MyGo app on the latest release, in clean architecture.")
    parser.add_argument("directory")
    parser.add_argument("--name")
    parser.add_argument("--module")
    args = parser.parse_args()
    if sys.platform == "darwin" and args.name and not args.name.isascii():
        sys.exit("use an ASCII --name: codesign crashes on such an executable name; the window title can be in any language")
    directory = os.path.abspath(args.directory)
    parent = os.path.dirname(directory)
    os.makedirs(parent, exist_ok=True)
    command = ["go", "run", CLI, "init", "-template", "native"]
    if args.name:
        command += ["-name", args.name]
    if args.module:
        command += ["-module", args.module]
    run(command + [directory], parent)
    for name in ("main.go", "main_test.go"):
        os.remove(os.path.join(directory, name))
    copy_starter(directory, module_of(directory))
    run(["go", "mod", "tidy"], directory)
    unformatted = run(["gofmt", "-l", "."], directory)
    if unformatted.strip():
        sys.exit("gofmt reports the files above")
    run(["go", "vet", "./..."], directory)
    run(["go", "tool", "mygo", "vet", "."], directory)
    outdated = run(["go", "tool", "mygo", "migrate-ui", "."], directory)
    if outdated.strip():
        sys.exit("migrate-ui would change the files above")
    environment_test = ["env", "CGO_ENABLED=0", "go", "test", "-count=1", "-timeout", "120s", "./..."]
    run(environment_test, directory)
    version = run(["go", "list", "-m", "github.com/egoist/mygo"], directory).strip()
    print(f"ready: {directory} on {version}")


main()
