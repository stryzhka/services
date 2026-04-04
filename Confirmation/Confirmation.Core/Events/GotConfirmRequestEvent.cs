using System.Text.Json.Serialization;

namespace Confirmation.Core.Events;

public class GotConfirmRequestEvent
{
    [JsonPropertyName("user_id")]
    public string UserId { get; set; }
}