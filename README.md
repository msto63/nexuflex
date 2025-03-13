# Nexuflex

Nexuflex is a universal platform that consists of Application Servers with modular Business Services and text-based client interfaces to support any business process. It enables users to interact with various business modules through a unified command-line interface, regardless of the specific business domain.

## Overview

Nexuflex follows a client-server architecture where the client provides a Text User Interface (TUI) and the server hosts modular business services. The core concept is a consistent command interface where users enter commands in a standardized format that the system processes through the appropriate business modules.

### Key Features

- **Universal Command Interface**: Standardized command format across all business domains
- **Modular Architecture**: Clear separation between client, application server, and business services
- **Business Service Plugins**: Easily extend the system with new business capabilities
- **Text-based UI**: Efficient keyboard-driven interface for power users
- **Cross-Platform Support**: Runs on Linux, macOS, and Windows
- **Multilingual Support**: Internationalization with language-specific message files
- **Alias Management**: Create shortcuts for frequently used commands
- **Automatic Server Discovery**: Easily connect to available servers on the network

## Architecture

Nexuflex consists of three main components:

1. **Client**: Text-based user interface that accepts commands and displays results
2. **Application Server**: Core server that handles authentication, command parsing, and routing
3. **Business Services**: Modular plugins that implement specific business functions

Communication between client and server uses gRPC for efficient, typed, and asynchronous messaging.

## Command Format

Nexuflex uses a standardized command format:

```
<BusinessService>.<Action>.<SubAction> <Parameter1> <Parameter2> ...
```

For example:
- `Finance.Create.Report Q4_2024 "Profit and Loss"` - Creates a financial report
- `HR.Update.Employee 12345 "John Doe" Role=Manager` - Updates employee information
- `Inventory.List.Items WarehouseA` - Lists items in a warehouse

## Project Structure

```
nexuflex/
├── shared/                  # Shared components and protocols
│   └── proto/               # gRPC protocol definitions
├── nexuflex-client/         # Client application
│   ├── config/              # Configuration management
│   ├── core/                # Core client functionality
│   ├── i18n/                # Internationalization
│   ├── ui/                  # User interface components
│   └── lang/                # Language files
├── nexuflex-server/         # Application server
│   ├── auth/                # Authentication and session management
│   ├── command/             # Command parsing and execution
│   ├── services/            # Service management
│   └── server/              # Server implementation
├── nexuflex-core-services/  # Core system services
│   ├── system/              # System management
│   ├── admin/               # Administrative functions
│   └── help/                # Help and documentation
└── business-services/       # Business-specific service modules
    ├── nexuflex-finance/    # Financial services
    ├── nexuflex-hr/         # Human resources services
    └── nexuflex-inventory/  # Inventory management services
```

## Getting Started

### Prerequisites

- Go 1.20 or higher
- Git

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/msto63/nexuflex.git
   cd nexuflex
   ```

2. Build the client and server:
   ```bash
   # Build the client
   cd nexuflex-client
   go build -o nexuflex-client
   
   # Build the server
   cd ../nexuflex-server
   go build -o nexuflex-server
   
   # Build core services
   cd ../nexuflex-core-services
   go build -o nexuflex-core-services
   ```

3. Start the server:
   ```bash
   ./nexuflex-server
   ```

4. Start the client:
   ```bash
   ./nexuflex-client
   ```

### Configuration

#### Client Configuration

The client can be configured through a `client.ini` file, which can be placed in:
- The current directory
- User config directory: `~/.config/nexuflex/client.ini` (Linux/macOS) or `%APPDATA%\nexuflex\client.ini` (Windows)

Example configuration:
```ini
[server]
address = localhost
port = 50051
use_tls = false
auto_discover = true
discovery_token = NEXUFLEX_DISCOVERY
discover_timeout_seconds = 5

[ui]
color_scheme = default
header_text = nexuflex Terminal
show_timestamps = true
enable_sounds = false
max_output_lines = 1000
max_history_entries = 100
auto_complete_enabled = true
auto_fill_service_prefix = true
language = en

[commands]
save_history = true
use_local_aliases = true
max_local_aliases = 50
enable_multiline_input = true
save_history_on_shutdown = true
```

#### Server Configuration

The server is configured through a `server.ini` file, which can be placed in:
- The current directory
- System config directory

#### Language Configuration

Language files are located in the `lang` directory and follow the naming convention `<language-code>.ini`.

Available languages:
- English (en.ini)
- German (de.ini)

To add a new language, create a new INI file based on the existing ones.

## Client Usage

### Command Line Arguments

```
Usage: nexuflex-client [options]

Options:
  -config string     Path to config file
  -server string     Server address (IP or hostname)
  -port int          Server port
  -discover          Enable automatic server discovery
  -discover-timeout  Timeout for server discovery in seconds (default 5)
  -debug             Enable debug output
  -lang string       Language code (e.g., 'en', 'de')
```

### Keyboard Shortcuts

- `Ctrl+H` - Show help
- `Ctrl+L` - Open login dialog
- `Ctrl+D` - Start server discovery
- `Ctrl+C` - Exit application
- `↑/↓` - Navigate through command history
- `Tab` - Command completion

### Basic Commands

- `help` or `?` - Show help
- `exit` or `quit` - Exit application
- `clear` or `cls` - Clear output
- `history` - Show command history
- `connect <host> [port]` - Connect to a server
- `disconnect` - Disconnect from server
- `login` - Open login dialog
- `logout` - Log out
- `alias` - Show all defined aliases
- `alias <name>=<command>` - Define a new alias
- `unalias <name>` - Delete an alias
- `use <service>` - Set service context

## Development

### Adding New Business Services

To create a new business service:

1. Create a new Go module for your service
2. Implement the `BusinessService` interface
3. Register your service with the Application Server

Example structure for a new business service:
```
nexuflex-myservice/
├── main.go            # Service registration
├── myservice.go       # Service implementation
├── commands.go        # Command definitions
└── models.go          # Data models
```

### Internationalization

To add support for a new language:

1. Create a new file in the `lang` directory named `<language-code>.ini`
2. Copy the content from an existing language file (e.g., `en.ini`)
3. Translate all messages

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- The project uses [tview](https://github.com/rivo/tview) for the text-based user interface
- Communication is handled by [gRPC](https://grpc.io/) 
