using Confirmation.Infrastructure.Messaging;
using Microsoft.Extensions.Hosting;

namespace Confirmation.Infrastructure.Worker;

public class KafkaWorkerService : BackgroundService
{
    private readonly KafkaConsumerHandler _handler;

    public KafkaWorkerService(KafkaConsumerHandler handler)
    {
        _handler = handler;
    }

    protected override Task ExecuteAsync(CancellationToken ct)
    {
        return _handler.ConsumeAsync("obj-to-users", ct);
    }
}