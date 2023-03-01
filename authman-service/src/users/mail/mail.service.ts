import { MailerService } from '@nestjs-modules/mailer';
import { Injectable } from '@nestjs/common';
import { EncryptService } from '../encrypt/encrypt.service';
import { JwtTokenService } from '../jwt/jwt.service';
import * as dotenv from 'dotenv';
dotenv.config();

@Injectable()
export class MailService {
  constructor(
    private mailerService: MailerService,
    private jwtService: JwtTokenService,
    private encryptService: EncryptService,
  ) {}

  async sendMail(email: string, payload: tokenPayload) {
    const token = await this.jwtService.generateEmailVerificationToken(payload);
    const encryptedToken = await this.encryptService.encrypt(token);
    try {
      await this.mailerService.sendMail({
        from: process.env.EMAIL,
        to: email,
        subject: `Welcome to Our World ${payload.username}`,
        template: './welcome.template.hbs',
        context: {
          ServiceName: 'Authman',
          VerificationLink: `http://localhost:3000/users/confirm-email/${encryptedToken}`,
        },
      });
    } catch (error) {
      /*
      Logger Should Handle this. 
      I dont want to throw an exception here,
      cos not sending a mail shouldnt break the flow.
      */
      console.log(error);
    }
  }
}
