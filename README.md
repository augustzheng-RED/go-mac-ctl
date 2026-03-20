# go-mac-ctl

`go-mac-ctl` is a small macOS CLI for browser and desktop actions.

It is designed for simple UI automation flows such as:
- open a page in Chrome
- find visible text on screen
- click a text target
- type text
- press keys
- use AI to choose the best matching on-screen target from a natural-language instruction

## Requirements

- macOS
- Go
- Google Chrome
- Python 3
- Tesseract OCR
- Accessibility permission for Terminal
- Screen Recording permission for Terminal

For `ai-click`:
- an OpenAI API key in `.env`

## Setup

Create a local env file from the template:

```bash
cp env.example .env
```

Then edit `.env` and set:

```env
OPENAI_API_KEY=your_real_key_here
OPENAI_MODEL=gpt-4.1-mini
```

The CLI loads `.env` automatically on startup.

## Build

Run directly:

```bash
go run .
```

Or build a binary:

```bash
go build .
```

## Commands

- `chrome-open [url]`
  Open a URL in Google Chrome and wait until the page is ready.

- `open [url-or-path]`
  Open a URL, file, or folder with the default macOS handler.

- `find-text [query]`
  Screenshot the current screen and return OCR matches as JSON.

- `click-text [query] [index]`
  Find text on screen and click the indexed OCR match.

- `type [text]`
  Type text into the currently focused field.

- `key [name]`
  Press a supported key. Current supported values: `enter`, `tab`, `escape`.

- `ai-click [instruction]`
  Screenshot the current screen, run OCR, let OpenAI choose the best candidate for the instruction, then click it.

## Examples

Open Apple in Chrome:

```bash
go run . chrome-open https://www.apple.com
```

Find all visible matches for `iPhone`:

```bash
go run . find-text iPhone
```

Click the first OCR match for `iPhone`:

```bash
go run . click-text iPhone 0
```

Type into the focused field and press Enter:

```bash
go run . type "apple.com"
go run . key enter
```

Use AI to choose the best `iPhone` navigation target:

```bash
go run . chrome-open https://www.apple.com
go run . ai-click "click iphone tab"
```

Use AI to select a more specific target on the current page:

```bash
go run . ai-click "click iphone 16e product"
```

## Notes

- Screenshots are stored in `tmp/screenshots/`.
- The CLI keeps only the newest 20 screenshots and removes older ones automatically.
- `env.example` is a placeholder template. Do not commit your real `.env`.
