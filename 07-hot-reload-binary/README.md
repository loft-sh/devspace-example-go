# Hot Reloading a Compiled Binary
To further speed up our development process, lets see an example of hot reloading. This example will configure DevSpace to synchronize a locally compiled binary to the deployed pod. This can have some caveats, since we'll need to ensure that the binary is compiled for the container architecture, however it may be the quickest options over synchronizing lots of files or attempting to compile it in a container with limited CPU and memory resources.

## Configure Pipelines
In order to rebuild our binary when files change, we'll use a pipeline function called `run_watch`. This will watch for file changes according to a match pattern, and run the given command when changes are detected. We'll also configure DevSpace so that it runs the regular `dev` pipeline in parallel with our binary compiling pipeline.

First, lets set up an overridden `dev` pipeline to run our parallel pipelines:
```yaml
pipelines:
  dev: |-
    run_pipelines compile-app default-dev
  compile-app: |-
  default-dev: |-
```

Now when running `devspace dev` DevSpace will start two pipelines that run in parallel. Lets make the `default-dev` pipeline run the built-in `devspace dev` functionality with `run_default_pipeline`:
```yaml
pipelines:
  dev: |-
    run_pipelines compile-app default-dev
  compile-app: |-
  default-dev: |-
    run_default_pipeline dev
```

Next, lets use `run_watch` to watch for file changes and rebuild our application's binary. We'll assume a local script file exists that has all the commands necessary to build the application.
```yaml
pipelines:
  dev: |-
    run_pipelines compile-app default-dev
  compile-app: |-
    # Do an initial build
    ./hack/build.sh
    
    # Rebuild on changes...
    run_watch --path "./**/*.go" -- ./hack/build.sh
  default-dev: |-
    run_default_pipeline dev
```

Note that `run_watch` won't build the application initially, so we've run our build command to ensure our binary exists. Also the `path` **must** be in quotes to prevent bash expansion from passing invalid arguments to `run_watch`.

The `./hack/build.sh` script builds the binary for the container environment, so you may not be able to run the same binary locally:
```shell
#!/bin/sh

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/server main.go
```

Now we have our two pipelines:
- `compile-app` does the initial compile of our binary, and re-builds it on `**/*.go` file changes.
- `default-dev` does everything the original `dev` pipeline would have had we not overridden it.

We're almost done. We only need to configure DevSpace to synchronize the binary to the container and restart the application.

## Run the Application
Let's change the `dev` configuration to start the compiled binary instead of showing us a terminal with the `command` property. We'll likely still want to see the application's log output, so we'll add `log` configuration as well.
```yaml
dev:
  hello-world:
    imageSelector: image(hello-world):tag(hello-world)
    devImage: loftsh/go
    command:
      - /app/server
    logs: {}
```

But wait... how does `/app/server` make it into the container? We'll add `sync` configuration to copy it from the local environment:
```yaml
dev:
  hello-world:
    imageSelector: image(hello-world):tag(hello-world)
    devImage: loftsh/go
    # ... omitted for clarity
    sync:
      - path: ./build/server:/app/server
        file: true
```
Notice that since we're synchronizing a single file, we've set `file: true` so that DevSpace handles it properly.

Finally, we'll tell DevSpace to restart the container whenever the binary is synchronized to the container:
```yaml
dev:
  hello-world:
    # ... omitted for clarity
    sync:
      - path: ./build/server:/app/server
        file: true
        onUpload:
          restartContainer: true
```

## Run `devspace dev`!
With the two pipelines running in parallel, we have one that performs the file synchronization, port-forwarding, and container restarting according to the `dev` configuration, and another that watches for file changes and recompiles our binary! With the combination, we have a hot reloading development experience that uses a compiled binary.

On some systems, you may not see the application restart as changes are made. Run `devspace dev --debug` to gather more information. If you encounter a `text file busy` error when syncing the binary directly, a simple workaround can be used. Instead of running the synced binary directly, create a copy and run the copy instead:
```yaml
dev:
  hello-world:
    command:
      - sh
      - -c
      - cp /app/server /app/server-run && /app/server-run
```

To avoid surprises in mixed environments, this is the command that used in the final `devspace.yaml`.