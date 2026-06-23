# .NET Migration Example

## Overview

This directory contains an example of migrating a .NET application from
the upstream Microsoft .NET images to the equivalent Chainguard images.

This example is based on the [.NET demo in Chainguard's `cs-workshop`
repository](https://github.com/chainguard-demo/cs-workshop/tree/main/dotnet).

The example uses a simple application that is adapted from one of the [Microsoft
.NET
examples](https://github.com/dotnet/dotnet-docker/tree/main/samples/dotnetapp).
The application has been modified slightly to write output both to stderr and
to a file on disk.

## Steps

### Using Upstream Image

1. Build the image.

```
docker build -t dotnet-example:notlinky -f notlinky.Dockerfile .
```

2. Run the image.

```
docker run --rm dotnet-example:notlinky
```

3. Scan the image.

```
grype dotnet-example:notlinky
```

### Takeaways

1. Note that the upstream image runs as root by default, so in the Dockerfile
   they set a non-root user.
2. Vulnerabilities from the grype scan
3. Number of packages, files, etc.

### Using Chainguard Images

1. Build the image.

```
docker build -t dotnet-example:linky -f linky.Dockerfile .
```

2. Run the image.

```
docker run --rm dotnet-example:linky
```

3. Scan the image:

```
grype dotnet-example:linky
```

### Compare Image Sizes

```
docker images dotnet-example
```

### Takeaways

1. Note that Chainguard image does not run as a root user (by default)
2. 0 Vulnerabilities from the grype scan
3. Number of packages, files, etc.

## Compare Dockerfiles

Diff the two Dockerfiles with `git`:

```
git diff --no-index -U1000 notlinky.Dockerfile linky.Dockerfile
```

1. Note that container registries are different.
2. Both Dockerfiles use multistage builds.
3. It's necessary to switch to the root user to do `dotnet restore` because, unlike
   the Microsoft image, the Chainguard image runs as a non-root user.
4. In the runtime stage the Microsoft example switches to a non-root user. This
   is not needed in the Chainguard version because, by default, Chainguard
   images do not run as root.
5. We ensure that the workdir is owned by the runtime user using `--chown`, so
   that the application can write to a file under that path. The name of the
   default user in Chainguard images may be different across different images
   (i.e., `nonroot`, `app`, `build`), so it's less ambiguous to use the UID, which
   is always `65532`.
