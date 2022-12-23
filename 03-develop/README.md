# Develop Application

Now that we have a deployment and a pod available, lets use it to develop our application.

## Terminal

The first step is to get a terminal into the pod. Update your `devspace.yaml` with the following configuration:

```yaml
dev:
  hello-world:
    imageSelector: loftsh/go
    terminal: {}
```

The `imageSelector` tells DevSpace how to find the container to connect to. You could also use a `labelSelector` which might be more familiar from using `kubectl`. The following would work as a `labelSelector` too:

```yaml
dev:
  hello-world:
    labelSelector:
      app.kubernetes.io/component: hello-world
    terminal: {}
```

The `app.kubernetes.io/component` label was added by the [component-chart](https://www.devspace.sh/component-chart/docs/introduction) we used to deploy the application. `imageSelector` may not work if you were to have multiple containers using the same image. In this case you could add a `container` property to specify the container to connect to. For this example, we'll stick with `imageSelector` since it's simple and works for our situation.

Now when we run `devspace dev` DevSpace won't exit immediately and instead leave us with a terminal to the running container of our deployment.

```sh
❯ devspace dev
info Using namespace 'russ'
info Using kube context 'kind-kind'
root@hello-world-devspace-68d996b7d5-jsfvt:/app#
```

This container has a number of tools pre-installed. For a complete list as well as other supported languages, see the [devspace-containers](https://github.com/loft-sh/devspace-containers) repository. To exit DevSpace, enter `exit` into the terminal and press enter.

## File Syncing

A terminal to a container with some tools installed is alright, but we'll also need code to build our application. We could use `vim` and `touch` commands to create the files, but using an IDE is much easier. Let's configure file syncing so that you can use any editor you'd like.

```yaml
dev:
  hello-world:
    labelSelector:
      app.kubernetes.io/component: hello-world
    terminal: {}
    sync:
      - path: ./:/app
```

Now create a file in the current directory named `main.go` and paste the following content into it:
```go
package main

import (
	"fmt"
	"net/http"
)

const (
	port = "8080"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, world!")
	})

	http.ListenAndServe(":"+port, nil)
}
```

Run `devspace dev` again, and run `pwd` then `ls` in the prompt, and you should see the three files in your workspace directory in the `/app` folder inside the container.
```sh
root@hello-world-devspace-68d996b7d5-jsfvt:/app# pwd
/app
root@hello-world-devspace-68d996b7d5-jsfvt:/app# ls
README.md  devspace.yaml  main.go
```

Since we don't need `README.md` and `devspace.yaml` to run our application, lets configure syncing to ignore them. DevSpace has [several options](https://www.devspace.sh/docs/configuration/dev/connections/file-sync#config-reference) for excluding files from syncing, but to keep it simple we'll use `excludePaths`:

```yaml
dev:
  hello-world:
    labelSelector:
      app.kubernetes.io/component: hello-world
    terminal: {}
    sync:
      - path: ./:/app
        excludePaths:
          - devspace.yaml
          - README.md
```
If you were to simply exit devspace and run `devspace dev` again, these files would still exist in the container because DevSpace would connect to the same running container as before. To see the `excludePaths` configuration in effect, we'll introduce the `devspace purge` command.

## `devspace purge`
 To clean up an old container, or just to save resources when not actively developing, you can run `devspace purge`. This will remove any deployments created by `devspace dev`. Here's an example output:
 
```sh
$ devspace purge
info Using namespace 'russ'
info Using kube context 'kind-kind'
dev:hello-world Stopping dev hello-world
dev:hello-world Scaling up Deployment hello-world...
purge:hello-world Deleting deployment hello-world...
purge:hello-world Successfully deleted deployment hello-world
```

Run `devspace dev` again, followed by listing the files in the container should show that these files were excluded from file syncing:

```sh
root@hello-world-devspace-68d996b7d5-slhpc:/app# ls
main.go
```

## Making Changes
Now that we have a terminal and file syncing configured, lets make some changes. Lets run our initial application with `go run main.go`:
```sh
root@hello-world-devspace-68d996b7d5-slhpc:/app# go run main.go

```

It doesn't look like anything is happening. Lets add a message so that we know that the program ran. Exit the application by pressing `control + C`, and add the print statement like below.
```go
package main

import (
	"fmt"
	"net/http"
)

const (
	port = "8080"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, world!")
	})

  fmt.Printf("Listening on port %s\n", port)
	http.ListenAndServe(":"+port, nil)
}
```
Running `go run main.go` again should yield some text output now:
```sh
root@hello-world-devspace-68d996b7d5-slhpc:/app# go run main.go
Listening on port 8080

```

You can see how the combination of file syncing and a terminal can make for a quick development cycle. We can now make changes and re-run our application in a matter of seconds. Compare this to a longer development cycle of making changes, building a docker image, deploying the application, and finally see the changes. Next we'll see how to send HTTP requests, since we are making a web application after all.

## Port Forwarding

So far our application is listening on port `8080`, but only requests coming from inside the cluster would be routed to our pod. We could configure a Kubernetes `Service` and `Ingress` and route traffic to our service, but port forwarding is an easier option for development. This has the advantage of not needing to coordinate domains, subdomains or paths between many users that might be developing apps on the same cluster. Eventually you'll want to configure ingress but usually that's scoped to one or two publicly available domains. Add the `ports` configuration to your `dev.hello-world` section, like below:

```yaml
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
      - port: "8080"
```

Now when you run `devspace dev` and run `go run main.go` in your pod again, DevSpace will portforward from `localhost:8080` to your pod. Go ahead and open http://localhost:8080 in your browser, and you should see `Hello, world!`.

Chances are you've already used port `8080` on your system before. This was the case for me, so I saw the following output:
```sh
$ devspace dev
info Using namespace 'russ'
info Using kube context 'kind-kind'
deploy:hello-world Skipping deployment hello-world
dev:hello-world Waiting for pod to become ready...
dev:hello-world Selected pod hello-world-devspace-68d996b7d5-slhpc
dev:hello-world sync  Sync started on: ./ <-> /app
dev:hello-world sync  Waiting for initial sync to complete
dev:hello-world sync  Initial sync completed
start_dev: forward ports: unable to listen on any of the requested ports: [{8080 8080}]
fatal exit status 1
```

A simple workaround is to map a different port to your pod. You won't need to change anything about the code, since you can map a different port using the format: `[LOCAL]:[REMOTE]`. Specifying `8080` was a shorthand for `8080:8080`. I've updated my configuration to use `9090:8080`, and now I can load http://localhost:9090 and see our application working!