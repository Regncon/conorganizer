'use server';
import { FirebaseCollectionNames, RoomName } from '$lib/enums';
import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import { Interest, PoolEvent, PoolPlayer } from '$lib/types';
import { collection, doc, getDoc, getDocs, setDoc, updateDoc } from 'firebase/firestore';

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
    const poolPlayerRef = doc(db, FirebaseCollectionNames.poolPlayers, poolPlayerId);
    await updateDoc(poolPlayerRef, { isAssigned, isGameMaster });
}

async function createPoolPlayer(
    participantId: string,
    poolEventId: string,
    isAssigned: boolean,
    isGameMaster: boolean
) {
    const { db } = await getAuthorizedAuth();

    if (!db) {
        throw new Error('Database is undefined');
    }

    const poolEventRef = doc(db, FirebaseCollectionNames.poolEvents, poolEventId);
    const poolEvent = (await getDoc(poolEventRef)).data() as PoolEvent;

    if (!poolEvent) {
        throw new Error('Pool event does not exist');
    }

    const interestRef = doc(db, FirebaseCollectionNames.poolEvents, FirebaseCollectionNames.interests, participantId);
    const interest = (await getDoc(interestRef)).data() as Interest;

    const newPlayer: PoolPlayer = {
        participantId: participantId,
        poolEventId: poolEventId,
        isGameMaster: isGameMaster,
        firstName: interest.participantFirstName,
        lastName: interest.participantLastName,
        interestLevel: interest.interestLevel,
        poolEventTitle: poolEvent.title,
        poolName: poolEvent.poolName,
        roomId: '',
        roomName: RoomName.NotSet,
        isPublished: false,
        isFirstChoice: false,
        createdAt: '',
        createdBy: '',
        updateAt: '',
        updatedBy: '',
    };
}

// get participant
// get pool event
// create pool player
// add pool player to pool event
// add pool player to participant
