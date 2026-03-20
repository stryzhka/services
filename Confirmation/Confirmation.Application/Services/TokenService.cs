using Confirmation.Application.Interactions.Out;
using Confirmation.Application.Services.Interfaces;
using Confirmation.Core.Models;
using Microsoft.IdentityModel.JsonWebTokens;
using Microsoft.IdentityModel.Tokens;  // SecurityTokenDescriptor lives here
using System.Security.Claims;
using System.Text; // if you need ClaimsIdentity later
// using Confirmation.Core.Models.Authentication;
// using Microsoft.IdentityModel.Tokens;
// using Microsoft.IdentityModel.JsonWebTokens;
// using System.Security.Claims;
// using System.Security.Cryptography;
// using System.Text;
// using Confirmation.Application.Interactions.In;
// using Confirmation.Application.Services.Interfaces;
// using Confirmation.Core.Models;

namespace Confirmation.Application.Services;

public class TokenService : ITokenService
{
    private JsonWebTokenHandler _handler;
    public TokenService()
    {
        _handler = new JsonWebTokenHandler();
    }

    public TokenResponse GenerateToken(User user)
    {
        var claims = new Dictionary<string, object>
        {
            [JwtRegisteredClaimNames.Sub] = user.Name,
            ["uid"] = user.Id,
        };
        
        var key = new SymmetricSecurityKey(Encoding.UTF8.GetBytes("super-secret-stab-fuck-suka-kogda-eto-konchitsya"));
        var tokenDescriptor = new SecurityTokenDescriptor
        {
            Claims = claims,
            Expires = DateTime.UtcNow.AddMinutes(60),
            SigningCredentials = new SigningCredentials(key, SecurityAlgorithms.HmacSha256)
        };
        string token = _handler.CreateToken(tokenDescriptor);
        // Console.WriteLine(token);
        return new TokenResponse { Token = token };
    }

}