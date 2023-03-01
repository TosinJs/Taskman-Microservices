import { Controller, Post, Body, Get, Param, Redirect } from '@nestjs/common';
import { UsersService } from './users.service';
import { CreateUserDto } from './dto/create-user.dto';
import { LoginUserDto } from './dto/login-user.dto';
import { RefreshTokenDto } from './dto/refresh-token.dto';
import { RefreshUserPasswordDto } from './dto/resfresh-password.dto';
import { EmailTokenDto } from './dto/email-token.dto';

@Controller('users')
export class UsersController {
  constructor(private readonly usersService: UsersService) {}

  @Post('/signup')
  async create(@Body() createUserDto: CreateUserDto) {
    return this.usersService.create(createUserDto);
  }

  @Post('/login')
  async login(@Body() loginUserDto: LoginUserDto) {
    return this.usersService.login(loginUserDto);
  }

  @Post('/refresh')
  async refresh(@Body() refreshToken: RefreshTokenDto) {
    return this.usersService.refresh(refreshToken);
  }

  @Get('/reset-password')
  async reset(@Body() refreshUserPasswordDto: RefreshUserPasswordDto) {
    return this.usersService.reset(refreshUserPasswordDto);
  }

  @Get('/confirm-email/:token')
  @Redirect('/', 301)
  async confirmEmail(@Param() token: EmailTokenDto) {
    return this.usersService.confirmEmail(token);
  }
}
