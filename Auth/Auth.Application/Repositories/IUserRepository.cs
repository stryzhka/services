using Auth.Core.Models;

namespace Auth.Application.Repositories;

public interface IUserRepository
{
    Task<User?> GetByNameAsync(string name);
    Task CreateAsync(User user);
}
