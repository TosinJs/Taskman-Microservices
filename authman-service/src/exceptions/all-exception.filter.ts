import {
  ArgumentsHost,
  BadRequestException,
  Catch,
  ExceptionFilter,
  HttpException,
  HttpStatus,
} from '@nestjs/common';
import { BaseExceptionFilter } from '@nestjs/core';
import { Request, Response } from 'express';
import * as dotenv from 'dotenv';
dotenv.config();

@Catch()
export class AllExceptionsFilter implements ExceptionFilter {
  catch(exception: Error, host: ArgumentsHost) {
    const ctx = host.switchToHttp();
    const response = ctx.getResponse<Response>();
    const request = ctx.getRequest<Request>();

    let status =
      exception instanceof HttpException
        ? exception.getStatus()
        : HttpStatus.INTERNAL_SERVER_ERROR;
    let message: string | object =
      exception instanceof HttpException
        ? exception.message
        : 'Internal Server Error';

    if (exception instanceof BadRequestException) {
      message = exception.getResponse();
    }

    const devResponse = {
      statusCode: status,
      timestamp: new Date().toISOString(),
      path: request.url,
      trace: exception?.stack?.split('/n'),
      error: exception,
    };

    //log devResponse

    const prodResponse = {
      statusCode: status,
      message: message,
      path: request.url,
    };
    response
      .status(status)
      .json(process.env.ENV == 'DEV' ? devResponse : prodResponse);
  }
}

@Catch()
export class AllExceptionsFilterExtendsBase extends BaseExceptionFilter {
  catch(exception: unknown, host: ArgumentsHost): void {
    super.catch(exception, host);
  }
}
