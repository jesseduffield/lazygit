#!/usr/bin/env python3
import re
from collections import defaultdict
from subprocess import check_output

README_FILE = "README.md"


lines = check_output(["go", "run", "./cmd/chroma/main.go", "--list"]).decode("utf-8").splitlines()
lines = [line.strip() for line in lines if line.startswith("  ") and not line.startswith("   ")]
lines = sorted(lines, key=lambda l: l.lower())

table = defaultdict(list)

for line in lines:
    table[line[0].upper()].append(line)

rows = []
for key, value in table.items():
    rows.append("{} | {}".format(key, ", ".join(value)))
tbody = "\n".join(rows)

with open(README_FILE, "r") as f:
    content = f.read()

with open(README_FILE, "w") as f:
    marker = re.compile(r"(?P<start>:----: \\| --------\n).*?(?P<end>\n\n)", re.DOTALL)
    replacement = r"\g<start>%s\g<end>" % tbody
    updated_content = marker.sub(replacement, content)
    f.write(updated_content)

print(tbody)
