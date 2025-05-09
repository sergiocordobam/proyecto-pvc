import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import * as dotenv from 'dotenv';

dotenv.config({ path: './src/Config/dev.env' });
//docker build -t interoperator-service .
//docker run -p 8080:3000 --env-file ./src/config/dev.env interoperator-service

//http://localhost:3000/comunication/operators
console.log('API_BASE_URL:', process.env.API_BASE_URL);
async function bootstrap() {
  const app = await NestFactory.create(AppModule);
    app.enableCors({
    origin: '*', // Cambia esto a los orígenes permitidos
    methods: ['GET', 'POST', 'PUT', 'DELETE'],
    allowedHeaders: ['Content-Type', 'Authorization'],
    credentials: true, // Habilita el uso de cookies/autenticación
  });
  app.setGlobalPrefix('comunication'); // Set global prefix for routes
  await app.listen(3000);
  console.log(`Server is running on http://localhost:3000`);
}
bootstrap();