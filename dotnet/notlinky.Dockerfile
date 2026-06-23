FROM mcr.microsoft.com/dotnet/sdk AS build

WORKDIR /source

COPY dotnetapp.csproj .
RUN dotnet restore

COPY . .
RUN dotnet publish --no-restore -o /app

FROM mcr.microsoft.com/dotnet/runtime

COPY --from=build --chown=app:app /app /app

WORKDIR /app

USER app

ENTRYPOINT ["./dotnetapp"]
