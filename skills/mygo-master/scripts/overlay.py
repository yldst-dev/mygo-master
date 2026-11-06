#!/usr/bin/env python3
import json
import os
import re
import shutil
import subprocess
import sys

USAGE = "usage: overlay.py <package directory> <work directory> <file>...\nprints the overlay file for go test -overlay"


def main():
    if len(sys.argv) < 4:
        sys.exit(USAGE)
    package = os.path.abspath(sys.argv[1])
    work = os.path.abspath(sys.argv[2])
    name = subprocess.run(["go", "list", "-f", "{{.Name}}", "."], cwd=package, capture_output=True, text=True, check=True).stdout.strip()
    copies = os.path.join(work, "overlay")
    shutil.rmtree(copies, ignore_errors=True)
    os.makedirs(copies)
    replace = {}
    for path in sys.argv[3:]:
        base = os.path.basename(path)
        with open(path) as f:
            text = re.sub(r"^package \w+", f"package {name}", f.read(), count=1, flags=re.M)
        copy = os.path.join(copies, base)
        with open(copy, "w") as f:
            f.write(text)
        replace[os.path.join(package, "zz_" + base)] = copy
    overlay = os.path.join(work, "overlay.json")
    with open(overlay, "w") as f:
        json.dump({"Replace": replace}, f, indent=1)
    print(overlay)


main()
