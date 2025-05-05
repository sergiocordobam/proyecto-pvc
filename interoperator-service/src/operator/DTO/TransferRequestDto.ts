import { IsNotEmpty, IsString } from 'class-validator';

export class TransferRequestDto {
    @IsNotEmpty()
    @IsString()
    userId: string; // ID of the user being transferred

    @IsNotEmpty()
    @IsString()
    sourceOperator: string; // ID of the source operator

    @IsNotEmpty()
    @IsString()
    targetOperator: string; // ID of the target operator
}