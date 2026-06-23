#! env bash
. ../base.sh

for i in $(docker images -q dotnet-example); do
  docker rmi $i
done

clear

banner "Here's a simple .NET application and its Dockerfile."
pe "$BATCAT Program.cs"
pe "$BATCAT notlinky.Dockerfile"

banner "Build the image and run it."
pe "docker build -t dotnet-example:notlinky -f notlinky.Dockerfile ."
pe "docker run --rm dotnet-example:notlinky"

banner "Check the size."
pe "docker images dotnet-example:notlinky"

banner "And the vulnerabilities."
pe "grype dotnet-example:notlinky"

banner "Let's migrate it to a Chainguard image."
pe "git diff --no-index -U1000 notlinky.Dockerfile linky.Dockerfile"
pe "docker build -t dotnet-example:linky -f linky.Dockerfile ."
pe "docker run --rm dotnet-example:linky"

banner "It should be a bit smaller."
pe "docker images dotnet-example"

banner "And have significantly less vulnerabilities."
pe "grype dotnet-example:linky"
