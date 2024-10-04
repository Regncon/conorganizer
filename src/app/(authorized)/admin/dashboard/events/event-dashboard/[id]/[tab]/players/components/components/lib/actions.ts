'use server';
import { FirebaseCollectionNames, InterestLevel, RoomName } from '$lib/enums';
import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import { Interest, Participant, PoolEvent, PoolPlayer } from '$lib/types';
import { addDoc, collection, doc, getDoc, getDocs, setDoc, updateDoc } from 'firebase/firestore';

export async function assignPlayer(
    participantId: string,
    poolEventId: string,
    isAssigned: boolean,
    isGameMaster: boolean,
    poolPlayerId?: string
) {
    console.log('assignPlayer', participantId, poolEventId, isAssigned, isGameMaster);
    if (poolPlayerId) {
        await updatePoolPlayer(poolPlayerId, poolEventId, isAssigned, isGameMaster);
        return;
    }
    await createPoolPlayer(participantId, poolEventId, isAssigned, isGameMaster);
    return;
}

async function updatePoolPlayer(poolPlayerId: string, poolEventId: string, isAssigned: boolean, isGameMaster: boolean) {
    if (isGameMaster) {
        isAssigned = true;
    }

    const { db, user } = await getAuthorizedAuth();
    if (!db || !user) {
        throw new Error('Database is undefined');
    }
    // Todo: add updated by and updated at
    // Add or delete pool and participant
    const poolPlayerRef = doc(db, FirebaseCollectionNames.poolPlayers, poolPlayerId);
    await updateDoc(poolPlayerRef, { isAssigned, isGameMaster });
}

async function createPoolPlayer(
    participantId: string,
    poolEventId: string,
    isAssigned: boolean,
    isGameMaster: boolean
) {
    const { db, user } = await getAuthorizedAuth();

    if (!db || !user) {
        throw new Error('Database is undefined');
    }
    console.log('createPoolPlayer', participantId, poolEventId, isAssigned, isGameMaster);

    const poolEventRef = doc(db, FirebaseCollectionNames.poolEvents, poolEventId);
    const poolEvent = (await getDoc(poolEventRef)).data() as PoolEvent;

    if (!poolEvent) {
        throw new Error('Pool event does not exist');
    }
    const roomsRef = collection(db, FirebaseCollectionNames.poolEvents, poolEventId, FirebaseCollectionNames.rooms);
    const rooms = await getDocs(roomsRef);

    const participantRef = doc(db, FirebaseCollectionNames.participants, participantId);
    const participant = (await getDoc(participantRef)).data() as Participant;

    const interestRef = doc(
        db,
        FirebaseCollectionNames.poolEvents,
        poolEventId,
        FirebaseCollectionNames.interests,
        participantId
    );
    const interest = (await getDoc(interestRef)).data() as Interest;

    const isFistChoice = isGameMaster ? false : interest?.interestLevel === InterestLevel.VeryInterested;
    const isPublished = isGameMaster;

    const newPlayer: PoolPlayer = {
        participantId: participantId,
        poolEventId: poolEventId,
        isGameMaster: isGameMaster,
        firstName: participant.firstName,
        lastName: participant.lastName,
        interestLevel: interest?.interestLevel ?? InterestLevel.NotInterested,
        poolEventTitle: poolEvent.title,
        poolName: poolEvent.poolName,
        roomId: rooms.docs[0]?.id ?? '',
        roomName: rooms.docs[0]?.data().name ?? RoomName.NotSet,
        isPublished: isPublished,
        isAssigned: isAssigned,
        isFirstChoice: isFistChoice,
        createdAt: Date.now().toString(),
        createdBy: user.uid,
        updateAt: Date.now().toString(),
        updatedBy: user.uid,
    };
    console.log('newPlayer', newPlayer);

    const poolPlayerRef = collection(db, FirebaseCollectionNames.players);
    const addPlayesResponce = await addDoc(poolPlayerRef, newPlayer);
    const newPoolPlayerId = addPlayesResponce.id;
    console.log('poolPlayerId', newPoolPlayerId);

    const poolEventPoolPlayerRef = doc(
        db,
        FirebaseCollectionNames.poolEvents,
        poolEventId,
        FirebaseCollectionNames.poolPlayers,
        newPoolPlayerId
    );
    await setDoc(poolEventPoolPlayerRef, newPlayer);

    const participantPoolPlayerRef = doc(
        db,
        FirebaseCollectionNames.participants,
        participantId,
        FirebaseCollectionNames.particitantPlayers,
        newPoolPlayerId
    );
    await setDoc(participantPoolPlayerRef, newPlayer);
}
