import { Injectable } from '@nestjs/common';

@Injectable()
export class TokenService {
    async saveToken(token: string): Promise<void> {
        process.env.OPERATOR_ID = token;
        console.log(`Token saved: ${token}`);
    }
}