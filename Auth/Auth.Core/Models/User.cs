using MongoDB.Bson;
using MongoDB.Bson.Serialization.Attributes;

namespace Auth.Core.Models;

public class User
{
    [BsonElement("_id")]
    [BsonGuidRepresentation(GuidRepresentation.Standard)]
    public Guid Id { get; set; }

    [BsonElement("name"), BsonRequired]
    public string Name { get; set; } = string.Empty;

    [BsonElement("password_hash"), BsonRequired]
    public string PasswordHash { get; set; } = string.Empty;
}
