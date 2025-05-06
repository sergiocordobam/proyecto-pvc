import { IsNotEmpty, IsString, IsArray } from 'class-validator';

export class RegisterOperatorDto {
    @IsNotEmpty()
    @IsString()
    name: string;

    @IsNotEmpty()
    @IsString()
    address: string;

    @IsNotEmpty()
    @IsString()
    contactMail: string;

    @IsArray()
    @IsString({ each: true })
    participants: string[];
}