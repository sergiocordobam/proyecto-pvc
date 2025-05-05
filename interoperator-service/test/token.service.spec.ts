import { Test, TestingModule } from '@nestjs/testing';
import { TokenService } from '../src/operator/services/token.service';

describe('TokenService', () => {
    let tokenService: TokenService;

    beforeEach(async () => {
        const module: TestingModule = await Test.createTestingModule({
            providers: [TokenService],
        }).compile();

        tokenService = module.get<TokenService>(TokenService);
    });

    it('should save a token', async () => {
        const token = 'test-token';

        await tokenService.saveToken(token);

        expect(process.env.OPERATOR_ID).toEqual(token);
    });
});