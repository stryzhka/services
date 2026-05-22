using MongoDB.Bson;
using MongoDB.Bson.Serialization.Attributes;

namespace Confirmation.Core.Models;

public class User
{
    [BsonElement("_id")]
    [BsonGuidRepresentation(GuidRepresentation.Standard)]
    public Guid Id { get; set; }

    [BsonElement("confirmed_objects"), BsonRequired]
    public int ConfirmedObjects { get; set; } = 0;
}
