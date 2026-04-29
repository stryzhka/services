using System.Diagnostics;
using System.Text.Json;
using Confirmation.Application.Services.Interfaces;
using Confirmation.Core.Events;
using Confluent.Kafka;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Logging;

namespace Confirmation.Infrastructure.Messaging;

public class KafkaConsumerHandler
{
    private readonly IServiceScopeFactory _scopeFactory;
    private readonly IConsumer<string, string> _consumer;
    
    public KafkaConsumerHandler(IServiceScopeFactory scopeFactory, IConfiguration config)
    {
        _scopeFactory = scopeFactory;
        var consumerConfig = new ConsumerConfig
        {
            BootstrapServers = config["Kafka:BootstrapServers"],
            GroupId = config["Kafka:GroupId"],
            AutoOffsetReset = AutoOffsetReset.Earliest,
            EnableAutoCommit = false, //а вообще иди наху   
        };
        _consumer = new ConsumerBuilder<string, string>(consumerConfig).Build();
    }

    public async Task ConsumeAsync(string topic, CancellationToken ct)
    {
        _consumer.Subscribe(topic);
        while (!ct.IsCancellationRequested)
        {
            var result = _consumer.Consume(ct);
            var sw = Stopwatch.StartNew();
            try
            {
                var @event = JsonSerializer.Deserialize<GotConfirmRequestEvent>(result.Message.Value);
                using var scope = _scopeFactory.CreateScope();
                var confirmService = scope.ServiceProvider.GetRequiredService<IConfirmService>();
                await confirmService.ConfirmObjectAsync(@event, ct);
                _consumer.Commit(result);
            }
            finally
            {
                KafkaMessagingMetrics.ObserveConsumerLatency(result.Topic, sw.Elapsed.TotalSeconds);
            }
        }
    }
}