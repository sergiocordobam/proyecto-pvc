import { Module } from '@nestjs/common';
import { DocumentService } from './services/Documents-service-conection';

@Module({
    providers: [DocumentService],
    exports: [DocumentService],
})
export class DocumentModule {}