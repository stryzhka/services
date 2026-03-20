using System.ComponentModel.DataAnnotations;
namespace Confirmation.Application.Interactions.In;

public record SigninRequest
{
    [Required, MaxLength(32)]
    public string Name { get; init; }
    
    [Required, MaxLength(32)]
    public string Password  { get; init; }
    
    
}