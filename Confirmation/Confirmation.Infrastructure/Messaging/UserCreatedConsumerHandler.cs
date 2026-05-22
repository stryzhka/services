using System.Diagnostics;
using System.Text.Json;
using Confirmation.Application.Repositories;
using Confirmation.Core.Events;
using Confluent.Kafka;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Logging;

namespace Confirmation.Infrastructure.Messaging;

public class UserCreatedConsumerHandler
{
    private readonly IServiceScopeFactory _scopeFactory;
    private readonly ILogger<UserCreatedConsumerHandler> _logger;
    private readonly IConsumer<string, string> _consumer;

    public UserCreatedConsumerHandler(
        IServiceScopeFactory scopeFactory,
        IConfiguration config,
        ILogger<UserCreatedConsumerHandler> logger)
    {
        _scopeFactory = scopeFactory;
        _logger = logger;
        var consumerConfig = new ConsumerConfig
        {
            BootstrapServers = config["Kafka:BootstrapServers"],
            GroupId = config["Kafka:UsersGroupId"] ?? "confirmation-users",
            AutoOffsetReset = AutoOffsetReset.Earliest,
            EnableAutoCommit = false,
        };
        _consumer = new ConsumerBuilder<string, string>(consumerConfig).Build();
    }

    public async Task ConsumeAsync(string topic, CancellationToken ct)
    {
        _consumer.Subscribe(topic);
        while (!ct.IsCancellationRequested)
        {
            ConsumeResult<string, string>? result;
            try
            {
                result = _consumer.Consume(ct);
            }
            catch (OperationCanceledException)
            {
                break;
            }
            catch (ConsumeException ex)
            {
                _logger.LogWarning(ex, "Consume error on topic {Topic}, retrying", topic);
                await Task.Delay(TimeSpan.FromSeconds(2), ct);
                continue;
            }

            if (result is null) continue;

            var sw = Stopwatch.StartNew();
            try
            {
                var @event = JsonSerializer.Deserialize<UserCreatedEvent>(result.Message.Value);
                if (@event is null || !Guid.TryParse(@event.UserId, out var userId))
                {
                    _logger.LogWarning("Invalid user-created event payload: {Value}", result.Message.Value);
                    _consumer.Commit(result);
                    continue;
                }

                using var scope = _scopeFactory.CreateScope();
                var repo = scope.ServiceProvider.GetRequiredService<IUserRepository>();
                await repo.CreateIfNotExistsAsync(userId);
                _logger.LogInformation("Provisioned user_stats for {UserId} ({Name})", userId, @event.Name);
                _consumer.Commit(result);
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, "Failed to process user-created event");
            }
            finally
            {
                KafkaMessagingMetrics.ObserveConsumerLatency(result.Topic, sw.Elapsed.TotalSeconds);
            }
        }
    }
}
