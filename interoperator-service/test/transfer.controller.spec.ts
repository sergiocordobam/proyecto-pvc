import { Test, TestingModule } from '@nestjs/testing';
import { TransferController } from '../src/operator/controllers/transfer.controler';
import { TransferService } from '../src/operator/services/transfer.service';
import { ConfirmTransferDto } from '../src/operator/DTO/ConfirmTransferDto';
import { TransferRequestDto } from '../src/operator/DTO/TransferRequestDto';
import { RegisterCitizenDto } from '../src/operator/DTO/RegisterCitizenDTO';
describe('TransferController', () => {
  let controller: TransferController;
  let service: TransferService;

  const mockTransferService = {
    confirmTransfer: jest.fn(),
    processTransfer: jest.fn(),
    registerCitizenAndDocuments: jest.fn(),
  };

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      controllers: [TransferController],
      providers: [
        {
          provide: TransferService,
          useValue: mockTransferService,
        },
      ],
    }).compile();

    controller = module.get<TransferController>(TransferController);
    service = module.get<TransferService>(TransferService);
  });

  it('should call confirmTransfer on the service', async () => {
    const dto = { operatorId: 'op-1', transferId: '123', status: 'confirmed' };
    await controller.confirmCitizenTransfer(dto);
    expect(service.confirmTransfer).toHaveBeenCalledWith(dto);
  });

  it('should call processTransfer on the service', async () => {
    const dto = { userId: '123', sourceOperator: 'op-1', targetOperator: 'op-2' };
    await controller.processTransfer(dto);
    expect(service.processTransfer).toHaveBeenCalledWith(dto);
  });

  it('should call registerCitizenAndDocuments and return success message', async () => {
    const payload = {
      id: '123',
      citizenName: 'John Doe',
      citizenEmail: 'john@example.com',
      urlDocuments: {
        national: ['doc1'],
        extra: ['doc2'],
      },
      confirmAPI: 'http://confirm-api.com',
    };

    const result = await controller.registerCitizen(payload);
    expect(service.registerCitizenAndDocuments).toHaveBeenCalledWith(payload);
    expect(result).toEqual({ message: 'Citizen and documents registered successfully' });
  });
});
