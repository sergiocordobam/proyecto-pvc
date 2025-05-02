import { IsNotEmpty, IsString } from 'class-validator';

export class ConfirmTransferDto {
    @IsNotEmpty()
    @IsString()
    operatorId: string;

    @IsNotEmpty()
    @IsString()
    transferId: string;

    @IsNotEmpty()
    @IsString()
    status: string;
}