using Confirmation.Application.Interactions.In;
using Confirmation.Application.Repositories;
using Confirmation.Application.Services.Interfaces;
using Confirmation.Core.Models;

namespace Confirmation.Application.Services;

public class UserService(IUserRepository repository) : IUserService
{
     private readonly IUserRepository _repository = repository;
     public async Task<SignupResponse> SignupAsync(SignupRequest request)
     {
          User newUser = new User();
          newUser.Name = request.Name;
          newUser.Password = request.Password;
          newUser.Id = Guid.NewGuid();
          
          await _repository.CreateAsync(newUser);
          var resp = new SignupResponse(newUser);
          return resp;
          
     }
     
     public async Task<SigninResponse?> SigninAsync(SigninRequest request)
     {
          return null;
     }
}