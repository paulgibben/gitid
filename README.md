# Git Identity Manager (gitid)

A terminal-based tool that helps developers manage multiple Git identities easily through an interactive interface.

![demo](https://github.com/user-attachments/assets/8ec86e59-2cb1-47b7-9acd-54a7d0d8f20f)

## Features

- 🔄 Switch between multiple Git identities globally
- 📁 Set repository-specific local identities
- 🏷️ Optional nicknames for quick identity identification
- ➕ Add new identities interactively
- 🗑️ Delete unwanted identities
- 💻 Terminal-based UI with keyboard navigation
- 🔒 Uses Git's built-in configuration system
- 🔍 Smart identity matching by nickname, name, or email
- 📋 Command-line interface for automation and scripting
- ✨ Visual indicators for current global and local identities

## Installation


### From Binary

Download the appropriate binary for your platform from the [releases page](https://github.com/nathabonfim59/gitid/releases).

### From Package (Linux)

#### Debian/Ubuntu:
```bash
sudo dpkg -i gitid_*.deb
```

#### RedHat/Fedora:
```bash
sudo rpm -i gitid_*.rpm
```

#### Arch Linux (AUR):
```bash
# Pre-compiled binary
yay -S gitid-bin

# Build from source
yay -S gitid-git
```

### Building from Source

#### Prerequisites

- Go 1.21 or later
- Make
- NFPM (for package generation)

#### Build Commands

```bash
# Build for your current platform
make build

# Build static binary (Linux only)
make build-static

# Create releases for all platforms and packages
make release

# Clean build artifacts
make clean
```

## Usage

GitID provides both an interactive TUI and command-line interface for managing Git identities.

### Interactive Mode

Run `gitid` (without arguments) to start the interactive interface.

#### Keyboard Controls

- `↑`/`↓` or `j`/`k` - Navigate through identities
- `Enter` - Select identity to set as global
- `r` - Set selected identity as local for current repository
- `R` - Remove local identity (use global)
- `D` - Delete selected identity
- `e` - Edit nickname for selected identity
- `E` - Edit full identity (name, email, nickname)
- `←`/`→` - Navigate confirmation dialog
- `Esc` - Cancel current action
- `q` - Quit application

### Command Line Interface

GitID also provides a full CLI for automation and scripting:

#### Global Identity Management
```bash
# List all identities
gitid list

# Show current global (and local if in repo) identity
gitid current

# Switch global identity
gitid switch <nickname|name|email>
gitid use <nickname|name|email>        # Alias for switch

# Add new identity
gitid add "Full Name" "email@domain.com" [nickname]

# Delete identity
gitid delete <nickname|name|email>

# Set/update nickname
gitid nickname <identifier> <new-nickname>
```

#### Repository-Specific (Local) Identity Management
```bash
# Show current local identity for repository
gitid repo current

# Set existing identity as local for current repository
gitid repo use <nickname|name|email>

# Add new identity and set as local for current repository
gitid repo add "Full Name" "email@domain.com" [nickname]

# Remove local identity override (fall back to global)
gitid repo remove
```
