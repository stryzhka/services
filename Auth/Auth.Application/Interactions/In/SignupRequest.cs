using System.ComponentModel.DataAnnotations;

namespace Auth.Application.Interactions.In;

public record SignupRequest
{
    [Required, MaxLength(32)]
    public string Name { get; init; } = string.Empty;

    [Required, MaxLength(64)]
    public string Password { get; init; } = string.Empty;
}
