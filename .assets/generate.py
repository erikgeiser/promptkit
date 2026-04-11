#!/usr/bin/env python3

import os
import shutil
import subprocess
import sys
from concurrent.futures import ThreadPoolExecutor, as_completed
from pathlib import Path

ASSETS_DIR = Path(__file__).resolve().parent
EXAMPLES_DIR = ASSETS_DIR / "examples"
REPO_DIR = ASSETS_DIR.parent


def find_example_dir(tape: Path) -> Path | None:
    name = tape.stem  # e.g. "textinput_custom"
    candidate = REPO_DIR / "examples" / name
    return candidate if candidate.is_dir() else None


def build_example(example_dir: Path) -> None:
    cmd = ["go", "build", "-o", str(EXAMPLES_DIR) + "/", str(example_dir)]
    env = {**os.environ, "GOOS": "linux"}
    result = subprocess.run(cmd, env=env, capture_output=True, text=True, cwd=REPO_DIR)
    if result.returncode != 0:
        raise RuntimeError(
            f"build failed for {example_dir.name}:\n{result.stderr.strip()}"
        )


def render_tape(tape: Path) -> None:
    tmp_path = str(ASSETS_DIR / "tmp") + ":" + os.environ.get("PATH", "")
    cmd = [
        "docker",
        "run",
        "--rm",
        "-v",
        f"{ASSETS_DIR}:/vhs",
        "ghcr.io/charmbracelet/vhs",
        tape.name,
    ]
    env = {**os.environ, "PATH": tmp_path}
    result = subprocess.run(
        cmd, env=env, capture_output=True, text=True, cwd=ASSETS_DIR
    )
    if result.returncode != 0:
        raise RuntimeError(
            f"vhs render failed for {tape.name}:\n{result.stderr.strip()}"
        )


def process_tape(tape: Path) -> str:
    example_dir = find_example_dir(tape)
    if example_dir:
        print(f"  building  {example_dir.name} …")
        build_example(example_dir)
        print(f"  built     {example_dir.name}")
    else:
        print(f"  (no example dir for {tape.stem}, skipping build)")

    print(f"  rendering {tape.name} …")
    render_tape(tape)
    print(f"  rendered  {tape.name}")
    return tape.name


def main() -> None:
    selector = (
        "*.tape" if sys.argv == 0 else sys.argv[1].removesuffix(".tape") + ".tape"
    )
    print(selector)

    tapes = sorted(ASSETS_DIR.glob(selector))
    if not tapes:
        print("No .tape files found in", ASSETS_DIR)
        sys.exit(1)

    # Recreate the staging directory for compiled binaries
    if EXAMPLES_DIR.exists():
        shutil.rmtree(EXAMPLES_DIR)
    EXAMPLES_DIR.mkdir()

    print(f"Processing {len(tapes)} tape(s) in parallel …\n")

    errors: list[str] = []
    with ThreadPoolExecutor() as executor:
        futures = {executor.submit(process_tape, t): t for t in tapes}
        for future in as_completed(futures):
            tape = futures[future]
            try:
                future.result()
            except Exception as exc:
                errors.append(f"{tape.name}: {exc}")
                print(f"ERROR {tape.name}: {exc}", file=sys.stderr)

    # Clean up staging binaries regardless of errors
    shutil.rmtree(EXAMPLES_DIR)

    if errors:
        print(f"\n{len(errors)} error(s) occurred:", file=sys.stderr)
        for e in errors:
            print(f"  {e}", file=sys.stderr)
        sys.exit(1)

    print("\nDone.")


if __name__ == "__main__":
    main()
