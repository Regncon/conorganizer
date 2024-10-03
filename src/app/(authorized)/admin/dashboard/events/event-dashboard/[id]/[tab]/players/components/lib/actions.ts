import { adminDb, getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import { ConEvent, Interest, Participant, PlayerInterest } from '$lib/types';

export async function generatePoolPlayerInterestById(id: string) {
    const { db, user } = await getAuthorizedAuth();

    let poolInterests: PlayerInterest[] = [];

    const interestRef = await adminDb.collection('pool-events').doc(id).collection('interests').get();
    const interests = interestRef.docs.map((doc) => ({ id: doc.id, ...doc.data() })) as Interest[];

    const participants = (await adminDb.collection('participants').get()).docs.map((doc) => ({
        id: doc.id,
        ...doc.data(),
    })) as Participant[];

    interests.forEach((interest) => {
        const participant = participants.find((participant) => participant.id === interest.participantId);
        const playerInterest: PlayerInterest = {
            poolEventId: interest.poolEventId,
            interestLevel: interest.interestLevel,
            participantId: interest.participantId,
            firstName: interest.participantFirstName,
            lastName: interest.participantLastName,
            isOver18: participant?.over18 ?? false,
            ticketCategoryID: participant?.ticketCategoryId ?? 0,
            ticketCategory: participant?.ticketCategory ?? 'FEIL!',
            conPlayers: [],
            isGameMaster: true,
            isAssigned: true,
        };

        poolInterests.push(playerInterest);
    });

    return poolInterests;
}
