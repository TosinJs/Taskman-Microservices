import { IsNotEmpty } from 'class-validator';

export class LoginUserDto {
  @IsNotEmpty({
    message: 'Invalid Login Credentials',
  })
  username: string;

  @IsNotEmpty({
    message: 'Invalid Login Credentials',
  })
  password: string;
}
