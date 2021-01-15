# Debugging Dendrite

## VSCode

- Make sure Go extension is insalled (https://marketplace.visualstudio.com/items?itemName=golang.Go)
- Ctrl+Shift+P -> Go: Install/Update Tools -> dlv
- Specify a launch.json file:
    ```
    {
    // Use IntelliSense to learn about possible attributes.
    // Hover to view descriptions of existing attributes.
    // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Package",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "cwd": "${workspaceFolder}",
            "program": "${file}"
        }
    ]
    }
    ```
- Go to `cmd/dendrite-monolith-server/main.go` and start the debugger