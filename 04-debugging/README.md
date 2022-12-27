# Debugging
Now that we have port forwarding it opens up the possibility of connecting our IDE to a debugging session and stepping through our code. Our application is super simple at the moment, but as it evolves you may want to peek into details of an incoming request (for example). We'll show how to configure this using VS Code and Goland, so feel free to jump the section relevant to your IDE, but first, lets learn how to run `dlv` (delve) and configure DevSpace for debugging since it'll be the same for either IDE.

## DevSpace Configuration
Configuring devspace for debugging is as simple as starting port forwarding on the port we'll run the debugger on. We'll use port `2345` in our example. Since we don't want to use a privileged port for our local computer, let's map it from `23450`. Remember that we just need to add `23450:2345` to our `ports` configuration. The complete `dev` section should now look like:

```yaml
...
dev:
  hello-world:
    imageSelector: loftsh/go
    terminal: {}
    sync:
      - path: ./:/app
        excludePaths:
          - devspace.yaml
          - README.md
    ports:
      - port: "9090:8080"
      - port: "23450:2345"
```

## Delve (dlv)
[Delve](https://github.com/go-delve/delve) is the go-to debugger for Go. You won't need to install it, since our `loftsh/go` image comes with it pre-installed. We just need to run `devspace dev` and run `dlv` with the correct arguments once the container terminal is available.

```sh
$ dlv debug main.go --listen 0.0.0.0:2345 --headless --api-version=2 --output /tmp/__debug_bin
```

We use the `--listen` to tell dlv to listen for connections from any IP on port 2345. The `--headless` option tells it to wait for a connection before starting our application. `--output` tells dlv to write the compile application to a path that's typically accessible and won't be synced back to our workspace by the 2-way file synchronization. `--api-version=2` instructs `dlv` to use a protocol that's compatible with VS Code's debugger.

Now that `dlv` is started and waiting for a connection, lets configure VSCode or Goland to connect to it.

## VS Code Debugging
Perhaps the easiest way to set up debugging with VS Code is to go to the "Run and Debug" tab and click the "create a launch.json file" link under the "Run and Debug" button. Note that the "Run and Debug" button seems tempting, but it will run our application locally, and not inside the pod that we've been using to develop our application so far.

Once you've clicked that link, you'll be given a series of options. Use these values for each question:
1. Go
2. Go: Connect to Server
3. 127.0.0.1
4. 23450

Once all the questions have been answered, you should see a generated `.vscode/launch.json` file similar to the following contents. Alternatively you can create the file manually.
```json
{
  // Use IntelliSense to learn about possible attributes.
  // Hover to view descriptions of existing attributes.
  // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Connect to server",
      "type": "go",
      "request": "attach",
      "mode": "remote",
      "remotePath": "${workspaceFolder}",
      "port": 23450,
      "host": "127.0.0.1"
    }
  ]
}
```

This would _almost_ work, but unless we were to recreate the full path of our workspace inside the container, we need to tell VSCode what the `remotePath` should be. Change `remotePath` to `/app` and select the "Connect to server" launch configuration and press the play button next to it.

Inside your DevSpace terminal you should see the `Listening on port 8080` message get logged, showing us that the debugger connection was successful and that our application has started. Open `main.go` in VS Code and set a breakpoint inside the function passed to `http.HandleFunc`. Open a browser tab to http://localhost:9090/ and you should see the debugger stop on the line where you set your breakpoint.

From here you can use the debugger to inspect the incoming request (`http.Request`) and the response writer (`http.ResponseWriter`), and any other structs or data that you may add as you develop your application.

## Goland Debugging
To get started debugging with Goland, lets create a new run configuration. Under the "Run" menu, select "Debug..." and click "Add new run configuration...". A dropdown will appear where you should select "Go Remote". In the form that appears, enter the following values:
- Name: `Debug: Hello World`
- Host: `localhost`
- Port: `23450`

Leave the other values as-is and click the "Apply" button to save your changes. Now in a terminal, run `devspace dev` and run the `dlv` command. Goland will display a `dlv` command to run, but we'll use the command from the earlier section to that the compile output isn't synchronized to our local workspace:
```sh
$ dlv debug main.go --listen 0.0.0.0:2345 --headless --api-version=2 --output /tmp/__debug_bin
```

Once `dlv` is running we can use the drop down of run configurations in the top right to select our "Debug: Hello World" configuration and run it by pressing the green bug button. Inside your DevSpace terminal you should see the `Listening on port 8080` message get logged, showing us that the debugger connection was successful and that our application has started. Open `main.go` in VS Code and set a breakpoint inside the function passed to `http.HandleFunc`. Open a browser tab to http://localhost:9090/ and you should see the debugger stop on the line where you set your breakpoint.

## Summary
Goland and VS Code are just two examples of debugging with DevSpace and an IDE. Other popular IDEs should work similarly. The things to remember are to use a remote debugging configuration and running`dlv` is run with the correct arguments. Now you have the tools you need to get to the bottom of most development issues!