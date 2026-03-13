using Confirmation.Application.Interactions.In;
using Confirmation.Core.Models;

namespace Confirmation.Application.Services.Interfaces;

public interface IUserService
{
    Task<SignupResponse?> SignupAsync(SignupRequest request);
    Task<SigninResponse?> SigninAsync(SigninRequest request);
}