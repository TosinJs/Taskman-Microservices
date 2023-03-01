import { InjectModel } from '@nestjs/mongoose';
import { User, UserDocument } from './schemas/user.schema';
import { Model } from 'mongoose';

export class DBService {
  constructor(@InjectModel(User.name) private userModel: Model<UserDocument>) {}

  async createUser(user: Partial<User>) {
    const newUser = new this.userModel();
    newUser.username = user.username;
    newUser.password = user.password;
    newUser.email = user.email;
    await newUser.save();
    return newUser;
  }

  async findByUsername(username: string) {
    const user = await this.userModel.findOne({ username });
    return user;
  }

  async findByIdAndUpdate(id: string, update: Partial<User>) {
    await this.userModel.findByIdAndUpdate(id, update);
  }
}
