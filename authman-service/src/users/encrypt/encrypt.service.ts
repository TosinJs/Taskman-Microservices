import { Injectable } from '@nestjs/common';
import * as crypto from 'crypto-js';

@Injectable()
export class EncryptService {
  async encrypt(text: string): Promise<string> {
    return new Promise((resolve) => {
      const encodedValue = crypto.enc.Base64.stringify(
        crypto.enc.Utf8.parse(text),
      );
      resolve(encodedValue);
    });
  }

  async decrypt(data: string): Promise<string> {
    return new Promise((resolve) => {
      const decryptedValue = crypto.enc.Base64.parse(data).toString(
        crypto.enc.Utf8,
      );
      resolve(decryptedValue);
    });
  }
}
