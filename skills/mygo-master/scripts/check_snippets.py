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
    "bytes.": "bytes",
    "color.": "image/color",
    "fmt.": "fmt",
    "log.": "log",
    "mygo.": "github.com/egoist/mygo",
    "os.": "os",
    "filepath.": "path/filepath",
    "reflect.": "reflect",
    "strings.": "strings",
    "testing.": "testing",
    "time.": "time",
    "ui.": "github.com/egoist/mygo/ui",
    "chart.": "github.com/egoist/mygo/ui/chart",
}
HEX = re.compile(r"#[0-9A-Fa-f]{6}\b")
KIT = ("motion.go", "palette.go", "palette_test.go", "targets_test.go", "components.go", "components_test.go")
UNTRACKED = {"build", ".mygo", ".agents", ".DS_Store"}
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


def prepare(command, work):
    done = subprocess.run(command, cwd=work, capture_output=True, text=True)
    if done.returncode != 0:
        print(" ".join(command), "FAILED")
        print(done.stdout + done.stderr)
        sys.exit(1)


def compile_and_test(project, work):
    environment = go_environment(project)
    if project:
        for name in ("go.mod", "go.sum"):
            shutil.copy(os.path.join(project, name), work)
    else:
        prepare(["go", "mod", "init", "snippets"], work)
        prepare(["go", "get", "github.com/egoist/mygo@latest"], work)
    shutil.copy(os.path.join(SKILL, "assets", "motion.go"), work)
    shutil.copy(os.path.join(SKILL, "assets", "film_test.go"), os.path.join(work, "film_test.go"))
    plain, tests, count = collect()
    with open(os.path.join(work, "snippets.go"), "w") as f:
        f.write(header(plain) + plain + ("" if "func main()" in plain else "\nfunc main() {}\n"))
    with open(os.path.join(work, "snippets_test.go"), "w") as f:
        f.write(header(tests) + tests)
    if not project:
        prepare(["go", "mod", "tidy"], work)
    version = subprocess.run(["go", "list", "-m", "github.com/egoist/mygo"], cwd=work, env=environment, capture_output=True, text=True)
    if version.returncode != 0 or not version.stdout.split():
        print("go list -m github.com/egoist/mygo FAILED")
        print(version.stdout + version.stderr)
        return False
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
    ok = True
    for name, flags in (("new_app.py", []), ("new_app.py --kit dashboard", ["--kit", "dashboard"])):
        with tempfile.TemporaryDirectory() as work:
            done = subprocess.run(
                [sys.executable, os.path.join(SKILL, "scripts", "new_app.py"), os.path.join(work, "app")] + flags,
                capture_output=True,
                text=True,
            )
        print(name, "ok" if done.returncode == 0 else "FAILED")
        if done.returncode != 0:
            print(done.stdout + done.stderr)
            ok = False
    return ok


def check_example():
    with tempfile.TemporaryDirectory() as work:
        app = os.path.join(work, "dashboard")
        shutil.copytree(os.path.join(SKILL, "examples", "dashboard"), app)
        for command in (
            ["go", "get", "github.com/egoist/mygo@latest"],
            ["go", "mod", "tidy"],
            ["go", "vet", "."],
            ["go", "test", "-count=1", "-timeout", "120s", "."],
        ):
            done = subprocess.run(command, cwd=app, capture_output=True, text=True)
            if done.returncode != 0:
                print("examples/dashboard FAILED:", " ".join(command))
                print(done.stdout + done.stderr)
                return False
    print("examples/dashboard ok")
    return True


def read_tree(root):
    files = {}
    for directory, names, filenames in os.walk(root):
        names[:] = [name for name in names if name not in UNTRACKED]
        for name in filenames:
            if name not in UNTRACKED:
                path = os.path.join(directory, name)
                with open(path, "rb") as f:
                    files[os.path.relpath(path, root)] = f.read()
    return files


def check_kit():
    example = os.path.join(SKILL, "examples", "dashboard")
    problems = []
    for name in KIT:
        with open(os.path.join(SKILL, "assets", name)) as a, open(os.path.join(example, name)) as b:
            original, copy = a.read(), b.read()
        if name == "motion.go":
            kept = {line.strip() for line in original.splitlines()}
            for line in copy.splitlines():
                if line.strip() and line.strip() not in kept:
                    problems.append(f"examples/dashboard/motion.go has a line assets/motion.go does not: {line.strip()}")
        elif original != copy:
            problems.append(f"assets/{name} differs from examples/dashboard/{name}")
    with open(os.path.join(SKILL, "references", "design.md")) as f:
        colors = re.search(r"^## 색$(.*?)^## 글자$", f.read(), re.M | re.S).group(1)
    with open(os.path.join(example, "palette.go")) as f:
        palette = f.read().upper()
    for stale in sorted({code.upper() for code in HEX.findall(colors)} - {code for code in HEX.findall(palette)}):
        problems.append(f"design.md names {stale}, which palette.go does not have")
    repository = os.path.normpath(os.path.join(SKILL, "..", "..", "dashboard"))
    if os.path.isfile(os.path.join(repository, "mygo.json")):
        kept, original = read_tree(example), read_tree(repository)
        for name in sorted(set(kept) | set(original)):
            if kept.get(name) != original.get(name):
                problems.append(f"examples/dashboard/{name} differs from dashboard/{name}")
    for problem in problems:
        print(problem)
    print("kit and example", "ok" if not problems else "FAILED")
    return not problems


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
                if name not in FOREIGN and not re.match(r"Test[A-Z]", name) and not re.search(r"func\s+(\([^)]*\)\s*)?" + name + r"\b|\btype\s+" + name + r"\b|^\s*(?:\w+,\s*)*" + name + r"[\s,]", text, re.M):
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
    bundled = project is not None or all([check_starter(), check_example(), check_kit()])
    if temporary:
        temporary.cleanup()
    if not (built and named and bundled):
        sys.exit(1)


main()
