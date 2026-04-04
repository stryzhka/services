using Confirmation.Core.Events;

namespace Confirmation.Application.Services.Interfaces;

public interface IConfirmService
{
    Task ConfirmObjectAsync(GotConfirmRequestEvent @event, CancellationToken ct);
}