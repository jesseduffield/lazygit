# Flaggy Contribution Guidelines for Agents

This repository provides a zero-dependency command-line parsing library that must always rely exclusively on the Go standard library for runtime and test code. The core principles that follow preserve the project's lightweight nature and its focus on easily understandable, flat code.

## Core Principles
- Keep the codebase dependency-free beyond the Go standard library. Adding third-party modules for any purpose, including testing, is not permitted.
- Prefer flat, straightforward control flow. Avoid `else` statements when possible and limit indentation depth so examples remain approachable for beginners.
- Optimize every change for readability. Favor descriptive names, small logical blocks, and explanatory comments that teach users how the parser behaves.

## Documentation Expectations
- Maintain clear, beginner-friendly explanations throughout the codebase. Add comments to **every** function and test describing what they do and why they matter to the overall library.
- Annotate each stanza of code with concise comments, even when the logic appears self-explanatory.
- Keep primary documentation accurate. Update `README.md` and `CONTRIBUTING.md` whenever your modifications alter usage instructions, contribution workflows, or observable behavior.

## Tooling Requirements
- Always run `go fmt`, `go vet`, and `goimports` over affected packages before committing.
- Favor consistent formatting and import organization that highlight the minimal surface area of each example.

## Testing Guidance
- When writing tests, ensure the accompanying comment explains exactly what is being verified.
- Leave benchmarks and examples with clarifying comments so readers immediately understand the intent and scope of each scenario.

Following these guidelines keeps Flaggy's codebase welcoming to newcomers and aligned with its lightweight philosophy.
