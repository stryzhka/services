using Confirmation.Application.Interactions.In;
using Confirmation.Application.Interactions.Out;
using Confirmation.Application.Repositories;
using Confirmation.Application.Services.Interfaces;
using Confirmation.Core.Models;

namespace Confirmation.Application.Services;

public class UserService(IUserRepository repository, ITokenService tokenService) : IUserService
{
     private readonly IUserRepository _repository = repository;
     private readonly ITokenService _tokenService = tokenService;
     public async Task<SignupResponse> SignupAsync(string name, string password)
     {
          User newUser = new User();
          newUser.Name = name;
          newUser.Password = password;
          newUser.Id = Guid.NewGuid();
          
          await _repository.CreateAsync(newUser);
          var resp = new SignupResponse(newUser);
          return resp;
     }
     
     public async Task<TokenResponse?> SigninAsync(string name, string password)
     {
          User user = await _repository.VerifyUser(name, password);
          if  (user == null) return null;
          return _tokenService.GenerateToken(user);
     }
}