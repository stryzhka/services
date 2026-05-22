using Confirmation.Infrastructure.Messaging;
using Microsoft.Extensions.Hosting;

namespace Confirmation.Infrastructure.Worker;

public class UserCreatedWorkerService : BackgroundService
{
    private readonly UserCreatedConsumerHandler _handler;

    public UserCreatedWorkerService(UserCreatedConsumerHandler handler)
    {
        _handler = handler;
    }

    protected override Task ExecuteAsync(CancellationToken ct)
    {
        return _handler.ConsumeAsync("users", ct);
    }
}
