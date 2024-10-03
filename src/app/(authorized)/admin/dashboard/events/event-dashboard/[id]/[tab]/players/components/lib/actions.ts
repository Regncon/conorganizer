import { FirebaseCollectionNames } from '$lib/enums';
import { adminDb, getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import { ConEvent, Interest, Participant, PlayerInterest, PoolPlayer } from '$lib/types';
import { collection, getDocs } from 'firebase/firestore';

export async function generatePoolPlayerInterestById(id: string) {
    const { db, user } = await getAuthorizedAuth();

    let poolInterests: PlayerInterest[] = [];

    const interestRef = await adminDb.collection('pool-events').doc(id).collection('interests').get();
    const interests = interestRef.docs.map((doc) => ({ id: doc.id, ...doc.data() })) as Interest[];

    const participants = (await adminDb.collection('participants').get()).docs.map((doc) => ({
        id: doc.id,
        ...doc.data(),
    })) as Participant[];

    const poolPlayers = await getPoolPlayers();

    interests.forEach((interest) => {
        const participant = participants.find((participant) => participant.id === interest.participantId);
        const participantPoolPlayers = poolPlayers.filter(
            (poopPlayer) => poopPlayer.participantId === interest.participantId
        );
        const playerInterest: PlayerInterest = {
            poolEventId: interest.poolEventId,
            interestLevel: interest.interestLevel,
            participantId: interest.participantId,
            firstName: interest.participantFirstName,
            lastName: interest.participantLastName,
            isOver18: participant?.over18 ?? false,
            ticketCategoryID: participant?.ticketCategoryId ?? 0,
            ticketCategory: participant?.ticketCategory ?? 'FEIL!',
            poolPlayers: participantPoolPlayers ? participantPoolPlayers : [],
            isGameMaster: true,
            isAssigned: true,
        };

        poolInterests.push(playerInterest);
    });

    return poolInterests;
}
export async function getPoolPlayers() {
    const { db } = await getAuthorizedAuth();
    if (!db) {
        throw new Error('Database is undefined');
    }
    const poolPlayersRef = collection(db, FirebaseCollectionNames.poolPlayers);
    const poolPlayersSnapshot = await getDocs(poolPlayersRef);
    const poolPlayers = poolPlayersSnapshot.docs.map((doc) => ({ id: doc.id, ...doc.data() })) as PoolPlayer[];
    console.log('poolPlayers', poolPlayers);

    return poolPlayers;
}
