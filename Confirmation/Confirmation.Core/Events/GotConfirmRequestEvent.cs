using System.Text.Json.Serialization;

namespace Confirmation.Core.Events;

public class GotConfirmRequestEvent
{
    [JsonPropertyName("user_id")]
    public string UserId { get; set; }
    
    [JsonPropertyName("object_id")]
    public string ObjectId { get; set; }

}