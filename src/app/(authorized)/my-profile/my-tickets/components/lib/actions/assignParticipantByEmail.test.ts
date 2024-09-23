// assignParticipantByEmail.test.ts

import { getMyUserInfo } from "$app/(authorized)/my-events/lib/actions";
import { adminDb, getAuthorizedAuth } from "$lib/firebase/firebaseAdmin";
import { AssignParticipantByEmail, GetParticipantsByEmail, GetTicketsByEmail } from "./actions";

jest.mock('./auth');
jest.mock('./tickets');
jest.mock('./participantsData');
jest.mock('./user');
jest.mock('./database');

describe('AssignParticipantByEmail', () => {
    beforeEach(() => {
        jest.clearAllMocks();
    });

    it('should assign participants by email', async () => {
        // Mock implementations
        (getAuthorizedAuth as jest.Mock).mockResolvedValue({
            db: {},
            user: { email: 'test@example.com', uid: 'user123' },
        });

        (GetTicketsByEmail as jest.Mock).mockResolvedValue([
            /* Mock ticket data */
        ]);

        (GetParticipantsByEmail as jest.Mock).mockResolvedValue([
            /* Mock participant data */
        ]);

        (getMyUserInfo as jest.Mock).mockResolvedValue(null);

        // Mock database methods
        adminDb.collection = jest.fn().mockReturnThis();
        adminDb.doc = jest.fn().mockReturnThis();
        adminDb.add = jest.fn().mockResolvedValue({ id: 'newParticipantId' });
        adminDb.update = jest.fn().mockResolvedValue({});
        adminDb.set = jest.fn().mockResolvedValue({});

        // Call the function
        const result = await AssignParticipantByEmail();

        // Assertions
        expect(getAuthorizedAuth).toHaveBeenCalled();
        expect(GetTicketsByEmail).toHaveBeenCalled();
        expect(GetParticipantsByEmail).toHaveBeenCalled();
        expect(getMyUserInfo).toHaveBeenCalled();
        expect(adminDb.collection).toHaveBeenCalledWith('participants');
        expect(adminDb.add).toHaveBeenCalled();
        expect(adminDb.collection).toHaveBeenCalledWith('users');
        expect(adminDb.set).toHaveBeenCalled();
        expect(result).toBeDefined();
    });
});
