using Auth.Core.Models;

namespace Auth.Application.Interactions.Out;

public record SignupResponse
{
    public SignupResponse(User user)
    {
        Id = user.Id;
        Name = user.Name;
    }

    public Guid Id { get; init; }
    public string Name { get; init; }
}
