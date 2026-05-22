namespace Confirmation.Application.Repositories;
using Confirmation.Core.Models;

public interface IUserRepository
{
    Task<User?> IncrementConfirmedObjects(Guid id);
    Task CreateIfNotExistsAsync(Guid id);
}
