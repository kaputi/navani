# Structure
`
snippets/
├── config/
│   └── config.json          # App config (root directory, filters, etc.)
├── metadata/
│   ├── index.json           # Centralized index for search/filter (auto-generated)
│   └── <snippet-name>.meta.json  # Per-snippet metadata (e.g., binaryTree.ts.meta.json)
├── snippets/
│   ├── go/
│   │   ├── binaryTree.go
│   │   └── binaryTree.go.meta.json
│   ├── typescript/
│   │   ├── binaryTree.ts
│   │   └── binaryTree.ts.meta.json
│   └── ...                  # Other languages/directories
├── logs/
│   └── app.log              # Optional: logs for debugging
└── .gitignore               # Ignore metadata/index.json if using Git
`
# Key Points
- **config/**: 
    Stores app-level configuration (root directory, default filters, etc.).
- **metadata/**:
    - **index.json**: Centralized metadata for fast search/filter (auto-updated by background task).
    - **Per-snippet .meta.json files**: Individual metadata for each snippet, stored alongside the snippet file.
- **snippets/**:
    - Organized by language or user-defined directories.
    - Each snippet file has a corresponding .meta.json file in the same directory.
- **logs/**: Optional directory for app logs, useful for debugging.
- **.gitignore**: Ensures auto-generated files (like index.json) are not tracked by version control.
