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
        var collection = _client.GetDatabase("words").GetCollection<User>("users");
        var filter = Builders<User>.Filter.Eq("_id", id.ToString());
    
        var user = await collection.Find(filter).FirstOrDefaultAsync();
        return user;
    }
    
    public async Task<User> VerifyUser(string name, string password)
    {
        var collection = _client.GetDatabase("words").GetCollection<User>("users");
        
        var filter = Builders<User>.Filter.And(
                Builders<User>.Filter.Eq("name", name.ToString()),
                Builders<User>.Filter.Eq("password", password.ToString())        
            );
        var user = await collection.Find(filter).FirstOrDefaultAsync();
        Console.WriteLine(user);
        return user;
    }
}