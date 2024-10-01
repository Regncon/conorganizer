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

    const GetAllInterests = (await adminDb.collectionGroup('interests').get()).docs.map((doc) => {
        if (doc.ref.path.includes('participants')) {
            return {
                id: doc.id,
                ...doc.data(),
            };
        }
        return undefined;
    }) as Interest[];
    const GetAllInterests2 = (await adminDb.collectionGroup('participant-interests').get()).docs.map((doc) => {
        if (doc.ref.path.includes('participants')) {
            return { id: doc.id, ...doc.data() };
        }
        return undefined;
    }) as Interest[];

    const test = GetAllInterests.filter(Boolean).map((interest) => {
        return adminDb
            .collection('participants')
            .doc(interest.participantId)
            .collection('interests')
            .doc(interest.id ?? '')
            .get()
            .then((res) => res.ref.path);
    });

    //missing poolName
    // console.log(
    //     'GetAllParticipants',
    //     GetAllInterests.filter((participant) => participant.poolName === undefined)
    // );
    // console.log(
    //     'GetAllParticipants2',
    //     GetAllInterests2.filter((participant) => participant.poolName === undefined)
    // );

    // // missing firstname
    // console.log(
    //     'GetAllParticipants',
    //     GetAllInterests.filter((participant) => participant.participantFirstName === undefined)
    // );
    // console.log(
    //     'GetAllParticipants2',
    //     GetAllInterests2.filter((participant) => participant.participantFirstName === undefined)
    // );
    const test2 = await Promise.all(test);
    console.log(GetAllInterests.length, 'legnth');

    const onlyWrongNameProperty = GetAllInterests.filter(Boolean).filter(
        (participant) => participant.pardicipantFirstName !== undefined
    );
    const onlyWrongNameProperty2 = GetAllInterests2.filter(Boolean).filter(
        (participant) => participant.pardicipantFirstName !== undefined
    );

    // const updateThis = await adminDb
    //     .collection('participants')
    //     .doc('FqxDY3CjlGfLqW5PlfSA')
    //     .collection('interests')
    //     .doc('5IaJJcmNxarqkcnjqlcl')
    //     .update({ participantFirstName: 'test' });
    // console.log(updateThis);

    return {
        0: onlyWrongNameProperty,
        1: onlyWrongNameProperty2,
        2: test2,
    };

    // migrate inn poolnames to intrests (from 'interest' and 'participant-interests')

    // migrate in name from pardicipantFirstName to participantFirstName  (from 'interest' and 'participant-interests')
};
