using Auth.Application.Interactions.Out;
using Auth.Core.Models;

namespace Auth.Application.Services.Interfaces;

public interface ITokenService
{
    TokenResponse GenerateToken(User user);
}
