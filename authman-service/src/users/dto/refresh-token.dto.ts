import { IsNotEmpty } from 'class-validator';

export class RefreshTokenDto {
  @IsNotEmpty({
    message: 'No User Token Present',
  })
  refreshToken: string;
}
