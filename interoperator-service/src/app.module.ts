import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { OperatorModule } from './operator/operator.module';

@Module({
    imports: [ConfigModule.forRoot({ isGlobal: true }), OperatorModule],
})
export class AppModule {}