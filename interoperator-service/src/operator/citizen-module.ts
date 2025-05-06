import { Module } from '@nestjs/common';
import { CitizenService } from './services/Auth-service-Conection';

@Module({
    providers: [CitizenService],
    exports: [CitizenService],
})
export class CitizenModule {}