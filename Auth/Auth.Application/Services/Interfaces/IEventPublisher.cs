namespace Auth.Application.Services.Interfaces;

public interface IEventPublisher
{
    Task PublishAsync<T>(string topic, T message, CancellationToken ct);
}
