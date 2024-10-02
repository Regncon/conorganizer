import { adminDb } from '$lib/firebase/firebaseAdmin';
import { ConEvent, Interest, PlayerInterest } from '$lib/types';

export async function generatePoolPlayerInterestById(id: string) {
    let poolInterests: PlayerInterest[] = [];

    const interestRef = await adminDb.collection('pool-events').doc(id).collection('interests').get();
    const interests = interestRef.docs.map((doc) => ({ id: doc.id, ...doc.data() })) as Interest[];

    interests.forEach((interest) => {
        const playerInterest: PlayerInterest = {
            poolEventId: interest.poolEventId,
            interestLevel: interest.interestLevel,
            participantId: interest.participantId,
            firstName: interest.participantFirstName,
            lastName: interest.participantLastName,
            isOver18: false,
            ticketCategoryID: 0,
            ticketCategory: '',
            conPlayers: [],
        };

        poolInterests.push(playerInterest);
    });

    return poolInterests;
}
