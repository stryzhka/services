using MongoDB.Bson;
using MongoDB.Bson.Serialization.Attributes;

namespace Confirmation.Core.Models;

public class User
{
    [BsonElement("_id")]
    [BsonGuidRepresentation(GuidRepresentation.Standard)]
    public Guid Id { get; set; }
    [BsonElement("name"), BsonRequired]
    public string Name { get; set; }
    [BsonElement("password"), BsonRequired]
    public string Password { get; set; } = "dummy";
    [BsonElement("confirmed_objects"), BsonRequired]
    public int ConfirmedObjects { get; set; } = 0;
}

// public class CreateUser
// {
//     public string Name { get; set; }
//     public string Password { get; set; }
// }
//
// // public class UpdateUserModel
// // {
// //     public string Id { get; set; }
// //     public string Name { get; set; }
// // }
// //
// // public class DeleteUserModel
// // {
// //     public string Id { get; set; }
// // }
//
// public class GetUser
// {
//     public string Id { get; set; }
//     public string Name { get; set; }
//     public string Password { get; set; }
// }