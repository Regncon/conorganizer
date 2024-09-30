'use server';

import type { MyUserInfo } from '$app/(authorized)/my-events/lib/types';
import { adminDb } from '$lib/firebase/firebaseAdmin';
import type { Interest, InterestsInPool, Participant, PoolEvent } from '$lib/types';

export const migrateParticipantAndInterest = async () => {
    // const allUserData = (await adminDb.collection('users').get()).docs.map((doc) => ({
    //     id: doc.id,
    //     ...doc.data(),
    // })) as (MyUserInfo & { id: string })[];
    // const allPoolEvents = (await adminDb.collection('pool-events').get()).docs.map((doc) => ({
    //     id: doc.id,
    //     ...doc.data(),
    // })) as PoolEvent[];

    const GetAllParticipants = (await adminDb.collectionGroup('interests').get()).docs.map((doc) => ({
        id: doc.id,
        ...doc.data(),
    })) as Interest[];
    const GetAllParticipants2 = (await adminDb.collectionGroup('participant-interests').get()).docs.map((doc) => ({
        id: doc.id,
        ...doc.data(),
    })) as Interest[];

    //missing poolName
    console.log(
        'GetAllParticipants',
        GetAllParticipants.filter((participant) => participant.poolName === undefined)
    );
    console.log(
        'GetAllParticipants2',
        GetAllParticipants2.filter((participant) => participant.poolName === undefined)
    );

    // missing firstname
    console.log(
        'GetAllParticipants',
        GetAllParticipants.filter((participant) => participant.participantFirstName === undefined)
    );
    console.log(
        'GetAllParticipants2',
        GetAllParticipants2.filter((participant) => participant.participantFirstName === undefined)
    );

    // migrate inn poolnames to intrests (from 'interest' and 'participant-interests')

    // migrate in name from pardicipantFirstName to participantFirstName  (from 'interest' and 'participant-interests')
};
