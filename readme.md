# ai-code-editor

A command-line tool that uses AI to help make code changes to your project.

## Overview

ai-code-editor is a CLI tool that uses Ollama's AI models to help modify your codebase. You provide a prompt describing what changes you want to make, optionally specify which files to focus on, and the AI will suggest specific code modifications.

## Getting Started

1. Install [Ollama](https://ollama.ai/) on your system
2. Clone this repository
3. Create a `.env` file with your configuration
4. Run the tool:
   ```
   go run main.go <model> "your prompt" [files...]
   ```

## Components

* **main.go**: Entry point that processes CLI arguments and coordinates the editing flow
* **codeEditor/**: Core editing logic
  - `code-editor.go`: Handles the code modification process
  - `ai-response-parser.go`: Processes AI suggestions into file changes
  - `actions/`: File modification actions
* **services/**: Supporting functionality
  - `directory_tree.go`: Provides project structure context to the AI
  - `file_context_provider.go`: Reads and provides file contents
  - `base_prompt_provider.go`: Constructs AI prompts
* **ollama/**: AI integration
  - `client.go`: Handles communication with Ollama models
* **config/**: Configuration handling
  - `config.go`: Loads and manages configuration settings


## Requirements

* Go 1.20 or higher
* Ollama installed and running locally
* Valid .env configuration

## License

MIT License - See `LICENSE` file for details