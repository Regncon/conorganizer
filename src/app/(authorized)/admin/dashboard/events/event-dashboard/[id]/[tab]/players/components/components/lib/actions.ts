'use server';
import { FirebaseCollectionNames, InterestLevel, RoomName } from '$lib/enums';
import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import { Interest, Participant, PoolEvent, PoolPlayer } from '$lib/types';
import { collection, deleteDoc, doc, getDoc, getDocs, setDoc, updateDoc } from 'firebase/firestore';

export async function assignPlayer(
    poolPlayerId: string,
    participantId: string,
    poolEventId: string,
    isAssigned: boolean,
    isGameMaster: boolean
) {
    console.log('assignPlayer', participantId, poolEventId, isAssigned, isGameMaster);
    if (poolPlayerId) {
        await updatePoolPlayer(poolPlayerId, isAssigned, isGameMaster);
        return;
    }
    await createPoolPlayer(participantId, poolEventId, isAssigned, isGameMaster);
    return;
}

async function updatePoolPlayer(poolPlayerId: string, isAssigned: boolean, isGameMaster: boolean) {
    if (isGameMaster) {
        isAssigned = true;
    }

    const { db } = await getAuthorizedAuth();
    if (!db) {
        throw new Error('Database is undefined');
    }
    // Todo: add updated by and updated at
    // Add or delete pool and participant
    const poolPlayerRef = doc(db, FirebaseCollectionNames.poolPlayers, poolPlayerId);
    if (isAssigned === false) {
        await deleteDoc(poolPlayerRef);
        return;
    }

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

    const poolEventRef = doc(db, FirebaseCollectionNames.poolEvents, poolEventId);
    const poolEvent = (await getDoc(poolEventRef)).data() as PoolEvent;

    if (!poolEvent) {
        throw new Error('Pool event does not exist');
    }
    const roomsRef = collection(db, FirebaseCollectionNames.poolEvents, poolEventId, FirebaseCollectionNames.rooms);
    const rooms = await getDocs(roomsRef);
    const room = rooms.docs[0].data();

    const participantRef = doc(db, FirebaseCollectionNames.participants, participantId);
    const participant = (await getDoc(participantRef)).data() as Participant;

    const interestRef = doc(db, FirebaseCollectionNames.poolEvents, FirebaseCollectionNames.interests, participantId);
    const interest = (await getDoc(interestRef)).data() as Interest;

    const newPlayer: PoolPlayer = {
        participantId: participantId,
        poolEventId: poolEventId,
        isGameMaster: isGameMaster,
        firstName: participant.firstName,
        lastName: participant.lastName,
        interestLevel: interest?.interestLevel ?? InterestLevel.NotInterested,
        poolEventTitle: poolEvent.title,
        poolName: poolEvent.poolName,
        roomId: room.id,
        roomName: room.name,
        isPublished: false,
        isFirstChoice: interest?.interestLevel === InterestLevel.VeryInterested,
        createdAt: Date.now().toString(),
        createdBy: user.uid,
        updateAt: Date.now().toString(),
        updatedBy: user.uid,
    };
}

// get participant
// get pool event
// create pool player
// add pool player to pool event
// add pool player to participant
