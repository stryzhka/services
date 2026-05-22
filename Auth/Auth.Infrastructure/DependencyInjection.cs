using Auth.Application.Repositories;
using Auth.Application.Services.Interfaces;
using Auth.Infrastructure.Messaging;
using Auth.Infrastructure.Repositories;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using MongoDB.Bson;
using MongoDB.Bson.Serialization;
using MongoDB.Bson.Serialization.Serializers;
using MongoDB.Driver;

namespace Auth.Infrastructure;

public static class DependencyInjection
{
    private static bool _guidSerializerRegistered;

    public static IServiceCollection AddInfrastructure(this IServiceCollection services, IConfiguration configuration)
    {
        if (!_guidSerializerRegistered)
        {
            BsonSerializer.RegisterSerializer(new GuidSerializer(GuidRepresentation.Standard));
            _guidSerializerRegistered = true;
        }

        var connectionString = configuration.GetConnectionString("MongoDB");
        var mongoClient = new MongoClient(connectionString);

        services.AddSingleton(mongoClient);
        services.AddScoped<IUserRepository, MongoUserRepository>();
        services.AddSingleton<IEventPublisher, KafkaEventPublisher>();

        return services;
    }
}
