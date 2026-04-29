// using Confirmation.Api.Extensions;
using Confirmation.Application;
using Confirmation.Application.Services;
using Confirmation.Application.Services.Interfaces;
using Confirmation.Infrastructure;
using Confirmation.Infrastructure.Messaging;
using Confirmation.Infrastructure.Worker;
// using Confirmation.Infrastructure.Data;
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
builder.Services.AddSingleton<IEventPublisher, KafkaEventPublisher>();
// builder.Services.AddOpenApi();

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


app.UseHttpsRedirection();
// app.UseAuthentication();
// app.UseAuthorization();
app.MapControllers();
app.MapMetrics();

app.Run();

// // Configure the HTTP request pipeline.
// if (app.Environment.IsDevelopment())
// {
//     app.MapOpenApi();
// }
//
// app.UseHttpsRedirection();
//
//
// app.Run();
