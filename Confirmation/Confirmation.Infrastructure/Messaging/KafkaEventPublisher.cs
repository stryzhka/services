using System.Diagnostics;
using System.Text.Json;
using Confirmation.Application.Services.Interfaces;
using Confluent.Kafka;
using Microsoft.Extensions.Configuration;

namespace Confirmation.Infrastructure.Messaging;

public class KafkaEventPublisher : IEventPublisher, IDisposable
{
    private readonly IProducer<string, string> _producer;

    public KafkaEventPublisher(IConfiguration config)
    {
        var producerConfig = new ProducerConfig
        {
            BootstrapServers = config["Kafka:BootstrapServers"]
        };
        _producer = new ProducerBuilder<string, string>(producerConfig).Build();
    }

    public async Task PublishAsync<T>(string topic, T message, CancellationToken ct)
    {
        var json = JsonSerializer.Serialize(message);
        var sw = Stopwatch.StartNew();
        try
        {
            await _producer.ProduceAsync(topic, new Message<string, string>
            {
                Key = Guid.NewGuid().ToString(),
                Value = json
            }, ct);
        }
        finally
        {
            KafkaMessagingMetrics.ObserveProducerLatency(topic, sw.Elapsed.TotalSeconds);
        }
    }

    public void Dispose() => _producer.Dispose();
}