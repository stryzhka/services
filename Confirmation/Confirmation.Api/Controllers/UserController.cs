using Microsoft.AspNetCore.Mvc; 
using Confirmation.Application.Interactions.In;
using Confirmation.Application.Services.Interfaces;

namespace Confirmation.Api.Controllers;

[ApiController]
[Route("api/[controller]")]

public class UserController(IUserService userService) : ControllerBase
{
    [HttpPost]
    [ProducesResponseType(typeof(SignupResponse), StatusCodes.Status201Created)]
    [ProducesResponseType(StatusCodes.Status400BadRequest)]
    public async Task<SignupResponse> SignupAsync(SignupRequest request)
    {
        return await _userService.SignupAsync(request);
    }
    private readonly IUserService _userService = userService ;
}