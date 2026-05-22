using Auth.Application.Interactions.Out;

namespace Auth.Application.Services.Interfaces;

public interface IUserService
{
    Task<SignupResponse> SignupAsync(string name, string password);
    Task<TokenResponse?> SigninAsync(string name, string password);
}
