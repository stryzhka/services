using Auth.Application;
using Auth.Infrastructure;
using Microsoft.AspNetCore.Mvc;
using Microsoft.OpenApi;
using Prometheus;
using Scalar.AspNetCore;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddControllers();
builder.Services.AddOpenApi(options =>
{
    options.AddDocumentTransformer((document, _, _) =>
    {
        var publicUrl = builder.Configuration["OpenApi:ServerUrl"];
        if (!string.IsNullOrWhiteSpace(publicUrl))
        {
            document.Servers = new List<OpenApiServer> { new() { Url = publicUrl } };
        }
        return Task.CompletedTask;
    });
});

builder.Services.AddInfrastructure(builder.Configuration);
builder.Services.AddApplication(builder.Configuration);

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

app.MapOpenApi();
app.MapScalarApiReference("/swagger", o => o
    .WithTitle("Auth API"));

app.MapControllers();
app.MapMetrics();

app.Run();
