import { Test, TestingModule } from '@nestjs/testing';
import { TransferController } from '../src/operator/controllers/transfer.controler';
import { TransferService } from '../src/operator/services/transfer.service';

describe('TransferController', () => {
    let transferController: TransferController;
    let transferService: TransferService;

    beforeEach(async () => {
        const module: TestingModule = await Test.createTestingModule({
            controllers: [TransferController],
            providers: [
                {
                    provide: TransferService,
                    useValue: {
                        processTransfer: jest.fn(),
                        confirmTransfer: jest.fn(),
                        registerCitizenAndDocuments: jest.fn(),
                    },
                },
            ],
        }).compile();

        transferController = module.get<TransferController>(TransferController);
        transferService = module.get<TransferService>(TransferService);
    });

    it('should process transfer', async () => {
        const mockDto = { userId: '123', sourceOperator: 'op1', targetOperator: 'op2' };
        jest.spyOn(transferService, 'processTransfer').mockResolvedValue(undefined); // Return void

        const result = await transferController.processTransfer(mockDto);

        expect(transferService.processTransfer).toHaveBeenCalledWith(mockDto);
        expect(result).toBeUndefined(); // Since the method returns void
    });
});