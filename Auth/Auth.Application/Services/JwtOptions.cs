namespace Auth.Application.Services;

public class JwtOptions
{
    public const string SectionName = "Jwt";

    public string Key { get; set; } = string.Empty;
    public string Issuer { get; set; } = "auth-service";
    public string Audience { get; set; } = "internal-api";
    public int ExpiresMinutes { get; set; } = 60;
}
