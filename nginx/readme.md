# nginx Demo

This repository contains a simple demonstration of an nginx application using the Chainguard nginx Image. For a complete tutorial, please refer to the [Getting Started with the nginx Chainguard Image guide](https://www.chainguard.dev/docs/chainguard/containers/getting-started/web-and-data-services/nginx).

## Running the Demo

You must have nginx and Docker Engine installed on your machine to run this demonstration. 

Clone this repository to your local machine. Navigate to the `nginx` directory, where you will find a `Dockerfile`.

When the image is built, the `Dockerfile` will do the following:

1. Start a new build based on the `cgr.dev/chainguard/nginx:latest` image;
2. Expose port 8080 in the container for nginx to listen on;
3. Copy the HTML content from the data directory into the image.

You can now build the image with:

```shell
docker build . -t nginx-demo
```

Once the build is complete, run the image with:

```shell
docker run -d --name nginxcontainer -p 8080:8080 nginx-demo
```

The `-d` flag configures our container to run as a background process. The `--name` flag will name our container `nginxcontainer`, making it easy to identify from other containers. The `-p` flag publishes the port that the container listens on to a port on your local machine. This allows us to navigate to `localhost:8080` in a web browser of our choice to view the HTML content served by the container. You should see the same HTML page as before, with Linky and an octopus fun fact.

If you wish to publish to a different port on your machine, such as `1313`, you can do so by altering the command-line argument as shown:

```shell
docker run -d --name nginxcontainer -p 1313:8080 nginx-demo
```

When you are done with your container, you can stop it with the following command:

```shell
docker container stop nginxcontainer
```
