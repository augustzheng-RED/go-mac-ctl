#!/opt/homebrew/bin/python3
import csv
import json
import subprocess
import sys
import tempfile
from pathlib import Path

def preprocess_image(image_path: str, scale_factor: float) -> str:
    temp = tempfile.NamedTemporaryFile(delete=False, suffix=".png")
    temp_path = Path(temp.name)
    temp.close()

    width, height = image_size(image_path)
    subprocess.run(
        [
            "sips",
            "-z",
            str(max(1, int(height * scale_factor))),
            str(max(1, int(width * scale_factor))),
            image_path,
            "--out",
            str(temp_path),
        ],
        check=True,
        capture_output=True,
        text=True,
    )

    return str(temp_path)


def image_size(image_path: str) -> tuple[int, int]:
    proc = subprocess.run(
        ["sips", "-g", "pixelWidth", "-g", "pixelHeight", image_path],
        check=True,
        capture_output=True,
        text=True,
    )
    width = 0
    height = 0
    for line in proc.stdout.splitlines():
        line = line.strip()
        if line.startswith("pixelWidth:"):
            width = int(line.split(":", 1)[1].strip())
        elif line.startswith("pixelHeight:"):
            height = int(line.split(":", 1)[1].strip())

    if width <= 0 or height <= 0:
        raise ValueError(f"failed to read image size for {image_path}")

    return width, height


def run_tesseract_tsv(image_path: str) -> str:
    proc = subprocess.run(
        ["tesseract", image_path, "stdout", "tsv"],
        check=True,
        capture_output=True,
        text=True,
    )
    return proc.stdout


def parse_targets(tsv_text: str, query: str, scale_factor: float) -> list[dict]:
    query = query.strip().lower()
    rows = csv.DictReader(tsv_text.splitlines(), delimiter="\t")
    targets = []

    for row in rows:
        text = (row.get("text") or "").strip()
        if not text:
            continue

        if query and query not in text.lower():
            continue

        try:
            confidence = float(row["conf"])
            left = int(round(int(row["left"]) / scale_factor))
            top = int(round(int(row["top"]) / scale_factor))
            width = int(round(int(row["width"]) / scale_factor))
            height = int(round(int(row["height"]) / scale_factor))
        except (KeyError, TypeError, ValueError):
            continue

        targets.append(
            {
                "Text": text,
                "Confidence": confidence,
                "Bounds": {
                    "X": left,
                    "Y": top,
                    "Width": width,
                    "Height": height,
                },
            }
        )

    targets.sort(key=lambda item: item["Confidence"], reverse=True)
    return targets


def main() -> int:
    if len(sys.argv) not in (3, 4):
        print("usage: ocr_find_text.py <image-path> <query> [scale]", file=sys.stderr)
        return 2

    image_path, query = sys.argv[1], sys.argv[2]
    scale_factor = 2.0
    if len(sys.argv) == 4:
        try:
            scale_factor = float(sys.argv[3])
        except ValueError:
            print(f"invalid scale: {sys.argv[3]}", file=sys.stderr)
            return 2
        if scale_factor <= 0:
            print(f"scale must be positive: {scale_factor}", file=sys.stderr)
            return 2

    processed_path = None
    try:
        processed_path = preprocess_image(image_path, scale_factor)
        tsv_text = run_tesseract_tsv(processed_path)
        targets = parse_targets(tsv_text, query, scale_factor)
    except subprocess.CalledProcessError as exc:
        message = exc.stderr.strip() or exc.stdout.strip() or str(exc)
        print(message, file=sys.stderr)
        return exc.returncode or 1
    finally:
        if processed_path:
            Path(processed_path).unlink(missing_ok=True)

    print(json.dumps({"Targets": targets}))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
