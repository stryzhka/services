namespace Confirmation.Application.Repositories;
using Confirmation.Core.Models;

public interface IUserRepository
{
    Task CreateAsync(User user);
    Task<User?> GetByIdAsync(Guid id);
    Task<User?> VerifyUser(string name, string password);
    Task<User?> IncrementConfirmedObjects(Guid id);
}