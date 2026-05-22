using Auth.Application.Repositories;
using Auth.Core.Models;
using MongoDB.Driver;

namespace Auth.Infrastructure.Repositories;

public class MongoUserRepository : IUserRepository
{
    private readonly IMongoCollection<User> _collection;

    public MongoUserRepository(MongoClient client)
    {
        _collection = client.GetDatabase("auth").GetCollection<User>("users");
    }

    public async Task<User?> GetByNameAsync(string name)
    {
        var filter = Builders<User>.Filter.Eq(u => u.Name, name);
        return await _collection.Find(filter).FirstOrDefaultAsync();
    }

    public async Task CreateAsync(User user)
    {
        await _collection.InsertOneAsync(user);
    }
}
