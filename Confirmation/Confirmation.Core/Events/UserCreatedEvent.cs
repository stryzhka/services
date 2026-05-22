using System.Text.Json.Serialization;

namespace Confirmation.Core.Events;

public class UserCreatedEvent
{
    [JsonPropertyName("user_id")]
    public string UserId { get; set; } = string.Empty;

    [JsonPropertyName("name")]
    public string Name { get; set; } = string.Empty;
}
