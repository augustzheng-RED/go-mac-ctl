#!/usr/bin/env python3
import base64
import json
import mimetypes
import os
import sys
import urllib.error
import urllib.request
from pathlib import Path


def load_candidates() -> list[dict]:
    try:
        data = json.load(sys.stdin)
    except json.JSONDecodeError as exc:
        raise ValueError(f"invalid candidate JSON on stdin: {exc}") from exc

    if not isinstance(data, list):
        raise ValueError("candidate JSON must be a list")

    return data


def image_data_url(image_path: str) -> str:
    mime_type, _ = mimetypes.guess_type(image_path)
    if not mime_type:
        mime_type = "image/png"

    encoded = base64.b64encode(Path(image_path).read_bytes()).decode("ascii")
    return f"data:{mime_type};base64,{encoded}"


def build_prompt(instruction: str, candidates: list[dict]) -> str:
    prompt_candidates = []
    for list_index, candidate in enumerate(candidates):
        prompt_candidate = dict(candidate)
        prompt_candidate["list_index"] = list_index
        prompt_candidates.append(prompt_candidate)

    return (
        "You are selecting the best UI element to click on a macOS screenshot.\n"
        "Choose exactly one candidate from the OCR candidate list that best matches the user's intent.\n"
        "Prefer navigation tabs, buttons, or links that satisfy the instruction.\n"
        "If several candidates share the same text, use the screenshot and the bounds to choose the most likely one.\n"
        "Use list_index in your answer.\n"
        "If nothing is suitable, return candidate_index -1.\n"
        "Return JSON only with keys candidate_index and reason.\n\n"
        f"Instruction: {instruction}\n\n"
        f"OCR candidates:\n{json.dumps(prompt_candidates, ensure_ascii=False)}"
    )


def extract_content(message_content) -> str:
    if isinstance(message_content, str):
        return message_content

    if isinstance(message_content, list):
        parts = []
        for item in message_content:
            if isinstance(item, dict) and item.get("type") == "text":
                parts.append(item.get("text", ""))
        return "".join(parts)

    raise ValueError("unexpected message content format")


def extract_json_object(text: str) -> dict:
    text = text.strip()
    if text.startswith("```"):
        text = text.strip("`")
        text = text.replace("json\n", "", 1).strip()

    start = text.find("{")
    end = text.rfind("}")
    if start == -1 or end == -1 or end < start:
        raise ValueError(f"response did not contain a JSON object: {text!r}")

    return json.loads(text[start : end + 1])


def main() -> int:
    if len(sys.argv) != 3:
        print("usage: openai_pick_target.py <image-path> <instruction>", file=sys.stderr)
        return 2

    api_key = os.environ.get("OPENAI_API_KEY")
    if not api_key:
        print("OPENAI_API_KEY is not set", file=sys.stderr)
        return 2

    image_path = sys.argv[1]
    instruction = sys.argv[2]
    model = os.environ.get("OPENAI_MODEL", "gpt-4.1-mini")
    candidates = load_candidates()

    body = {
        "model": model,
        "response_format": {"type": "json_object"},
        "messages": [
            {
                "role": "user",
                "content": [
                    {"type": "text", "text": build_prompt(instruction, candidates)},
                    {"type": "image_url", "image_url": {"url": image_data_url(image_path)}},
                ],
            }
        ],
        "max_tokens": 200,
    }

    request = urllib.request.Request(
        "https://api.openai.com/v1/chat/completions",
        data=json.dumps(body).encode("utf-8"),
        headers={
            "Content-Type": "application/json",
            "Authorization": f"Bearer {api_key}",
        },
        method="POST",
    )

    try:
        with urllib.request.urlopen(request, timeout=60) as response:
            payload = json.loads(response.read().decode("utf-8"))
    except urllib.error.HTTPError as exc:
        detail = exc.read().decode("utf-8", errors="replace")
        print(detail, file=sys.stderr)
        return exc.code or 1
    except urllib.error.URLError as exc:
        print(str(exc), file=sys.stderr)
        return 1

    try:
        message = payload["choices"][0]["message"]["content"]
        content = extract_content(message)
        parsed = extract_json_object(content)
        candidate_index = int(parsed.get("candidate_index", -1))
        reason = str(parsed.get("reason", "")).strip()
    except (KeyError, IndexError, TypeError, ValueError, json.JSONDecodeError) as exc:
        print(f"failed to parse model response: {exc}", file=sys.stderr)
        return 1

    print(
        json.dumps(
            {
                "CandidateIndex": candidate_index,
                "Reason": reason,
                "Usage": payload.get("usage", {}),
            }
        )
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
