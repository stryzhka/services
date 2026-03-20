using Confirmation.Application.Interactions.In;
using Confirmation.Application.Interactions.Out;
using Confirmation.Core.Models;

namespace Confirmation.Application.Services.Interfaces;

public interface IUserService
{
    Task<SignupResponse?> SignupAsync(string name, string password);
    Task<TokenResponse?> SigninAsync(string name, string password);
}