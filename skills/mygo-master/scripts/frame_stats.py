#!/usr/bin/env python3
import re
import sys

LINE = re.compile(
    r"mygo: frame (\d+) took ([\d.]+) ms: build ([\d.]+)(?: \(([^)]*)\))?, "
    r"layout ([\d.]+), paint ([\d.]+), present ([\d.]+) \((.*?)\), other ([\d.]+)"
)


def percentile(ordered, fraction):
    return ordered[min(len(ordered) - 1, int(len(ordered) * fraction))]


def main():
    if len(sys.argv) < 2:
        sys.exit("usage: frame_stats.py <log> [budget ms, default 16.7]")
    budget = float(sys.argv[2]) if len(sys.argv) > 2 else 16.7
    frames = []
    with open(sys.argv[1], errors="replace") as f:
        for line in f:
            match = LINE.search(line)
            if match:
                frames.append(match.groups())
    if not frames:
        sys.exit("no frame lines: run the app with MYGO_FRAME_STATS=all")
    took = sorted(float(frame[1]) for frame in frames)
    average = sum(took) / len(took)
    print(
        f"{len(took)} frames: avg {average:.2f} ms, p50 {percentile(took, 0.5):.1f}, "
        f"p95 {percentile(took, 0.95):.1f}, p99 {percentile(took, 0.99):.1f}, max {took[-1]:.1f}"
    )
    for label, index in (("build", 2), ("layout", 4), ("paint", 5), ("present", 6), ("other", 8)):
        print(f"  {label:8} avg {sum(float(frame[index]) for frame in frames) / len(frames):.2f} ms")
    repainted = sum(1 for frame in frames if frame[3] == "painted again")
    print(f"  {repainted} frames only painted, without running the view")
    slow = sorted((frame for frame in frames if float(frame[1]) > budget), key=lambda frame: -float(frame[1]))
    print(f"{len(slow)} frames over {budget} ms")
    for frame in slow[:10]:
        print(
            f"  frame {frame[0]}: {frame[1]} ms, build {frame[2]}, layout {frame[4]}, "
            f"paint {frame[5]}, present {frame[6]}, other {frame[8]}"
        )


main()
