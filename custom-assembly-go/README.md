Chainguard's [Custom Assembly](https://www.chainguard.dev/docs/chainguard/containers/custom-assembly/overview) is a tool that allows customers to create customized containers with extra packages added. This enables customers to reduce their risk exposure by creating container images that are tailored to their internal organization and application requirements while still having few-to-zero CVEs.

You can use the Chainguard API to further customize your Custom Assembly builds and retrieve information about them. The demo application stored in this repository can be used to update a Custom Assembly container's configuration based on a provided YAML file.


## Prerequisites

In order to follow along with this guide, you will need the following:

* Access to a Custom Assembly container image. If your organization doesn't yet have access to Custom Assembly, reach out to your account team to start the process.
* The demo application used in this guide is written in Go, so you will need Go installed on your local machine. Refer to the [official Go documentation](https://go.dev/doc/install) for instructions on downloading and installing Go.
* You will also need [`chainctl` installed](https://www.chainguard.dev/docs/platform/chainctl-usage/how-to-install-chainctl) on your local machine to create a Chainguard token and authenticate to the Chainguard API.


## Running the Demo Application

Before you can run the demo application, there are a few steps you need to take in order for it to work properly.

First, run the following `go` commands:

```shell
go mod init github.com/chainguard-dev/sdk && go mod tidy
```

The `go mod init` command will initialize a new `go.mod` file in the current directory. Including the `github.com/chainguard-dev/sdk` URL tells Go to use that as the module path. The `go mod tidy` command ensures that the new `go.mod` file matches the source code in the module.

You must authenticate before you can interact with the Chainguard API. For this reason, this demo application expects an environment variable named `TOK` to be present when it's run. Create this environment variable with the following command:

```shell
export TOK=$(chainctl auth token)
```

Following that, open up `main.go` with your preferred text editor. This example uses `nano`:

```shell
nano main.go
```

From there, edit the following lines:

```
    	// Group and repository settings
    	defaultGroupName = "ORGANIZATION"
    	demoRepoName 	= "CUSTOM-IMAGE-NAME"
```

Replace `ORGANIZATION` with the name of your organization's repository within the Chainguard registry. This usually takes the form of a domain name, such as `example.com`. Additionally, replace `CUSTOM-IMAGE-NAME` with the name of a Custom Assembly image. This is typically a name like `custom-nginx` or `custom-python`. 

Save and close the `main.go` file. If you used `nano`, you can do so by pressing `CTRL+X`, `Y`, and then `ENTER`. 

Next, open up the `build.yaml` file:

```shell
nano build.yaml
```

This file will have the following content:

```
contents:
  packages:
	- wolfi-base
	- go
```

Here, replace `wolfi-base` and `go` with whatever packages you'd like to be included in the customized container image. Note that you can only add packages that your organization already has access to, based on the Chainguard Containers you have already purchased. Refer to the [Custom Assembly Overview](https://www.chainguard.dev/docs/chainguard/containers/custom-assembly/overview#limitations) for more details on the limitations of what packages you can add to a Custom Assembly image.

Save and close the `build.yaml` file. Finally, you can run the application to apply the configuration listed in the `build.yaml` file to your organization's Custom Assembly image:

```shell
go run main.go
```

The application will start by listing some information, including the specified organization's repositories and build reports for the chosen Custom Assembly image:

```
Group: example.com (ID: 45a0cEXAMPLE977f050c5fb9aEXAMPLEed764595)

All repositories in example.com:
- custom-assembly
- nginx
- curl

Repository: custom-assembly (ID: 45a0cEXAMPLE977f050c5fb9aEXAMPLEed764595/c375EXAMPLEb500c)

Build Reports for custom-node repository:

. . .

```

It will then prompt you to confirm that you want to apply the customization configuration listed in the `build.yaml` file:

```
About to apply customization using configuration file: build.yaml
Are you sure you want to update repository custom-node? (y/n): y
```

Enter `y` to confirm. Then, if everything was configured correctly, the application output will show successful build reports:

```
. . .

- Started: Mon, 28 Apr 2025 00:28:44 UTC, Result: Success, Digest: . . .
```

### Troubleshooting

Although the demo application has been tested to ensure that it works properly, there are several pitfalls one may encounter when they attempt to run it. 

For example, you may run into an error like the following:

```
Failed to list groups: rpc error: code = Internal desc = stream terminated by RST_STREAM with error code: PROTOCOL_ERROR
```

This may indicate that there is an issue with your Chainguard authentication token. To resolve this, try recreating the environment variable that holds the token:

```shell
export TOK=$(chainctl auth token)
```

You may also encounter errors like the following:

```
cannot find package "chainguard.dev/sdk/auth" . . .
```

This may indicate that the Chainguard SDK wasn't imported correctly. Be sure that you run the following commands to set this up:

```shell
go mod init github.com/chainguard-dev/sdk && go mod tidy
```

Again, the `main.go` file contains many comments that explain each portion of the code. If you encounter any errors, we encourage you to review the file closely to better understand how the application works and what might be going wrong.


## More Information

For a more detailed explanation of this demo application, we encourage you to follow the [associated tutorial in Chainguard Docs](https://www.chainguard.dev/docs/chainguard/containers/custom-assembly/custom-assembly-api-demo). You may also find our [overview of Custom Assembly](https://www.chainguard.dev/docs/chainguard/containers/custom-assembly/overview) to be of interest.