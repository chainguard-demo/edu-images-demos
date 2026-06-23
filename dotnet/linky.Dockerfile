FROM cgr.dev/chainguard/dotnet-sdk AS build

WORKDIR /source

COPY dotnetapp.csproj .
USER root
RUN dotnet restore

COPY . .
RUN dotnet publish --no-restore -o /app

FROM cgr.dev/chainguard/dotnet-runtime

COPY --from=build --chown=65532:65532 /app /app

WORKDIR /app

ENTRYPOINT ["./dotnetapp"]
