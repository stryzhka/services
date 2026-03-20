using System.Security.Claims;
using Confirmation.Application.Interactions.In;
using Confirmation.Application.Interactions.Out;
using Confirmation.Core.Models;

namespace Confirmation.Application.Services.Interfaces;

public interface ITokenService
{
    TokenResponse GenerateToken(User user);
    // Task<ClaimsPrincipal?>  GetClaimsPrincipalAsync(User user); 
}