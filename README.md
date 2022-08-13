# devspace-example-go

Welcome to the `devspace-example-go` project! In this repo, you'll find information about:
- [Debugging](#debugging)
- [Monitoring](#monitoring)
- [Logging](#logging)
- [Advanced Configuration](#advanced-configuration)

You'll also find example devspace.yaml files and projects to use as reference.

## Debugging

DevSpace has a debugging mode for Go projects using [delve](https://github.com/go-delve/delve) and [VSCode](https://code.visualstudio.com/)!

To run your project in debug mode, follow these steps:

### Modify `devspace.yaml`

In `devspace.yaml`, under the `dev` section, add a section with the ports you'd like your debugger process to listen on. 

Define a `DEBUGGER PORT` (the remote port to run your debugger process) and a `LOCAL PORT` (the port to forward the remote process to your local host.)

```yaml
...
dev:
  ...
  ports:
    - imageSelector: {PROJECT-IMAGE}
      forward:
        - port: {LOCAL-PORT}
          remotePort: {DEBUGGER-PORT}
...
```

In Vcluster, the [`devspace.yaml`](https://github.com/loft-sh/vcluster/blob/main/devspace.yaml#L70-L74) includes:

```yaml
...
dev:
  terminal:
    imageSelector: ${SYNCER_IMAGE}
    command: ["./devspace_start.sh"]
  ports:
    - imageSelector: ${SYNCER_IMAGE}
      forward:
        - port: 2346
          remotePort: 2345
...
```

(`${SYNCER_IMAGE}` in this case is defined under the `vars` section as `ghcr.io/loft-sh/loft-enterprise/dev-vcluster`.)

### Start the delve process

In your working directory, run:

```bash
dlv debug {PROJECT-NAME} --listen=0.0.0.0:{PORT} --api-version=2 --output /tmp/__debug_bin --headless --build-flags={BUILD-FLAGS} -- start
```

The command is defined in Loft's Vcluster in the [`devspace.sh` file](https://github.com/loft-sh/vcluster/blob/main/devspace_start.sh#L11).

```bash
dlv debug ./cmd/vcluster/main.go --listen=0.0.0.0:2345 --api-version=2 --output /tmp/__debug_bin --headless --build-flags=\"-mod=vendor\" -- start
```

This command will eventually output:

```bash
API server listening at: [::]:{PORT}
```

### Attach the debugger process to VSCode

To add the debugger process to a VSCode session, add this section to your VSCode `launch.json` configuration file:

```json
{
  "name": "Debug {APP-NAME} (localhost:{LOCAL-HOST-PORT})",
  "type": "go",
  "request": "attach",
  "mode": "remote",
  "port": {LOCAL-HOST-PORT},
  "host": "localhost",
  "substitutePath": [
    {
      "from": "${workspaceFolder}",
      "to": "{DEV-CONTAINER-NAME}",
    },
  ],
  "showLog": true,
  //"trace": "verbose", // use for debugging problems with delve (breakpoints not working, etc.)
}
```

You can find the Vcluster VSCode configuration file [here](https://github.com/loft-sh/vcluster/blob/main/.vscode/launch.json#L7-L22). The configuration section for the delve process looks like:

```json
{
  "name": "Debug vcluster (localhost:2346)",
  "type": "go",
  "request": "attach",
  "mode": "remote",
  "port": 2346,
  "host": "localhost",
  "substitutePath": [
    {
      "from": "${workspaceFolder}",
      "to": "/vcluster-dev",
    },
  ],
  "showLog": true,
  //"trace": "verbose", // use for debugging problems with delve (breakpoints not working, etc.)
},
```

### Stopping the Debugger

Once the delve process starts, it **must** connect to a VSCode debugger session. It cannot be stopped with a `CTRL-C` command.

Similarly, when the VSCode session is disconnected, the delve process is killed.

## Monitoring

## Logging

## Advanced Configuration