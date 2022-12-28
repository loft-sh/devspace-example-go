# Deploying to Production
We've seen how DevSpace can help you bootstrap a prototype application, provide a tight development loop, and help debug your application as problems arise. Now we'll see how to prepare your application for production. We've used a development optimized docker image, but now it's time to produce a production optimized image.

## Create a Dockerfile
Our first step is to create a Dockerfile that will be used to run our application. This image won't contain any development tools. This helps reduce the time it takes to deploy the image but also helps reduce the security vulnerabilities that could come with more software. Create a new file named `Dockerfile` and add the following content:

```Dockerfile
# Builder Image
FROM golang:1.17 as builder
ARG TARGETOS
ARG TARGETARCH
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /app/hello-world

# Production Image
FROM gcr.io/distroless/static AS production
COPY --from=builder --chown=nonroot:nonroot /app/hello-world /
USER nonroot:nonroot
ENTRYPOINT [ "/hello-world" ]
```

This uses a multi-stage build. The first stage `builder` produces an executable binary. The second stage `production` simply copies the binary to an execution environment and runs it as a `nonroot` user. There are even smaller images possible, but the `gcr.io/distroless/static` provides a convenient environment to run applications as a `nonroot` user without needing to agonize over certificates and user creation.

You might think that we'll need to run `docker build . -t ...` to produce our image, but this is something that DevSpace helps with. Let's update `devspace.yaml` with the information it needs to build the image. We'll use [ttl.sh](https://ttl.sh/), an anonymous, publicly available image registry for this example, but you'll likely want to use a more permanent registry for your own production images. This also gives an opportunity to show [DevSpace variables](https://www.devspace.sh/docs/configuration/variables) in action:

```yaml
vars:
  IMAGE_NAME: $(uuidgen | tr "[:upper:]" "[:lower:]")

images:
  hello-world:
    image: ttl.sh/${IMAGE_NAME}:2h
```

Here we've created a variable `IMAGE_NAME` that gets populated with the results of the command, `uuidgen` (after being lowercased). This variable is then used in our `images` configuration to make a unique image tag that expires in 2 hours, according to ttl.sh's instructions.

Typically, you would want to use your production image in your deployments, so lets go ahead and replace all references to `loftsh/go` with our production image. You might be thinking that your development workflow will be totally different from here on, but DevSpace has a [devImage](https://www.devspace.sh/docs/configuration/dev/modifications/dev-image) feature that will replace your production image with a development one whenever you're running `devspace dev`. When running `devspace deploy` your production image will be used as normal. We'll use this feature to maintain our development environment as it was before.

Additionally, you might find it awkward to repeat your production image in multiple places. DevSpace has a special placeholder for images using the format `image(name):tag(name)`. We'll use this format so that if we change our image repository, we'll only need to update it in one place. The complete `devspace.yaml` should now look like:

```yaml
version: v2beta1
name: hello-world

vars:
  IMAGE_NAME: $(uuidgen | tr "[:upper:]" "[:lower:]")

images:
  hello-world:
    image: ttl.sh/${IMAGE_NAME}:2h

deployments:
  hello-world:
    helm:
      values:
        containers:
          - image: image(hello-world):tag(hello-world)

dev:
  hello-world:
    imageSelector: image(hello-world):tag(hello-world)
    devImage: loftsh/go
    terminal: {}
    sync:
      - path: ./:/app
        excludePaths:
          - devspace.yaml
          - README.md
    ports:
      - port: "9090:8080"
```

You might be thinking that we've never run our production optimized image, and you would be correct. To see if it runs, we can restore the production optimized deployment with `devspace reset pods`, or we can use `devspace purge` followed by `devspace deploy`.