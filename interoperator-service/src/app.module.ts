/*import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { OperatorModule } from './operator/operator.module';

@Module({
    imports: [ConfigModule.forRoot({ isGlobal: true }), OperatorModule],
})
export class AppModule {}**/
import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { OperatorModule } from './operator/operator.module';
import { TransferModule } from './operator/transfer-module';


@Module({
    imports: [
        ConfigModule.forRoot({ isGlobal: true }),
        OperatorModule,
        TransferModule,
    ],
})
export class AppModule {}