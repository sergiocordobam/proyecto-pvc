import { Test, TestingModule } from '@nestjs/testing';
import { TransferService } from '../src/operator/services/transfer.service';
import { ClientProxyFactory } from '@nestjs/microservices';
import { of } from 'rxjs';

jest.mock('@nestjs/microservices', () => ({
    ...jest.requireActual('@nestjs/microservices'),
    ClientProxyFactory: {
        create: jest.fn(() => ({
            send: jest.fn(),
            emit: jest.fn(),
        })),
    },
}));

describe('TransferService', () => {
    let transferService: TransferService;

    beforeEach(async () => {
        const module: TestingModule = await Test.createTestingModule({
            providers: [TransferService],
        }).compile();

        transferService = module.get<TransferService>(TransferService);

        // Mock the send and emit methods for all ClientProxy instances
        jest.spyOn(transferService['fetchCitizenInfoClient'], 'send').mockReturnValue(
            of({ name: 'John Doe', email: 'john.doe@example.com' }),
        );
        jest.spyOn(transferService['fetchDocumentUrlsClient'], 'send').mockReturnValue(
            of(['http://example.com/doc1', 'http://example.com/doc2']),
        );
        jest.spyOn(transferService['transferConfirmationsClient'], 'emit').mockReturnValue(
            of(null),
        );
    });

    it('should process transfer', async () => {
        const mockMessage = { citizenId: '123', operatorId: 'op1' };

        await transferService.processTransfer(mockMessage);

        expect(transferService['fetchCitizenInfoClient'].send).toHaveBeenCalledWith(
            'fetch_citizen_info',
            { citizenId: '123' },
        );
        expect(transferService['fetchDocumentUrlsClient'].send).toHaveBeenCalledWith(
            'fetch_document_urls',
            { citizenId: '123' },
        );
        expect(transferService['transferConfirmationsClient'].emit).toHaveBeenCalledWith(
            'transfer_confirmations',
            { transferId: '123', status: 'confirmed' },
        );
    });
});