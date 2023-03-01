import { IsNotEmpty, Length, IsOptional, IsEmail } from 'class-validator';

export class CreateUserDto {
  @IsNotEmpty()
  @Length(4, 35, {
    message: 'Username should be between 4 and 25 Characters long',
  })
  username: string;

  @IsNotEmpty()
  @Length(5, 200, {
    message: 'Password should be at least 5 characters long',
  })
  password: string;

  @IsNotEmpty()
  @IsEmail()
  @Length(5, 200, {
    message: 'Invalid Email',
  })
  email: string;

  @IsOptional()
  avatar: string;
}
