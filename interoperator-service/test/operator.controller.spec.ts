import { Test, TestingModule } from '@nestjs/testing';
import { OperatorController } from '../src/operator/controllers/operator.controller';
import { OperatorFetchService } from '../src/operator/services/operator-fetch.service';
import { OperatorRegistrationService } from '../src/operator/services/operator-registration.service';
import { TokenService } from '../src/operator/services/token.service';

describe('OperatorController', () => {
    let operatorController: OperatorController;
    let fetchService: OperatorFetchService;
    let registrationService: OperatorRegistrationService;
    let tokenService: TokenService;

    beforeEach(async () => {
        const module: TestingModule = await Test.createTestingModule({
            controllers: [OperatorController],
            providers: [
                {
                    provide: OperatorFetchService,
                    useValue: {
                        getOperators: jest.fn(),
                        getOperatorByName: jest.fn(),
                        getOperatorById: jest.fn(),
                    },
                },
                {
                    provide: OperatorRegistrationService,
                    useValue: {
                        registerOperator: jest.fn(),
                        registerEndPoint: jest.fn(),
                    },
                },
                {
                    provide: TokenService,
                    useValue: {
                        saveToken: jest.fn(),
                    },
                },
            ],
        }).compile();

        operatorController = module.get<OperatorController>(OperatorController);
        fetchService = module.get<OperatorFetchService>(OperatorFetchService);
        registrationService = module.get<OperatorRegistrationService>(OperatorRegistrationService);
        tokenService = module.get<TokenService>(TokenService);
    });

    it('should fetch operators', async () => {
        const mockOperators = [{ id: '1', name: 'Operator1' }];
        jest.spyOn(fetchService, 'getOperators').mockResolvedValue(mockOperators);

        // Mock the `res` object
        const res = {
            status: jest.fn().mockReturnThis(),
            json: jest.fn(),
        };

        await operatorController.fetchOperators({} as any, res as any);

        expect(fetchService.getOperators).toHaveBeenCalled();
        expect(res.status).toHaveBeenCalledWith(200);
        expect(res.json).toHaveBeenCalledWith(mockOperators);
    });

    it('should handle errors when fetching operators', async () => {
        jest.spyOn(fetchService, 'getOperators').mockRejectedValue(new Error('Test error'));

        // Mock the `res` object
        const res = {
            status: jest.fn().mockReturnThis(),
            json: jest.fn(),
        };

        await operatorController.fetchOperators({} as any, res as any);

        expect(fetchService.getOperators).toHaveBeenCalled();
        expect(res.status).toHaveBeenCalledWith(500);
        expect(res.json).toHaveBeenCalledWith({ error: 'Failed to fetch operators' });
    });
});