import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import * as dotenv from 'dotenv';

dotenv.config({ path: './src/Config/dev.env' });

async function bootstrap() {
  const app = await NestFactory.create(AppModule);
  app.setGlobalPrefix('comunication'); // Set global prefix for routes
  await app.listen(3000);
  console.log(`Server is running on http://localhost:3000`);
}
bootstrap();