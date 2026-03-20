using Confirmation.Application.Repositories;
using Confirmation.Infrastructure.Repositories;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using MongoDB.Driver;

namespace Confirmation.Infrastructure;

public static class DependencyInjection
{
    public static IServiceCollection AddInfrastructure(this IServiceCollection services, IConfiguration configuration)
    {
        var connectionString = configuration.GetConnectionString("MongoDB");
        var mongoClient = new MongoClient(connectionString);
        
        services.AddSingleton<MongoClient>(mongoClient);
        services.AddScoped<IUserRepository, MongoUserRepository>();
        
        return services;
    }
}