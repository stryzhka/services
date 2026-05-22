using Auth.Application.Interactions.In;
using Auth.Application.Interactions.Out;
using Auth.Application.Services.Interfaces;
using Microsoft.AspNetCore.Mvc;

namespace Auth.Api.Controllers;

[ApiController]
[Route("api/[controller]")]
public class AuthController : ControllerBase
{
    private readonly IUserService _userService;

    public AuthController(IUserService userService)
    {
        _userService = userService;
    }

    [HttpPost("signup")]
    [ProducesResponseType(typeof(SignupResponse), StatusCodes.Status201Created)]
    [ProducesResponseType(StatusCodes.Status400BadRequest)]
    public async Task<IActionResult> SignupAsync([FromBody] SignupRequest request)
    {
        try
        {
            var response = await _userService.SignupAsync(request.Name, request.Password);
            return StatusCode(StatusCodes.Status201Created, response);
        }
        catch (InvalidOperationException ex)
        {
            return BadRequest(new { error = ex.Message });
        }
    }

    [HttpPost("signin")]
    [ProducesResponseType(typeof(TokenResponse), StatusCodes.Status200OK)]
    [ProducesResponseType(StatusCodes.Status401Unauthorized)]
    public async Task<IActionResult> SigninAsync([FromBody] SigninRequest request)
    {
        var token = await _userService.SigninAsync(request.Name, request.Password);
        if (token is null) return Unauthorized(new { error = "invalid credentials" });
        return Ok(token);
    }
}
