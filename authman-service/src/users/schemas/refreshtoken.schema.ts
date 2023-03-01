import { Prop, Schema, SchemaFactory } from '@nestjs/mongoose';
import mongoose from 'mongoose';
import { User } from './user.schema';

export type RefreshTokenDocument = RefreshToken & mongoose.Document;
@Schema()
export class RefreshToken {
  @Prop({ type: mongoose.Schema.Types.ObjectId, ref: 'User' })
  userId: User;

  @Prop({ required: true })
  roles: string;
}

export const RefreshTokenSchema = SchemaFactory.createForClass(RefreshToken);
