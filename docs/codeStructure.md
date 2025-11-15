# Structure
`
cmd/
├── main.go                # Entry point, initializes Bubble Tea app and config
internal/
├── app/
│   ├── app.go             # Main Bubble Tea app logic, state management
│   ├── ui/                # UI components (views, widgets, rendering)
│   │   ├── explorer.go    # Tree explorer view
│   │   ├── editor.go      # Snippet editor view
│   │   └── search.go      # Search/filter view
│   └── handlers/          # Event handlers (key presses, mouse, etc.)
├── config/
│   ├── config.go          # Config loading/saving, defaults
│   └── configtypes.go     # Config struct definitions
├── filesystem/
│   ├── watcher.go         # fsnotify logic for file changes
│   ├── metadata.go        # Metadata file creation/reading/writing
│   ├── index.go           # Centralized index management
│   └── utils.go           # Helper functions for file operations
├── models/
│   ├── snippet.go         # Snippet struct and metadata definitions
│   └── index.go           # Index struct for search/filter
└── utils/
    ├── logger.go          # Logging utilities
    └── helpers.go         # General helper functions
`

# Key Points

- **cmd/main.go**: 
    Entry point, sets up the Bubble Tea app and loads config.
- **internal/app/**
  - **app.go**: 
        Main app logic, state management, and Bubble Tea integration.
  - **ui/**:
        UI components for different views (explorer, editor, search).
  - **handlers/**: 
        Event handlers for user input.
- **internal/config/**: 
    Config loading/saving and struct definitions.
- **internal/filesystem/**: 
    File watching, metadata management, and index updates.
- **internal/models/**:
    Struct definitions for snippets and index.
- **internal/utils/**:
    Logging and general helpers.

This structure keeps your code modular, making it easier to maintain and extend as your app grows.
