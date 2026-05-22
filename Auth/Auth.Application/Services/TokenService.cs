using System.Text;
using Auth.Application.Interactions.Out;
using Auth.Application.Services.Interfaces;
using Auth.Core.Models;
using Microsoft.Extensions.Options;
using Microsoft.IdentityModel.JsonWebTokens;
using Microsoft.IdentityModel.Tokens;

namespace Auth.Application.Services;

public class TokenService : ITokenService
{
    private readonly JsonWebTokenHandler _handler = new();
    private readonly JwtOptions _options;

    public TokenService(IOptions<JwtOptions> options)
    {
        _options = options.Value;
        if (string.IsNullOrWhiteSpace(_options.Key))
            throw new InvalidOperationException("Jwt:Key is not configured");
    }

    public TokenResponse GenerateToken(User user)
    {
        var expiresAt = DateTime.UtcNow.AddMinutes(_options.ExpiresMinutes);

        var claims = new Dictionary<string, object>
        {
            [JwtRegisteredClaimNames.Sub] = user.Name,
            ["uid"] = user.Id,
        };

        var key = new SymmetricSecurityKey(Encoding.UTF8.GetBytes(_options.Key));
        var descriptor = new SecurityTokenDescriptor
        {
            Claims = claims,
            Issuer = _options.Issuer,
            Audience = _options.Audience,
            Expires = expiresAt,
            SigningCredentials = new SigningCredentials(key, SecurityAlgorithms.HmacSha256)
        };

        var token = _handler.CreateToken(descriptor);
        return new TokenResponse { Token = token, ExpiresAt = expiresAt };
    }
}
