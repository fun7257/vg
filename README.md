# vg

Virtual Go - A Fast and Flexible Go Version Manager

## Installation

### Quick Install

```bash
curl -fsSL https://raw.githubusercontent.com/fun7257/vg/main/install.sh | bash
```

Or download the install script and run it:

```bash
./install.sh
```

## Development

### Setup

1. Clone the repository:
```bash
git clone https://github.com/fun7257/vg.git
cd vg
```

2. Install Git hooks (runs golangci-lint on each commit):
```bash
make install-hooks
# or
./scripts/install-hooks.sh
```

3. Install golangci-lint (required for pre-commit hook):
```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Build

```bash
make build
```

The binary will be in `.tmp/build/vg`.

## Git Hooks

The project includes a pre-commit hook that runs `golangci-lint` before each commit to ensure code quality.

To install the hooks:
```bash
make install-hooks
```

To skip the hook (not recommended):
```bash
git commit --no-verify
```
