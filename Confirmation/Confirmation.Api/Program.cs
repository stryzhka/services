using Confirmation.Application;
using Confirmation.Application.Services;
using Confirmation.Application.Services.Interfaces;
using Confirmation.Infrastructure;
using Confirmation.Infrastructure.Messaging;
using Confirmation.Infrastructure.Worker;
using Microsoft.AspNetCore.Mvc;
using Prometheus;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddControllers();
builder.Services.AddEndpointsApiExplorer();

builder.Services.AddInfrastructure(builder.Configuration);
builder.Services.AddApplication(builder.Configuration);
builder.Services.AddScoped<IConfirmService, ConfirmService>();
builder.Services.AddSingleton<KafkaConsumerHandler>();
builder.Services.AddHostedService<KafkaWorkerService>();
builder.Services.AddSingleton<UserCreatedConsumerHandler>();
builder.Services.AddHostedService<UserCreatedWorkerService>();
builder.Services.AddSingleton<IEventPublisher, KafkaEventPublisher>();

var app = builder.Build();

app.UseExceptionHandler(error => error.Run(async context =>
{
    context.Response.StatusCode = StatusCodes.Status500InternalServerError;
    context.Response.ContentType = "application/problem+json";

    await context.Response.WriteAsJsonAsync(new ProblemDetails
    {
        Status = StatusCodes.Status500InternalServerError,
        Title = "Internal Server Error"
    });
}));

app.MapControllers();
app.MapMetrics();

app.Run();
