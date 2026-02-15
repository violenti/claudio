# Claudio :chicken:


[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8.svg)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](https://choosealicense.com/licenses/gpl-3.0/)

**Claudio** is a lightweight AI assistant powered by Anthropic and  OpenAI (just drop in your API key). Support for other LLMs is on the roadmap.

## :rocket: Features

- [x] Anthropic (Claude) integration.
- [x] OpenAI support.
- [x] Streaming support.
- [ ] Local LLM support via Ollama (coming soon).


## :hammer_and_wrench: Prerequisites

You need **Go 1.25** or higher installed on your machine.

**macOS:**
```bash
brew install go
```

Linux/Windows: Check the official installation guide.


## :inbox_tray: Installation

Clone the repository and build the binary:

```bash

git clone https://github.com/violenti/claudio.git

cd claudio

make

```
Or you could execute: 

``` bash 
curl -fsSL https://raw.githubusercontent.com/violenti/claudio/main/install.sh | bash

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

 :computer: Usage 

 ```bash

./claudio

```

Pro tip: Move the binary to your path to run it from anywhere: sudo mv claudio /usr/local/bin/


## Contributing

Pull requests are welcome. For major changes, please open an issue first
to discuss what you would like to change.

Please make sure to update tests as appropriate.

:page_facing_up: License 

[GPL V3](https://choosealicense.com/licenses/gpl-3.0/)