import { IsNotEmpty, IsString } from 'class-validator';

export class RefreshUserPasswordDto {
  @IsNotEmpty()
  @IsString({
    message: 'Invalid User Credentials',
  })
  username: string;
}
