using Confirmation.Application.Services.Interfaces;
using Confirmation.Core.Events;
using Microsoft.Extensions.Logging;

namespace Confirmation.Application.Services;

public class ConfirmService : IConfirmService
{
    private readonly ILogger<ConfirmService> _logger;
    
    public ConfirmService(ILogger<ConfirmService> logger) 
    {
        _logger = logger;
    }

    public async Task ConfirmObjectAsync(GotConfirmRequestEvent @event, CancellationToken ct)
    {
        _logger.LogInformation("Confirm Object got message: {UserId}",  @event.UserId);
        // Console.WriteLine("Confirm Object got message " +  @event.UserId);
        await Task.CompletedTask;
    }
    
    
}