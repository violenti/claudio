# Claudio :chicken:


[![Go Version](https://img.shields.io/badge/go-1.25+-00ADD8.svg)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](https://choosealicense.com/licenses/gpl-3.0/)

**Claudio** is a lightweight AI assistant powered by Anthropic and OpenAI (just drop in your API key), with local LLM support via Ollama. Support for other providers is on the roadmap.

## :rocket: Features

- [x] Anthropic (Claude) integration.
- [x] OpenAI support.
- [x] Streaming support.
- [x] Local LLM support via Ollama.
- [x] Configurable Foghorn Leghorn quotes while Claudio thinks. :rooster:

## :hammer_and_wrench: Prerequisites

Go is **only** required if you want to build Claudio from source. If you just want to run the prebuilt binary via `install.sh`, you can skip this section.

You need **Go 1.25** or higher installed on your machine.

**macOS:**
```bash
brew install go
```

Linux/Windows: Check the official installation guide.


## :inbox_tray: Installation

**Option 1: Prebuilt binary (no Go required)**

```bash
curl -fsSL https://raw.githubusercontent.com/violenti/claudio/main/install.sh | bash
```

This downloads the latest prebuilt `claudio` binary for your platform and installs it to `/usr/local/bin`.

**Option 2: Build from source (requires Go, see Prerequisites)**

```bash

git clone https://github.com/violenti/claudio.git

cd claudio

make

```

This will generate the claudio binary in the root directory.


:gear: Configuration
Claudio relies on environment variables for authentication.

Export your Anthropic API Key:

```bash

 export ANTHROPIC_API_KEY= ""

 ```


 ```bash

 export OPENAI_API_KEY= ""

 ```

### :wrench: Config file

Claudio reads its model list and system prompt from `~/.claudio/config.json`. See [`config.example.json`](config.example.json) for the expected format.

If you installed via `install.sh`, a default config is set up automatically from `config.example.json` (an existing file is never overwritten, so your customizations are safe). If you built from source, copy the example config to your home directory:

```bash
mkdir -p ~/.claudio
cp config.example.json ~/.claudio/config.json
```

**Windows:** `install.sh` does not run on Windows, so set up the config file manually. Claudio looks for it at `%USERPROFILE%\.claudio\config.json`. In PowerShell:

```powershell
New-Item -ItemType Directory -Force "$env:USERPROFILE\.claudio" | Out-Null
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/violenti/claudio/main/config.example.json" -OutFile "$env:USERPROFILE\.claudio\config.json"
```

(If you cloned the repo, you can copy `config.example.json` there instead of downloading it.)

Edit the file to add or remove models per provider, or to tweak Claudio's system prompt.

### :rooster: Foghorn Leghorn quotes

While Claudio is thinking, it displays a random Foghorn Leghorn quote. The quotes are loaded from `~/.claudio/foghorn_quotes.json`, so you can fully customize them.

If you installed via `install.sh`, the quotes file is set up automatically (an existing file is never overwritten, so your customizations are safe). If you built from source, copy the bundled quotes file to your home directory:

```bash
mkdir -p ~/.claudio
cp assets/foghorn_quotes.json ~/.claudio/foghorn_quotes.json
```

**Windows:** `install.sh` does not run on Windows, so set up the quotes file manually. Claudio looks for it at `%USERPROFILE%\.claudio\foghorn_quotes.json`. In PowerShell:

```powershell
New-Item -ItemType Directory -Force "$env:USERPROFILE\.claudio" | Out-Null
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/violenti/claudio/main/assets/foghorn_quotes.json" -OutFile "$env:USERPROFILE\.claudio\foghorn_quotes.json"
```

(If you cloned the repo, you can copy `assets\foghorn_quotes.json` there instead of downloading it.)

The file is a simple JSON object with a `quotes` array — add, remove, or replace entries as you like:

```json
{
  "quotes": [
    "I say, I say, boy — pay attention when I'm talkin' to ya!",
    "That's a joke, son, a flag-waver.",
    "Go away, son, ya bother me."
  ]
}
```

Each time Claudio processes a prompt, it picks one quote at random from the list. This feature is optional: if the file is missing, invalid, or the `quotes` array is empty, Claudio simply skips the quote and shows the spinner as usual.

 :computer: Usage 

 ```bash

./claudio

```

Pro tip: If you built from source, move the binary to your path to run it from anywhere: sudo mv claudio /usr/local/bin/ (not needed if you installed via install.sh, it already does this for you).


## Contributing

Pull requests are welcome. For major changes, please open an issue first
to discuss what you would like to change.

Please make sure to update tests as appropriate.

:page_facing_up: License 

[GPL V3](https://choosealicense.com/licenses/gpl-3.0/)