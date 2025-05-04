import { IsNotEmpty, IsString } from 'class-validator';

export class ConfirmTransferDto {
    @IsNotEmpty()
    @IsString()
    operatorId: string; // ID of the target operator

    @IsNotEmpty()
    @IsString()
    transferId: string; // ID of the transfer to confirm

    @IsNotEmpty()
    @IsString()
    status: string; // Status of the transfer (e.g., "confirmed", "rejected")
}