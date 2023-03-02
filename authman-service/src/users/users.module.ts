import * as redisStore from 'cache-manager-redis-store';
import { CacheModule, Module } from '@nestjs/common';
import { UsersService } from './users.service';
import { UsersController } from './users.controller';
import { JwtTokenService } from './jwt/jwt.service';
import { MongooseModule } from '@nestjs/mongoose';
import { User, UserSchema } from './database/schemas/user.schema';
import * as bcrypt from 'bcrypt';
import { MailerModule } from '@nestjs-modules/mailer';
import { MailService } from './mail/mail.service';
import { HandlebarsAdapter } from '@nestjs-modules/mailer/dist/adapters/handlebars.adapter';
import { join } from 'path';
import { EncryptService } from './encrypt/encrypt.service';
import * as dotenv from 'dotenv';
import { DBService } from './database/db.service';
dotenv.config();

@Module({
  imports: [
    MongooseModule.forFeatureAsync([
      {
        name: User.name,
        useFactory: () => {
          const schema = UserSchema;
          schema.pre('save', async function (next) {
            try {
              const salt = await bcrypt.genSalt(
                Number(process.env.SALT_ROUNDS),
              );
              const hash = await bcrypt.hash(this.password, salt);
              this.password = hash;
              return next();
            } catch (error) {
              next(error);
            }
          });
          return schema;
        },
      },
    ]),
    CacheModule.register({
      host: process.env.REDIS_HOST,
      port: process.env.REDIS_PORT,
      password: process.env.REDIS_PASSWORD,
      store: redisStore.redisStore,
    }),
    MailerModule.forRootAsync({
      useFactory: async () => ({
        transport: {
          host: process.env.SMTP_HOST,
          port: Number(process.env.SMTP_PORT),
          secure: true,
          auth: {
            user: process.env.SMTP_USER,
            pass: process.env.SMTP_PASS,
          },
          tls: {
            ciphers: 'SSLv3',
          },
        },
        defaults: {
          from: process.env.EMAIL,
        },
        template: {
          dir: join(__dirname, './mail/templates'),
          adapter: new HandlebarsAdapter(),
          options: {
            strict: false,
          },
        },
      }),
    }),
  ],
  controllers: [UsersController],
  providers: [
    UsersService,
    JwtTokenService,
    MailService,
    EncryptService,
    DBService,
  ],
})
export class UsersModule {}
