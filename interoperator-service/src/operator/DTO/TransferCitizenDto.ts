import { IsNotEmpty, IsString } from 'class-validator';

export class TransferCitizenDto {
    @IsNotEmpty()
    @IsString()
    operatorId: string; // ID of the target operator

    @IsNotEmpty()
    @IsString()
    citizenId: string; // ID of the citizen to be transferred
}