using Confirmation.Application.Repositories;
using Confirmation.Application.Services.Interfaces;
using Confirmation.Core.Events;
using Confirmation.Core.Models;
using Microsoft.Extensions.Logging;

namespace Confirmation.Application.Services;

public class ConfirmService : IConfirmService
{
    private readonly ILogger<ConfirmService> _logger;
    private readonly IUserRepository _userRepository;
    
    public ConfirmService(ILogger<ConfirmService> logger, IUserRepository userRepository) 
    {
        _logger = logger;
        _userRepository = userRepository;
    }

    public async Task ConfirmObjectAsync(GotConfirmRequestEvent @event, CancellationToken ct)
    {
        _logger.LogInformation("Confirm Object got message: {UserId} {ObjectId}",  @event.UserId, @event.ObjectId);
        if (!Guid.TryParse(@event.UserId, out var id))
        {
            _logger.LogError("Invalid UserId: {UserId}", @event.UserId);
            return;
        }

        var user = await _userRepository.IncrementConfirmedObjects(id);

        if (user is null)
        {
            _logger.LogWarning("User {UserId} not found or some troubles", @event.UserId);
            return;
        }

        _logger.LogInformation("doing things...");
        await Task.CompletedTask;
    }
    
    
}