# Add a Helm Deployment

So far we've only learned how to configure the kubernetes context and namespace that DevSpace will use. Now it's time to bootstrap our web application. Update your `devspace.yaml`:

```yaml
version: v2beta1
name: hello-world

deployments:
  hello-world:
    helm:
      values:
        containers:
          - image: loftsh/go
            command:
              - tail
              - -f
              - /dev/null
```

This instructs DevSpace to use a helm chart to create a deployment named `hello-world` with a single image `loftsh/go`. We've added a command to this container, because otherwise it would simply exit and the deployment would never become available. Now run `devspace dev` again:

```sh
❯ devspace dev
info Using namespace 'russ'
info Using kube context 'kind-kind'
deploy:hello-world Deploying chart /Users/russellcentanni/.devspace/component-chart/component-chart-0.8.5.tgz (hello-world) with helm...
deploy:hello-world Deployed helm chart (Release revision: 1)
deploy:hello-world Successfully deployed hello-world with helm
```

We see that by default, DevSpace will use a [component-chart](). This is a generic pre-made helm chart that can quickly bootstrap prototype applications. We'll use this chart for now, but later we'll see how to use a different chart and how to make a new helm chart of our own.

Our application doesn't really do anything, so lets configure DevSpace to help us develop it in the next section.

