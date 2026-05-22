using Confirmation.Application.Repositories;
using Confirmation.Core.Models;
using MongoDB.Driver;

namespace Confirmation.Infrastructure.Repositories;

public class MongoUserRepository : IUserRepository
{
    private readonly IMongoCollection<User> _collection;

    public MongoUserRepository(MongoClient client)
    {
        _collection = client.GetDatabase("confirmation").GetCollection<User>("user_stats");
    }

    public async Task<User?> IncrementConfirmedObjects(Guid id)
    {
        var filter = Builders<User>.Filter.Eq(u => u.Id, id);
        var update = Builders<User>.Update
            .SetOnInsert(u => u.Id, id)
            .Inc(u => u.ConfirmedObjects, 1);

        var options = new FindOneAndUpdateOptions<User>
        {
            ReturnDocument = ReturnDocument.After,
            IsUpsert = true,
        };

        return await _collection.FindOneAndUpdateAsync(filter, update, options);
    }

    public async Task CreateIfNotExistsAsync(Guid id)
    {
        var filter = Builders<User>.Filter.Eq(u => u.Id, id);
        var update = Builders<User>.Update
            .SetOnInsert(u => u.Id, id)
            .SetOnInsert(u => u.ConfirmedObjects, 0);

        await _collection.UpdateOneAsync(filter, update, new UpdateOptions { IsUpsert = true });
    }
}
