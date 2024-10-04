'use server';
import { FirebaseCollectionNames, InterestLevel, RoomName } from '$lib/enums';
import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import { Interest, Participant, PoolEvent, PoolPlayer } from '$lib/types';
import { addDoc, collection, deleteDoc, doc, getDoc, getDocs, setDoc, updateDoc } from 'firebase/firestore';

export async function assignPlayer(
    participantId: string,
    poolEventId: string,
    isAssigned: boolean,
    isGameMaster: boolean,
    poolPlayerId?: string
) {
    console.log('assignPlayer', participantId, poolEventId, isAssigned, isGameMaster);
    if (poolPlayerId) {
        await updatePoolPlayer(poolPlayerId, poolEventId, participantId, isAssigned, isGameMaster);
        return;
    }
    await createPoolPlayer(participantId, poolEventId, isAssigned, isGameMaster);
    return;
}

async function updatePoolPlayer(
    poolPlayerId: string,
    poolEventId: string,
    participantId: string,
    isAssigned: boolean,
    isGameMaster: boolean
) {
    if (isGameMaster) {
        isAssigned = true;
    }
    console.log('updatePoolPlayer', poolPlayerId, poolEventId, isAssigned, isGameMaster);

    const { db, user } = await getAuthorizedAuth();
    if (!db || !user) {
        throw new Error('Database is undefined');
    }

    const playerRef = doc(db, FirebaseCollectionNames.players, poolPlayerId);

    const poolEventPoolPlayerRef = doc(
        db,
        FirebaseCollectionNames.poolEvents,
        poolEventId,
        FirebaseCollectionNames.poolPlayers,
        poolPlayerId
    );
    const participantPoolPlayerRef = doc(
        db,
        FirebaseCollectionNames.participants,
        participantId,
        FirebaseCollectionNames.particitantPlayers,
        poolPlayerId
    );

    if (isAssigned || isGameMaster) {
        console.log('updating player', poolPlayerId);

        await updateDoc(playerRef, {
            isAssigned,
            isGameMaster,
            updateAt: new Date().toISOString(),
            updatedBy: user.uid,
        });

        await updateDoc(poolEventPoolPlayerRef, {
            isAssigned,
            isGameMaster,
            updateAt: new Date().toISOString(),
            updatedBy: user.uid,
        });

        await updateDoc(participantPoolPlayerRef, {
            isAssigned,
            isGameMaster,
            updateAt: new Date().toISOString(),
            updatedBy: user.uid,
        });
    }
    if (isAssigned === false && isGameMaster === false) {
        console.log('deleteting player', poolPlayerId);

        await deleteDoc(playerRef);
        await deleteDoc(poolEventPoolPlayerRef);
        await deleteDoc(participantPoolPlayerRef);
    }
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
        isAssigned: isAssigned || isGameMaster,
        isFirstChoice: isFistChoice,
        createdAt: new Date().toISOString(),
        createdBy: user.uid,
        updateAt: new Date().toISOString(),
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
