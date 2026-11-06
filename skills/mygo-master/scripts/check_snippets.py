#!/usr/bin/env python3
import argparse
import glob
import os
import re
import shutil
import subprocess
import sys
import tempfile

SKILL = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
FENCE = re.compile(r"^```go( test)?\n(.*?)^```", re.M | re.S)
IMPORTS = {
    "bufio.": "bufio",
    "bytes.": "bytes",
    "color.": "image/color",
    "fmt.": "fmt",
    "image.": "image",
    "log.": "log",
    "math.": "math",
    "mygo.": "github.com/egoist/mygo",
    "os.": "os",
    "filepath.": "path/filepath",
    "reflect.": "reflect",
    "sort.": "sort",
    "strings.": "strings",
    "testing.": "testing",
    "time.": "time",
    "ui.": "github.com/egoist/mygo/ui",
    "chart.": "github.com/egoist/mygo/ui/chart",
}
TOKEN = re.compile(r"`([^`\n]+)`")
PACKAGE_NAME = re.compile(r"\b(ui|chart|mygo)\.([A-Z]\w*)")
METHOD_NAME = re.compile(r"(?:^|[\s(])(?:[a-z]\w*)?\.([A-Z]\w*)\(?")
FOREIGN = {"Sleep", "BROKEN", "MISSED", "Client", "RGBAAt", "ReadFile", "TempDir", "Dir", "Tasks", "TaskStore", "Fatal"}
PACKAGES = {
    "ui": "github.com/egoist/mygo/ui",
    "chart": "github.com/egoist/mygo/ui/chart",
    "mygo": "github.com/egoist/mygo",
}


def documents():
    return [os.path.join(SKILL, "SKILL.md")] + sorted(glob.glob(os.path.join(SKILL, "references", "*.md")))


def header(code):
    used = sorted(path for prefix, path in IMPORTS.items() if re.search(r"(?<![\w.])" + re.escape(prefix), code))
    return "package main\n\nimport (\n" + "".join(f'\t"{path}"\n' for path in used) + ")\n\n"


def collect():
    plain, tests = [], []
    for path in documents():
        with open(path) as f:
            for tag, code in FENCE.findall(f.read()):
                (tests if tag else plain).append(code)
    return "\n".join(plain), "\n".join(tests), len(plain) + len(tests)


def go_environment(pinned):
    return dict(os.environ, GOFLAGS="-mod=mod", GOPROXY="off") if pinned else dict(os.environ)


def compile_and_test(project, work):
    environment = go_environment(project)
    if project:
        for name in ("go.mod", "go.sum"):
            shutil.copy(os.path.join(project, name), work)
    else:
        subprocess.run(["go", "mod", "init", "snippets"], cwd=work, capture_output=True, check=True)
        subprocess.run(["go", "get", "github.com/egoist/mygo@latest"], cwd=work, capture_output=True, check=True)
    shutil.copy(os.path.join(SKILL, "assets", "motion.go"), work)
    shutil.copy(os.path.join(SKILL, "assets", "film_test.go"), os.path.join(work, "film_test.go"))
    plain, tests, count = collect()
    with open(os.path.join(work, "snippets.go"), "w") as f:
        f.write(header(plain) + plain + ("" if "func main()" in plain else "\nfunc main() {}\n"))
    with open(os.path.join(work, "snippets_test.go"), "w") as f:
        f.write(header(tests) + tests)
    if not project:
        subprocess.run(["go", "mod", "tidy"], cwd=work, capture_output=True, check=True)
    version = subprocess.run(["go", "list", "-m", "github.com/egoist/mygo"], cwd=work, env=environment, capture_output=True, text=True)
    print("checking against", version.stdout.strip())
    ok = True
    for command in (["go", "vet", "."], ["go", "test", "-count=1", "-timeout", "120s", "."]):
        done = subprocess.run(command, cwd=work, env=environment, capture_output=True, text=True)
        print(" ".join(command), "ok" if done.returncode == 0 else "FAILED")
        if done.returncode != 0:
            print(done.stdout + done.stderr)
            ok = False
    module = version.stdout.split()[-1]
    outdated = subprocess.run(
        ["go", "run", f"github.com/egoist/mygo/cmd/mygo@{module}", "migrate-ui", "."], cwd=work, capture_output=True, text=True
    )
    print("migrate-ui", "ok" if outdated.returncode == 0 and not outdated.stdout.strip() else "FAILED")
    if outdated.returncode != 0 or outdated.stdout.strip():
        print(outdated.stdout + outdated.stderr)
        ok = False
    print(f"{count} code blocks")
    return ok


def check_starter():
    with tempfile.TemporaryDirectory() as work:
        done = subprocess.run(
            [sys.executable, os.path.join(SKILL, "scripts", "new_app.py"), os.path.join(work, "app")],
            capture_output=True,
            text=True,
        )
    print("new_app.py", "ok" if done.returncode == 0 else "FAILED")
    if done.returncode != 0:
        print(done.stdout + done.stderr)
    return done.returncode == 0


def source_text(project):
    done = subprocess.run(
        ["go", "list", "-m", "-f", "{{.Dir}}", "github.com/egoist/mygo"], cwd=project, capture_output=True, text=True, check=True
    )
    root = done.stdout.strip()
    text = ""
    for pattern in ("*.go", "ui/*.go", "ui/chart/*.go", "plugins/automation/*.go"):
        for path in glob.glob(os.path.join(root, pattern)):
            if not path.endswith("_test.go"):
                with open(path, errors="replace") as f:
                    text += f.read()
    return text


def exported(project, package, name):
    done = subprocess.run(["go", "doc", f"{PACKAGES[package]}.{name}"], cwd=project, capture_output=True, text=True)
    return done.returncode == 0


def check_names(project):
    text = source_text(project)
    missing = set()
    known = {}
    for path in documents():
        with open(path) as f:
            document = FENCE.sub("", f.read())
        where = os.path.basename(path)
        for token in TOKEN.findall(document):
            for package, name in PACKAGE_NAME.findall(token):
                if (package, name) not in known:
                    known[(package, name)] = exported(project, package, name)
                if not known[(package, name)]:
                    missing.add(f"{package}.{name} in {where}")
            if "/" in token or (" " in token and "(" not in token):
                continue
            stripped = PACKAGE_NAME.sub("", token)
            for name in METHOD_NAME.findall(stripped):
                if name not in FOREIGN and not re.search(r"func\s+\([^)]*\)\s*" + name + r"\b|^\s*(?:\w+,\s*)*" + name + r"[\s,]", text, re.M):
                    missing.add(f".{name} in {where}")
            if re.fullmatch(r"[A-Z]\w+(\(.*\))?", token):
                name = re.match(r"[A-Z]\w+", token).group(0)
                if name not in FOREIGN and not re.search(r"func\s+(\([^)]*\)\s*)?" + name + r"\b|\btype\s+" + name + r"\b|^\s*(?:\w+,\s*)*" + name + r"[\s,]", text, re.M):
                    missing.add(f"{name} in {where}")
    for item in sorted(missing):
        print("name not found in this version:", item)
    print(f"{len(known)} package names checked, {len(missing)} names missing")
    return not missing


def main():
    parser = argparse.ArgumentParser(description="Check this skill's code and API names against MyGo: the latest release, or a project's pin.")
    parser.add_argument("--project", help="a project whose go.mod pins the version to check against")
    parser.add_argument("--keep", help="a directory to keep the generated module in")
    args = parser.parse_args()
    project = os.path.abspath(args.project) if args.project else None
    temporary = None
    if args.keep:
        work = os.path.abspath(args.keep)
        os.makedirs(work, exist_ok=True)
    else:
        temporary = tempfile.TemporaryDirectory()
        work = temporary.name
    built = compile_and_test(project, work)
    named = check_names(project or work)
    started = project is not None or check_starter()
    if temporary:
        temporary.cleanup()
    if not (built and named and started):
        sys.exit(1)


main()
