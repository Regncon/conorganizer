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

    // ----------------------------------------------------------------
    //missing poolName in interests
    // ----------------------------------------------------------------
    // console.log(
    //     'GetAllInterests',
    //     GetAllInterests.filter((interest) => interest.poolName === undefined)
    // );
    // console.log(
    //     'GetAllInterests2',
    //     GetAllInterests2.filter((interest) => interest.poolName === undefined)
    // );

    // // missing firstname
    // console.log(
    //     'GetAllInterests',
    //     GetAllInterests.filter((interest) => interest.participantFirstName === undefined)
    // );
    // console.log(
    //     'GetAllInterests2',
    //     GetAllInterests2.filter((interest) => interest.participantFirstName === undefined)
    // );

    // const interestOnlyWrongNameProperty = GetAllInterests.filter(Boolean).filter(
    //     (interest) => interest.pardicipantFirstName !== undefined
    // );
    // const interestOnlyWrongNameProperty2 = GetAllInterests2.filter(Boolean).filter(
    //     (interest) => interest.pardicipantFirstName !== undefined
    // );

    // fix name issues in interests
    // const interestOnlyWrongNamePropertyUpdateDb = interestOnlyWrongNameProperty.map((interest) => {
    //     return adminDb
    //         .collection('participants')
    //         .doc(interest.participantId)
    //         .collection('interests')
    //         .doc(interest.id ?? '')
    //         .update({ participantFirstName: interest.pardicipantFirstName });
    // });

    // const interestOnlyWrongNameProperty2UpdateDb = interestOnlyWrongNameProperty2.map((interest) => {
    //     return adminDb
    //         .collection('participants')
    //         .doc(interest.participantId)
    //         .collection('participant-interests')
    //         .doc(interest.id ?? '')
    //         .update({ participantFirstName: interest.pardicipantFirstName });
    // });

    // return {
    //     0: interestOnlyWrongNameProperty,
    //     1: interestOnlyWrongNameProperty2,
    // };

    // ----------------------------------------------------------------
    // migrate inn poolnames to intrests (from 'interest')
    // ----------------------------------------------------------------

    // const interestOnlyWrongNameProperty = GetAllInterests.filter(Boolean).filter(
    //     (interest) => interest.poolName === undefined
    // );

    // const allPoolEvents = (await adminDb.collection('pool-events').get()).docs.map((doc) => ({
    //     id: doc.id,
    //     ...doc.data(),
    // })) as PoolEvent[];

    // const allPoolEventsWithInterestsPoolNamesThatArMissing = allPoolEvents.filter((poolEvent) =>
    //     interestOnlyWrongNameProperty.some((interest) => interest.poolEventId === poolEvent.id)
    // );
    // const interestOnlyWrongNamePropertyUpdateDb = interestOnlyWrongNameProperty.map((interest) => {
    //     const poolEvent = allPoolEvents.find((poolEvent) => poolEvent.id === interest.poolEventId);
    //     if (poolEvent === undefined) {
    //         return;
    //     }
    //     return adminDb
    //         .collection('participants')
    //         .doc(interest.participantId)
    //         .collection('interests')
    //         .doc(interest.id ?? '')
    //         .update({ poolName: poolEvent.poolName, updatedBy: 'migrate interests and poolNames' });
    // });

    // const updateResult = await Promise.all(interestOnlyWrongNamePropertyUpdateDb);
    // console.log(updateResult);

    // const updateTest = await adminDb
    //     .collection('participants')
    //     .doc('0bNx792Iz1qk2eKS5OZx')
    //     .collection('interests')
    //     .doc('MnmNgvVQuay4S5B3mMLG')
    //     // .update({ poolName: 'fridayEvening' });
    //     .get()
    //     .then((doc) => doc.data());
    // console.log(updateTest);
    // return {
    //     0: interestOnlyWrongNameProperty,
    //     // 1: interestOnlyWrongNameProperty2,
    //     2: allPoolEvents,
    //     3: allPoolEventsWithInterestsPoolNamesThatArMissing,
    // };

    // ----------------------------------------------------------------
    // migrate in poolName from 'poolEvents' to 'interests in pool-event'
    // ----------------------------------------------------------------

    // const GetAllPoolEventsInterests = (await adminDb.collectionGroup('interests').get()).docs.map((doc) => {
    //     if (doc.ref.path.includes('pool-events')) {
    //         return {
    //             id: doc.id,
    //             ...doc.data(),
    //         };
    //     }
    //     return undefined;
    // }) as Interest[];
    // const poolEventInterestOnlyWrongNameProperty = GetAllPoolEventsInterests.filter(Boolean).filter(
    //     (interest) => interest.poolName === undefined
    // );

    // const allPoolEvents = (await adminDb.collection('pool-events').get()).docs.map((doc) => ({
    //     id: doc.id,
    //     ...doc.data(),
    // })) as PoolEvent[];

    // const allPoolEventsWithPoolEventsInterestsPoolNamesThatArMissing = allPoolEvents.filter((poolEvent) =>
    //     poolEventInterestOnlyWrongNameProperty.some((interest) => interest.poolEventId === poolEvent.id)
    // );
    // const interestOnlyWrongNamePropertyUpdateDb = poolEventInterestOnlyWrongNameProperty.map((interest) => {
    //     const poolEvent = allPoolEvents.find((poolEvent) => poolEvent.id === interest.poolEventId);
    //     if (poolEvent === undefined) {
    //         return;
    //     }
    //     return adminDb
    //         .collection('pool-events')
    //         .doc(interest.poolEventId)
    //         .collection('interests')
    //         .doc(interest.id ?? '')
    //         .update({ poolName: poolEvent.poolName, updatedBy: 'migrate interests and poolNames' });
    // });

    // const updateResult = await Promise.all(interestOnlyWrongNamePropertyUpdateDb);
    // console.log(updateResult);

    // const updateTest = await adminDb
    //     .collection('pool-events')
    //     .doc('qZmV7evJ749YjyAiRKMe')
    //     .collection('interests')
    //     .doc('BhTrf9kH5nzy4rxbV4ru')
    //     // .update({ poolName: 'fridayEvening', updatedBy: 'migrate interests and poolNames' });
    // .get()
    // .then((doc) => doc.data());
    // console.log(updateTest);

    // return {
    //     0: poolEventInterestOnlyWrongNameProperty,
    //     // 1: interestOnlyWrongNameProperty2,
    //     2: allPoolEvents,
    //     3: allPoolEventsWithPoolEventsInterestsPoolNamesThatArMissing,
    // };
    // migrate in name from pardicipantFirstName to participantFirstName  (from 'interest' and 'participant-interests')
};
