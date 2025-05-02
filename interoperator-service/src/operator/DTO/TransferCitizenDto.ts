import { IsNotEmpty, IsString } from 'class-validator';

export class TransferCitizenDto {
    @IsNotEmpty()
    @IsString()
    operatorId: string;

    @IsNotEmpty()
    @IsString()
    citizenId: string;

    @IsNotEmpty()
    @IsString()
    data: string;
}