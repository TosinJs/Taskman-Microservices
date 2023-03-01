import { IsNotEmpty, IsString } from 'class-validator';

export class EmailTokenDto {
  @IsNotEmpty({
    message: 'No User Token Present',
  })
  @IsString()
  token: string;
}
