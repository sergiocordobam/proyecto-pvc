/*import { Controller, Post, Body, HttpCode, HttpStatus } from '@nestjs/common';
import { TransferService } from '../services/transfer.service';
import { TransferCitizenDto } from '../DTO/TransferCitizenDto';
import { ConfirmTransferDto } from '../DTO/ConfirmTransferDto';
import { TransferRequestDto } from '../DTO/TransferRequestDto';

@Controller('transfers')
export class TransferController {
    constructor(private readonly transferService: TransferService) {}

    @Post('transfer-citizen')
    @HttpCode(HttpStatus.OK)
    async transferCitizen(@Body() dto: TransferCitizenDto): Promise<any> {
        return this.transferService.transferCitizen(dto);
    }

    @Post('confirm-citizen-transfer')
    @HttpCode(HttpStatus.OK)
    async confirmCitizenTransfer(@Body() dto: ConfirmTransferDto): Promise<any> {
        return this.transferService.confirmTransfer(dto);
    }

    @Post('process-transfer')
    @HttpCode(HttpStatus.OK)
    async processTransfer(@Body() dto: TransferRequestDto): Promise<any> {
        return this.transferService.processTransfer(dto);
    }
}*/
import { Injectable } from '@nestjs/common';
import { EventPattern, Payload } from '@nestjs/microservices';

@Injectable()
export class TransferController {
    @EventPattern('fetch_citizen_info')
    async fetchCitizenInfo(@Payload() data: { citizenId: string }): Promise<any> {
        console.log(`Fetching citizen info for ID: ${data.citizenId}`);
        // Simulate fetching citizen info
        return { id: data.citizenId, name: 'John Doe', email: 'john.doe@example.com' };
    }

    @EventPattern('delete_citizen')
    async deleteCitizen(@Payload() data: { citizenId: string }): Promise<void> {
        console.log(`Deleting citizen with ID: ${data.citizenId}`);
        // Simulate deletion logic
    }
}