FROM mcr.microsoft.com/dotnet/aspnet:10.0 AS base
USER $APP_UID
WORKDIR /app
EXPOSE 8082
EXPOSE 8083
EXPOSE 5016

FROM mcr.microsoft.com/dotnet/sdk:10.0 AS build
ARG BUILD_CONFIGURATION=Release
WORKDIR /src
COPY ["Confirmation.Api/Confirmation.Api.csproj", "Confirmation.Api/"]
RUN dotnet restore "Confirmation.Api/Confirmation.Api.csproj"
COPY . .
WORKDIR "/src/Confirmation.Api"
RUN dotnet build "./Confirmation.Api.csproj" -c $BUILD_CONFIGURATION -o /app/build

FROM build AS publish
ARG BUILD_CONFIGURATION=Release
RUN dotnet publish "./Confirmation.Api.csproj" -c $BUILD_CONFIGURATION -o /app/publish /p:UseAppHost=false

FROM base AS final
WORKDIR /app
COPY --from=publish /app/publish .
ENTRYPOINT ["dotnet", "Confirmation.Api.dll"]
