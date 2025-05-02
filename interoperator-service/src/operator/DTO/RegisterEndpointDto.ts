import { IsNotEmpty, IsString } from 'class-validator';

export class RegisterEndpointDto {
    @IsNotEmpty()
    @IsString()
    idOperator: string;

    @IsNotEmpty()
    @IsString()
    endPoint: string;

    @IsNotEmpty()
    @IsString()
    endPointConfirm: string;
}