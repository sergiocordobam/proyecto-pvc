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
import { CitizenModule } from './operator/citizen-module';
import { DocumentModule } from './operator/document-module';

@Module({
    imports: [
        ConfigModule.forRoot({ isGlobal: true }),
        OperatorModule,
        TransferModule,
        CitizenModule,
        DocumentModule,
    ],
})
export class AppModule {}