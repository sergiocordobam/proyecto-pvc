import { Test, TestingModule } from '@nestjs/testing';
import { OperatorController } from '../src/operator/controllers/operator.controller';
import { OperatorFetchService } from '../src/operator/services/operator-fetch.service';
import { OperatorRegistrationService } from '../src/operator/services/operator-registration.service';
import { TokenService } from '../src/operator/services/token.service';

describe('OperatorController', () => {
  let controller: OperatorController;

  const mockFetchService = {
    getOperators: jest.fn(),
    getSelfOperator: jest.fn(),
    getOperatorByName: jest.fn(),
  };

  const mockRegistrationService = {
    registerOperator: jest.fn(),
    registerEndPoint: jest.fn(),
  };

  const mockTokenService = {
    saveToken: jest.fn(),
  };

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      controllers: [OperatorController],
      providers: [
        { provide: OperatorFetchService, useValue: mockFetchService },
        { provide: OperatorRegistrationService, useValue: mockRegistrationService },
        { provide: TokenService, useValue: mockTokenService },
      ],
    }).compile();

    controller = module.get<OperatorController>(OperatorController);
  });

  it('should fetch all operators and respond with 200', async () => {
    const req = {};
    const res = {
      status: jest.fn().mockReturnThis(),
      json: jest.fn(),
    };
    const data = [{ name: 'Operator 1' }];
    mockFetchService.getOperators.mockResolvedValue(data);

    await controller.fetchOperators(req as any, res as any);

    expect(mockFetchService.getOperators).toHaveBeenCalled();
    expect(res.status).toHaveBeenCalledWith(200);
    expect(res.json).toHaveBeenCalledWith(data);
  });

  it('should fetch self operator and respond with 200', async () => {
    const req = {};
    const res = {
      status: jest.fn().mockReturnThis(),
      json: jest.fn(),
    };
    const data = { name: 'Self Operator' };
    mockFetchService.getSelfOperator.mockResolvedValue(data);

    await controller.fetchSelfOperator(req as any, res as any);

    expect(mockFetchService.getSelfOperator).toHaveBeenCalled();
    expect(res.status).toHaveBeenCalledWith(200);
    expect(res.json).toHaveBeenCalledWith(data);
  });
});
