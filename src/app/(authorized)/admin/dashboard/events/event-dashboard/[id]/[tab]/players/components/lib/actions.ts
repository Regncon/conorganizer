import { FirebaseCollectionNames, RoomName } from '$lib/enums';
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
    const filterdInterests = interests.filter((interest) => interest.interestLevel > 0);
    const playerInterestsAsPoolPlayers = filterdInterests.map((interest) => {
        const test = { 
            interestLevel: interest.interestLevel,
            poolEventId: interest.poolEventId,
            participantId: interest.participantId,
            firstName: interest.participantFirstName,
            lastName: interest.participantLastName,
            poolName: interest.poolName,
            roomId: '',
            roomName: RoomName.NotSet,
            isPublished: false,
            isFirstChoice: false,
            isGameMaster: false,
            isAssigned: false,
            createdAt: '',
            createdBy: '',
            updateAt: '',
            updatedBy: '',
            poolEventTitle: '',
            id: undefined

        } as PoolPlayer;
        return test;
        
    });

    const filterdPoolPlayers = poolPlayers.filter((poolPlayer) => poolPlayer.poolEventId === id);


    let poolPlayersWithInterests = [...filterdPoolPlayers, ...playerInterestsAsPoolPlayers];
    // remove duplicates
    poolPlayersWithInterests = poolPlayersWithInterests.filter(
        (v, i, a) => a.findIndex((t) => t.participantId === v.participantId) === i
    );




    poolPlayersWithInterests
        .forEach((interest) => {
            const participant = participants.find((participant) => participant.id === interest.participantId);
            const participantPoolPlayers = poolPlayers.filter(
                (poolPlayer) => poolPlayer.participantId === interest.participantId
            );
            const currentPoolPlayer = poolPlayers.find(
                (poolPlayer) =>
                    poolPlayer.poolEventId === interest.poolEventId &&
                    poolPlayer.participantId === interest.participantId
            );

            const isAssigned = currentPoolPlayer ? currentPoolPlayer.isAssigned : false;
            const isGameMaster = currentPoolPlayer ? currentPoolPlayer.isGameMaster : false;
            // check if the player is already assigned to another poolievent in this pool
            const isAssignedInSamePool = poolPlayers.some(
                (poolPlayer) =>
                    poolPlayer.participantId === interest.participantId &&
                    poolPlayer.poolEventId !== interest.poolEventId &&
                    poolPlayer.poolName === interest.poolName &&
                    poolPlayer.isAssigned
            );

            const playerInterest: PlayerInterest = {
                poolEventId: interest.poolEventId,
                interestLevel: interest.interestLevel,
                participantId: interest.participantId,
                currentPoolPlayerId: currentPoolPlayer?.id,
                firstName: interest.firstName,
                lastName: interest.lastName,
                isOver18: participant?.over18 ?? false,
                ticketCategoryID: participant?.ticketCategoryId ?? 0,
                ticketCategory: participant?.ticketCategory ?? 'FEIL!',
                playerInPools: participantPoolPlayers ? participantPoolPlayers : [],
                isAssigned: isAssigned,
                isGameMaster: isGameMaster,
                isAlredyPlayerInPool: isAssignedInSamePool || isAssigned || isGameMaster,
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
    const poolPlayersRef = collection(db, FirebaseCollectionNames.players);
    const poolPlayersSnapshot = await getDocs(poolPlayersRef);
    const poolPlayers = poolPlayersSnapshot.docs.map((doc) => ({ id: doc.id, ...doc.data() })) as PoolPlayer[];
    // console.log('poolPlayers', poolPlayers);

    return poolPlayers;
}
