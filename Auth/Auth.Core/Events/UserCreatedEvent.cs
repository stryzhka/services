using System.Text.Json.Serialization;

namespace Auth.Core.Events;

public class UserCreatedEvent
{
    [JsonPropertyName("user_id")]
    public string UserId { get; set; } = string.Empty;

    [JsonPropertyName("name")]
    public string Name { get; set; } = string.Empty;
}
