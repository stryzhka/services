using Microsoft.AspNetCore.Mvc; 
using Confirmation.Application.Interactions.In;
using Confirmation.Application.Interactions.Out;
using Confirmation.Application.Services.Interfaces;

namespace Confirmation.Api.Controllers;

[ApiController]
[Route("api/[controller]")]

public class UserController(IUserService userService) : ControllerBase
{
    [HttpPost]
    [Route("signup")]
    [ProducesResponseType(typeof(SignupResponse), StatusCodes.Status201Created)]
    [ProducesResponseType(StatusCodes.Status400BadRequest)]
    public async Task<SignupResponse> SignupAsync(SignupRequest request)
    {
        return await _userService.SignupAsync(request.Name, request.Password);
    }
    [HttpPost]
    [Route("signin")]
    [ProducesResponseType(typeof(TokenResponse), StatusCodes.Status200OK)]
    [ProducesResponseType(StatusCodes.Status400BadRequest)]
    public async Task<TokenResponse> SigninAsync(SigninRequest request)
    {
        return await _userService.SigninAsync(request.Name, request.Password);
    }
    private readonly IUserService _userService = userService ;
}