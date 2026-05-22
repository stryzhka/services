using Auth.Application.Interactions.Out;
using Auth.Application.Repositories;
using Auth.Application.Services.Interfaces;
using Auth.Core.Events;
using Auth.Core.Models;

namespace Auth.Application.Services;

public class UserService : IUserService
{
    private readonly IUserRepository _repository;
    private readonly ITokenService _tokenService;
    private readonly IEventPublisher _publisher;

    public UserService(IUserRepository repository, ITokenService tokenService, IEventPublisher publisher)
    {
        _repository = repository;
        _tokenService = tokenService;
        _publisher = publisher;
    }

    public async Task<SignupResponse> SignupAsync(string name, string password)
    {
        var existing = await _repository.GetByNameAsync(name);
        if (existing is not null)
            throw new InvalidOperationException("user already exists");

        var user = new User
        {
            Id = Guid.NewGuid(),
            Name = name,
            PasswordHash = PasswordHasher.Hash(password),
        };

        await _repository.CreateAsync(user);

        await _publisher.PublishAsync("users", new UserCreatedEvent
        {
            UserId = user.Id.ToString(),
            Name = user.Name,
        }, CancellationToken.None);

        return new SignupResponse(user);
    }

    public async Task<TokenResponse?> SigninAsync(string name, string password)
    {
        var user = await _repository.GetByNameAsync(name);
        if (user is null) return null;
        if (!PasswordHasher.Verify(password, user.PasswordHash)) return null;

        return _tokenService.GenerateToken(user);
    }
}
