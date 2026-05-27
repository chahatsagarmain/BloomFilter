# Contributing to Simple Bloom Filter

Thank you for your interest in contributing to this project! We welcome all contributions, from bug reports and documentation updates to performance improvements and new features.

Follow this guide to get started.

---

## Development Setup

### Prerequisites

To build and run this project, make sure you have the following installed:

- **Go**: Version 1.18 or higher.
- **Make**: (Optional but highly recommended) for executing shortcuts.

### Getting Started

1. **Fork** the repository on GitHub.
2. **Clone** your fork locally:
   ```bash
   git clone https://github.com/YOUR-USERNAME/BloomFilter.git
   cd BloomFilter
   ```
3. **Add the upstream remote** to keep your local branch synced with the main project:
   ```bash
   git remote add upstream https://github.com/chahatsagarmain/BloomFilter.git
   ```

---

## Makefile Workflow

We use a simple `Makefile` to automate common development workflows. Ensure your changes compile and pass tests by using the following commands:

- **Format Code**: Formats Go source files.
  ```bash
  make fmt
  ```
- **Run Tests**: Runs the package unit and statistical tests.
  ```bash
  make test
  ```
- **Build Binaries**: Compiles the CLI application.
  ```bash
  make build
  ```
- **Run CLI**: Starts the interactive CLI tool directly.
  ```bash
  make run
  ```
- **Clean Up**: Deletes compiled binaries and clean the workspace.
  ```bash
  make clean
  ```

---

## Contribution Guidelines

To maintain code quality and consistency, please follow these guidelines:

### 1. Code Formatting
All Go source files must conform to the standard style formatting. Always run `make fmt` before committing your changes.

### 2. Testing
- Write corresponding unit tests for any new features or bug fixes.
- Ensure that the empirical/statistical tests still succeed.
- Run `make test` locally to verify that your changes do not break any existing code.

### 3. Clean Commit Messages
Write clear, concise commit messages that describe the intent of your change. For example:
- `feat: add option to select custom hash functions`
- `fix: prevent counter underflow in counting bloom filter`
- `docs: update installation instructions`

---

## Pull Request Process

1. **Create a new branch** from the `main` branch for your work:
   ```bash
   git checkout -b feature/your-feature-name
   ```
2. Make your changes, formatting them with `make fmt`.
3. Add tests to cover your changes and run `make test` to ensure everything is green.
4. **Push** your branch to your GitHub fork:
   ```bash
   git push origin feature/your-feature-name
   ```
5. **Open a Pull Request** (PR) on GitHub against the `main` branch of the upstream repository.
6. Provide a description of the changes and link any related issues in the PR template.
