import { IsString, IsEmail, IsArray, IsObject } from 'class-validator';

export class RegisterCitizenDto {
    @IsString()
    id: string;

    @IsString()
    citizenName: string;

    @IsEmail()
    citizenEmail: string;

    @IsObject()
    urlDocuments: Record<string, string[]>;
}