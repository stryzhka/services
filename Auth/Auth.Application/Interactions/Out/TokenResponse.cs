namespace Auth.Application.Interactions.Out;

public record TokenResponse
{
    public required string Token { get; init; }
    public required DateTime ExpiresAt { get; init; }
}
