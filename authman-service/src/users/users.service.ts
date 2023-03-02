import { HttpException, HttpStatus, Injectable } from '@nestjs/common';
import { CreateUserDto } from './dto/create-user.dto';
import { EmailTokenDto } from './dto/email-token.dto';
import { LoginUserDto } from './dto/login-user.dto';
import { RefreshTokenDto } from './dto/refresh-token.dto';
import { RefreshUserPasswordDto } from './dto/resfresh-password.dto';
import { EncryptService } from './encrypt/encrypt.service';
import { JwtTokenService } from './jwt/jwt.service';
import { MailService } from './mail/mail.service';
import { DBService } from './database/db.service';

@Injectable()
export class UsersService {
  constructor(
    private jwtService: JwtTokenService,
    private mailService: MailService,
    private encryptService: EncryptService,
    private dbService: DBService,
  ) {}

  private async genToken(payload): Promise<responseObj> {
    const idToken = await this.jwtService.generateIdToken(payload, '1d');
    const refreshToken = await this.jwtService.generateRefreshToken(
      payload,
      '365d',
    );
    const accessToken = await this.jwtService.generateAuthToken(payload, '25m');
    return { idToken, refreshToken, accessToken, tokenType: 'Bearer' };
  }

  async create(createUserDto: CreateUserDto) {
    const user = await this.dbService.createUser(createUserDto);

    const payload: tokenPayload = {
      id: user.id,
      username: user.username,
      email: user.email,
      isVerified: false,
    };
    await this.mailService.sendMail(user.email, payload);
    const idToken = await this.jwtService.generateIdToken(payload, '1d');
    return { idToken, tokenType: 'Bearer' };
  }

  async login(loginUserDto: LoginUserDto): Promise<responseObj> {
    const user = await this.dbService.findByUsername(loginUserDto.username);
    if (!user) {
      throw new HttpException(
        'Invalid User Credentials',
        HttpStatus.BAD_REQUEST,
      );
    }
    if (!(await user.comparePassword(loginUserDto.password))) {
      throw new HttpException(
        'Invalid User Credentials',
        HttpStatus.BAD_REQUEST,
      );
    }
    const payload: tokenPayload = {
      id: user.id,
      username: user.username,
      email: user.email,
      isVerified: user.emailVerified,
    };
    return this.genToken(payload);
  }

  async refresh(refreshToken: RefreshTokenDto) {
    try {
      const userTokenPayload = await this.jwtService.verifyRefreshToken(
        refreshToken.refreshToken,
      );
      return this.genToken(userTokenPayload);
    } catch (error) {
      if (error.message == 'Invalid Token Credentials') {
        throw new HttpException(
          'Invalid User Credentials',
          HttpStatus.BAD_REQUEST,
        );
      }
      throw new error();
    }
    /**
     * Whats Next,
     * Add refresh Token to Redis Cache
     * And Update the Refresh Handler
     * Create a Database to store the authTokens and their credentials
     * Send a reset mail to a user that allows the user to reset their password
     * Create a verifiable link that a user recieves to reset their main
     * Fix the sending of error messages from the class-validator package
     */
  }

  async confirmEmail(inputToken: EmailTokenDto) {
    const token = await this.encryptService.decrypt(inputToken.token);
    //dont proceed if theres no token
    const userTokenPayload = await this.jwtService.verifyRefreshToken(token);
    await this.dbService.findByIdAndUpdate(userTokenPayload.id, {
      emailVerified: true,
    });
    return this.genToken(userTokenPayload);
  }

  async reset(refreshUserPasswordDto: RefreshUserPasswordDto) {
    console.log(refreshUserPasswordDto.username);
    return 'password has been reset';
  }
}
