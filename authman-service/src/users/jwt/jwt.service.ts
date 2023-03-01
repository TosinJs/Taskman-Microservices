import { CACHE_MANAGER, Inject, Injectable } from '@nestjs/common';
import { Cache } from 'cache-manager';
import * as jwt from 'jsonwebtoken';
import * as dotenv from 'dotenv';
dotenv.config();

@Injectable()
export class JwtTokenService {
  constructor(@Inject(CACHE_MANAGER) private cacheManager: Cache) {}

  generateIdToken(payload: tokenPayload, expiryTime: string): Promise<string> {
    return new Promise((resolve, reject) => {
      jwt.sign(
        payload,
        process.env.JWT_SECRET,
        {
          expiresIn: expiryTime,
          issuer: 'RilHomie',
        },
        (err, token) => {
          if (err) reject(err);
          resolve(token);
        },
      );
    });
  }

  generateRefreshToken(
    payload: tokenPayload,
    expiryTime: string,
  ): Promise<string> {
    return new Promise((resolve, reject) => {
      jwt.sign(
        payload,
        process.env.JWT_SECRET,
        {
          expiresIn: expiryTime,
          issuer: 'RilHomie',
        },
        async (err, token) => {
          if (err) reject(err);
          await this.cacheManager.set(payload.id, token, 3600 * 24 * 12);
          resolve(token);
        },
      );
    });
  }

  generateEmailVerificationToken(payload: tokenPayload): Promise<string> {
    return new Promise((resolve, reject) => {
      jwt.sign(
        payload,
        process.env.JWT_SECRET,
        {
          expiresIn: '1d',
          issuer: 'RilHomie',
        },
        async (err, token) => {
          if (err) reject(err);
          await this.cacheManager.set(payload.id, token, 3600 * 24);
          resolve(token);
        },
      );
    });
  }

  generateAuthToken(
    payload: tokenPayload,
    expiryTime: string,
  ): Promise<string> {
    return new Promise((resolve, reject) => {
      jwt.sign(
        payload,
        process.env.JWT_SECRET,
        {
          expiresIn: expiryTime,
          issuer: 'RilHomie',
        },
        (err, token) => {
          if (err) reject(err);
          resolve(token);
        },
      );
    });
  }

  verifyRefreshToken(refreshToken: string): Promise<tokenPayload> {
    return new Promise((resolve, reject) => {
      jwt.verify(
        refreshToken,
        process.env.JWT_SECRET,
        async (err, payload: any) => {
          if (err) reject(err);
          const userId = payload.id;
          const cachedToken = await this.cacheManager.get(userId);
          if (refreshToken != cachedToken) {
            return reject(new Error('Invalid Token Credentials'));
          }
          const returnPayload: tokenPayload = {
            id: payload.id,
            username: payload.username,
          };
          resolve(returnPayload);
        },
      );
    });
  }
}
