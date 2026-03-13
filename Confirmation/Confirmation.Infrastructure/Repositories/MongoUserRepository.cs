using Confirmation.Application.Repositories;
using Confirmation.Core.Models;
using MongoDB.Bson;
using MongoDB.Driver;

namespace Confirmation.Infrastructure.Repositories;

public class MongoUserRepository(MongoClient client) : IUserRepository
{
    private readonly MongoClient _client = client;
    public async Task CreateAsync(User user)
    {
        var collection = _client.GetDatabase("words").GetCollection<BsonDocument>("users");
        await collection.InsertOneAsync(user.ToBsonDocument());
        
    }

    public async Task<User> GetByIdAsync(Guid id)
    {
        var collection = _client.GetDatabase("confirmation").GetCollection<User>("users");
        var filter = Builders<User>.Filter.Eq("_id", id.ToString());
    
        var user = await collection.Find(filter).FirstOrDefaultAsync();
        return user;
    }
}