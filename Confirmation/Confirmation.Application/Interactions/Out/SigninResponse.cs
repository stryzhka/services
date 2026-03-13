using System.ComponentModel.DataAnnotations;
using Confirmation.Core.Models;

namespace Confirmation.Application.Interactions.In;

public record SigninResponse
{
    public SigninResponse(User user)
    {
        Id = user.Id;
        Name = user.Name;
        ConfirmedObjects = user.ConfirmedObjects;
    }
    public Guid Id { get; init; }
    public string Name { get; init; }
    public int ConfirmedObjects { get; init; }
}