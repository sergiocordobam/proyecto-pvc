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
import { Controller, Post, Body, HttpCode, HttpStatus } from '@nestjs/common';
import { TransferService } from '../services/transfer.service';
import { TransferCitizenDto } from '../DTO/TransferCitizenDto';
import { ConfirmTransferDto } from '../DTO/ConfirmTransferDto';
import { TransferRequestDto } from '../DTO/TransferRequestDto';
import { RegisterCitizenDto } from '../DTO/RegisterCitizenDTO';

@Controller('transfers')
export class TransferController {
    constructor(private readonly transferService: TransferService) {}
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
    
    @Post('transfer-citizen')
    @HttpCode(HttpStatus.OK)
    async registerCitizen(@Body() payload: RegisterCitizenDto): Promise<{ message: string }> {
        await this.transferService.registerCitizenAndDocuments(payload);
        return { message: 'Citizen and documents registered successfully' };
    }
}