using Prometheus;

namespace Confirmation.Infrastructure.Messaging;

internal static class KafkaMessagingMetrics
{
    private static readonly Histogram ProducerLatencySeconds = Metrics.CreateHistogram(
        "confirmation_kafka_producer_latency_seconds",
        "Latency of producing a message to Kafka (serialize + ProduceAsync).",
        new HistogramConfiguration
        {
            LabelNames = ["topic"],
            Buckets = [0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10]
        });

    private static readonly Histogram ConsumerLatencySeconds = Metrics.CreateHistogram(
        "confirmation_kafka_consumer_latency_seconds",
        "Latency of handling one consumed message (deserialize, business logic, commit).",
        new HistogramConfiguration
        {
            LabelNames = ["topic"],
            Buckets = [0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10]
        });

    public static void ObserveProducerLatency(string topic, double elapsedSeconds) =>
        ProducerLatencySeconds.WithLabels(topic).Observe(elapsedSeconds);

    public static void ObserveConsumerLatency(string topic, double elapsedSeconds) =>
        ConsumerLatencySeconds.WithLabels(topic).Observe(elapsedSeconds);
}
