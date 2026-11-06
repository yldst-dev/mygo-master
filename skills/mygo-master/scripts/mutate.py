#!/usr/bin/env python3
import json
import os
import subprocess
import sys
import tempfile

USAGE = "usage: mutate.py <project> <mutations.json>"


def go_test(directory, *args):
    command = ["go", "test", "-count=1", "-timeout", "180s", *args, "."]
    done = subprocess.run(command, cwd=directory, capture_output=True, text=True)
    return done.returncode, done.stdout + done.stderr


def package_directory(project, mutation):
    if "package" in mutation:
        return os.path.normpath(os.path.join(project, mutation["package"]))
    return os.path.dirname(os.path.join(project, mutation["file"]))


def classify(code, output):
    if "[build failed]" in output or "[setup failed]" in output:
        return "BROKEN", "the mutation does not compile"
    if "no tests to run" in output:
        return "BROKEN", "no test matches this name"
    if code == 0:
        return "MISSED", ""
    if "panic: test timed out" in output:
        return "caught", "by a hang"
    return "caught", ""


def main():
    if len(sys.argv) != 3:
        sys.exit(USAGE)
    project = os.path.abspath(sys.argv[1])
    with open(sys.argv[2]) as f:
        mutations = json.load(f)
    directories = {package_directory(project, m) for m in mutations}
    for directory in sorted(directories):
        code, output = go_test(directory)
        if code != 0:
            print(output)
            sys.exit(f"the tests fail before any mutation in {directory}")
    counts = {"caught": 0, "MISSED": 0, "BROKEN": 0}
    with tempfile.TemporaryDirectory() as work:
        for index, mutation in enumerate(mutations):
            source = os.path.join(project, mutation["file"])
            test = mutation.get("test", "")
            with open(source) as f:
                text = f.read()
            found = text.count(mutation["old"])
            if found != 1:
                verdict, note = "BROKEN", f"the old text appears {found} times"
            else:
                mutated = os.path.join(work, f"{index}.go")
                with open(mutated, "w") as f:
                    f.write(text.replace(mutation["old"], mutation["new"]))
                overlay = os.path.join(work, f"{index}.json")
                with open(overlay, "w") as f:
                    json.dump({"Replace": {source: mutated}}, f)
                arguments = ["-overlay", overlay]
                if test:
                    arguments += ["-run", f"^({test})$"]
                code, output = go_test(package_directory(project, mutation), *arguments)
                verdict, note = classify(code, output)
            counts[verdict] += 1
            name = test or "all tests"
            suffix = f"  ({note})" if note else ""
            print(f"{index:3} {verdict:7} {name}{suffix}")
    print(f"{counts['caught']} caught, {counts['MISSED']} missed, {counts['BROKEN']} broken")
    if counts["MISSED"] or counts["BROKEN"]:
        sys.exit(1)


main()
