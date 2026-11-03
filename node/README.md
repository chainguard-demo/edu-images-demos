# node-demo

This directory holds a sample JavaScript application to be used within a Node environment. For a more detailed explanation of this example application we encourage you to refer to our [Getting Started with the Node Chainguard Image](https://www.chainguard.dev/docs/chainguard/containers/getting-started/languages-and-runtimes/node) tutorial.

## Running the demo

From within this directory, run the following command to create a new `package.json` file:

```shell
npm init -y
```

Next, install the application dependencies. Specifically, you'll need `ronin-server` and `ronin-mocks`. These will create a "mock" server that saves JSON data in memory and returns it in subsequent GET requests to the same endpoint.

```shell
npm install ronin-server ronin-mocks
```

Following that, the the demo application can be built into a container image using the Dockerfile included in this example repository. 

This Dockerfile will perform the following actions:

1. Start a new image based on the `cgr.dev/chainguard/node:latest` image;
2. Set the work dir to `/app` inside the container;
3. Copy application files from the current directory to the `/app` location in the container;
4. Run `npm install` to install production-only dependencies;
5. Set up additional arguments to the default entrypoint (`node`), specifying which script to run.

Build the application image with the following command:

```shell
docker build . -t wolfi-node-server
```

Once the build is finished, run the image:

```shell
docker run --rm -it -p 8000:8000 wolfi-node-server
```

Although the application is running from within a container, this command will cause it to block your terminal we set up a port redirect to receive requests on `localhost:8000` as the application waits for connections on port `8000`.

From a new terminal window, run the following command. This will make a POST request to your application sending a JSON payload:

```shell
curl --request POST \
  --url http://localhost:8000/test \
  --header 'content-type: application/json' \
  --data '{"msg": "testing node wolfi image" }'
```

If the connection is successful, you will receive output like this in the terminal where the application is running:

```shell
2023-02-07T15:48:54:2450  INFO: POST /test
```

You can now query the same endpoint to receive the data that was stored in memory when you run the previous command:

```shell
curl http://localhost:8000/test
```
```shell
{"code":"success","meta":{"total":1,"count":1},"payload":[{"msg":"testing node wolfi image","id":"6011f987-b9f8-4442-8253-d54166df5966","createDate":"2023-02-07T15:57:23.520Z"}]}
```

When you're finished, you can close the application by pressing `CTRL+C`. 