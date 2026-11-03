# C/C++ Demo

This repository contains three examples showing single and multi-stage builds for C-based programs using Chainguard Images. For a complete tutorial, please refer to the [Getting Started with the C/C++ Chainguard Images guide](https://www.chainguard.dev/docs/chainguard/containers/getting-started/languages-and-runtimes/c).

# Running the Demo

To run this demo on your machine, you will need to have both [Docker Engine](https://docs.docker.com/engine/install) and the [GNU Compiler Collection](https://gcc.gnu.org/) installed locally.

Clone this repository locally and navigate to the `c` directory in your terminal. Here, you can test the `hello.c` program by compiling it with `gcc`.

```sh
gcc -Wall -o hello hello.c
```

You can execute the resultant binary with the following command:

```sh
./hello
```

You will see the following output in your terminal.

```Output
Hello, world!
I am a demo from Chainguard Docs.
My code was written in C.
```

Now, with Docker Engine running, you can compile this program inside of the `gcc-glibc` Chainguard Image. You will use `Dockerfile1` for this image build.

`Dockerfile1` will:
1. Use the `gcc-glibc:latest` Chainguard Image as the base image;
2. Create and set the current working directory to `/home/build`;
3. Copy the `hello.c` program code to the current directory;
4. Compile the program and name it `hello`;
5. Copy the compiled binary to `/usr/bin`;
6. Set the image to run as a non-root user; and,
7. Execute the compiled binary when the container is started.

Execute the following command to initiate the image build.
```sh
docker build -f Dockerfile1 -t example1:latest .
```

Once your image build it complete, run the following command to start a container.
```sh
docker run --name example1 example1:latest
```

You will see output in your terminal identical to that of the binary you compiled locally.

If you wish to follow along with the other examples in this demonstration, please check out our complete [Getting Started with the C/C++ Chainguard Images Guide](https://www.chainguard.dev/docs/chainguard/containers/getting-started/languages-and-runtimes/c) in [Chainguard Docs](https://www.chainguard.dev/docs)