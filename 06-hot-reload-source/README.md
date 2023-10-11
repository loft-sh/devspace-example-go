# Hot Reloading from Source
To further speed up our development process, lets see an example of hot reloading. This example will configure DevSpace to synchronize the golang source to the container, and use the container environment to run the application using `go run`. This has the advantage of only synchronizing incremental changes instead of a potentially large binary, and automatically building for the container environment.

## Synchronize the golang source
In earlier examples, we simply synchronized everything from the current directory to the `/app` folder. Let's be a little more selective, and only synchronize things that should trigger a rebuild of the golang application, namely any `*.go` files and `go.mod` and `go.sum`:

```yaml
dev:
  hello-world:
    imageSelector: image(hello-world):tag(hello-world)
    devImage: loftsh/go
    sync:
      - path: ./:/app
        excludePaths:
          - '**'
          - '!**/*.go'
          - '!go.mod'
          - '!go.sum'
```

The sync `path` stayed the same, however we've updated our `excludePaths` to use some wildcard matching and negation features. We start by excluding everything with a recursive wildcard `'**'`. Then we want to _include_ all golang source files with `'**/*.go'`. To do this, we negate the exclusion with `'!**/*.go'`. Similarly, we add `!go.mod` and `!go.sum` to the list.

## Run the Application
In earlier examples, we simply used the terminal to run the application be entering `go run main.go`. This manual restarting of the application might get repetitive, so let's have DevSpace start the application for us, and restart it when changes are made.

We'll start with setting the development container's start command:
```yaml
dev:
  hello-world:
    imageSelector: image(hello-world):tag(hello-world)
    devImage: loftsh/go
    command:
      - go
      - run
      - main.go
```

Next, we'll tell DevSpace to restart the container whenever a file is synchronized:
```yaml
dev:
  hello-world:
  # ... excluded for clarity
    sync:
      - path: ./:/app
        excludePaths:
          - '**'
          - '!**/*.go'
          - '!go.mod'
          - '!go.sum'
        onUpload:
          restartContainer: true
```

Now, whenever a `**/*.go` file is changed, DevSpace will kill the current process and run `go run main.go` again. It may seem like `restartContainer` will restart the pod container, however it uses a special restart helper script that only reruns the given command. This allows us to take advantage of the build cache on the container as well.

# See Application Logs

Finally, we may not need the terminal, and might prefer to see the log output of our application. To do this, we only need to remove the `terminal` config, and add `logs`. Here's what our final `dev` configuration should look like:
```yaml
dev:
  hello-world:
    imageSelector: image(hello-world):tag(hello-world)
    devImage: loftsh/go
    logs: {}
    command:
      - go
      - run
      - main.go
    sync:
      - path: ./:/app
        excludePaths:
          - '**'
          - '!**/*.go'
          - '!go.mod'
          - '!go.sum'
        onUpload:
          restartContainer: true
    ports:
      - port: "9090:8080"
      - port: "23450:2345"
```

As before, this can be started with `devspace dev`. I've included some `*.txt` files so that you can see that the container does not restart when those are changed. Go ahead and give it a try!