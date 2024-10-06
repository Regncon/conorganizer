'use server';
import { adminDb } from '$lib/firebase/firebaseAdmin';
import { revalidatePath } from 'next/cache';
import type { ConEvent, Participant, PoolEvent, InterestsInPool, Interest } from '$lib/types';
import { PoolName, InterestLevel } from '$lib/enums';
import { createIconOptions } from './helpers/icons';
import { FieldPath, Filter } from 'firebase-admin/firestore';
import type { MyUserInfo } from '$app/(authorized)/my-events/lib/types';

export async function getAllEvents() {
    const eventRef = await adminDb.collection('events').get();
    const events = eventRef.docs.map((doc) => ({ id: doc.id, ...doc.data() })) as ConEvent[];
    return events;
}

export async function GetAllUsers() {
    const userRef = await adminDb.collection('users').get();
    const users = userRef.docs.map((doc) => ({ id: doc.id, ...doc.data() }));
    return users;
}

export async function GetAllParticipants() {
    const participantsRef = await adminDb.collection('participants').get();
    const participants = participantsRef.docs.map((doc) => ({ id: doc.id, ...doc.data() }));
    return participants;
}
export async function GetAllParticipantsSnapshot() {
    const participantsRef = await adminDb.collection('participants').get();
    const participants = participantsRef;
    return participants;
}

const initialSortedEventMap = new Map<PoolName, ConEvent[]>([
    [PoolName.fridayEvening, []],
    [PoolName.saturdayMorning, []],
    [PoolName.saturdayEvening, []],
    [PoolName.sundayMorning, []],
]);
export async function getAllEventsSortedByPoolName() {
    const events = await getAllEvents();
    const sortedEventDay = events.reduce((accumulator, poolEvent) => {
        poolEvent.poolIds.forEach((poolChildRef) => {
            const currentAccumulator = accumulator.get(poolChildRef.poolName);

            if (currentAccumulator !== undefined) {
                accumulator.set(poolChildRef.poolName, [...currentAccumulator.values(), poolEvent]);
            }
        });
        return accumulator;
    }, new Map<PoolName, ConEvent[]>(initialSortedEventMap));
    return sortedEventDay;
}

export type PoolEvents = Awaited<ReturnType<typeof getAllPoolEvents>>;

const initialSortedPoolMap = new Map<PoolName, PoolEvent[]>([
    [PoolName.fridayEvening, []],
    [PoolName.saturdayMorning, []],
    [PoolName.saturdayEvening, []],
    [PoolName.sundayMorning, []],
]);

export async function getAllPoolEvents() {
    const allPoolEventsRef = await adminDb.collection('pool-events').get();
    const allPoolEvents = allPoolEventsRef.docs.map((doc) => ({ id: doc.id, ...doc.data() })) as PoolEvent[];

    const sortedPoolEventDay = allPoolEvents.reduce((accumulator, poolEvent) => {
        const currentAccumulator = accumulator.get(poolEvent.poolName);

        if (currentAccumulator !== undefined) {
            accumulator.set(poolEvent.poolName, [...currentAccumulator.values(), poolEvent]);
        }

        return accumulator;
    }, new Map<PoolName, PoolEvent[]>(initialSortedPoolMap));

    return sortedPoolEventDay;
}
export async function getAdjacentPoolEventsById(id: string, day: PoolName) {
    const poolDayEvents = await getAllPoolEvents();
    const getPoolEventsByDay = poolDayEvents.get(day);

    if (getPoolEventsByDay) {
        const poolEvents = getPoolEventsByDay.filter((event) => event.published);
        const eventIndex = poolEvents.findIndex((event) => event.id === id);
        const prevNavigationId = poolEvents[eventIndex - 1]?.id;
        const nextNavigationId = poolEvents[eventIndex + 1]?.id;
        return { prevNavigationId, nextNavigationId };
    }

    return { prevNavigationId: undefined, nextNavigationId: undefined };
}
export async function getPoolEventById(id: string) {
    const poolEvent = (await adminDb.collection('pool-events').doc(id).get()).data() as PoolEvent;

    let icons = createIconOptions(
        poolEvent.adultsOnly,
        poolEvent.childFriendly,
        poolEvent.beginnerFriendly,
        poolEvent.lessThanThreeHours,
        poolEvent.moreThanSixHours,
        poolEvent.possiblyEnglish,
        poolEvent.gameType
    );
    poolEvent.icons = icons;
    return { ...poolEvent, id };
}

export async function getEventInterestById(id: string) {
    const event = (await adminDb.collection('events').doc(id).get()).data() as ConEvent;
    if (!event) {
        return;
    }

    let poolInterests: InterestsInPool[] = [];

    await Promise.all(
        event.poolIds.map(async (pool) => {
            const interestRef = await adminDb.collection('pool-events').doc(pool.id).collection('interests').get();
            const interests = interestRef.docs.map((doc) => ({ id: doc.id, ...doc.data() })) as Interest[];

            const poolInterest: InterestsInPool = {
                poolId: pool.id,
                poolName: pool.poolName,
                interests: interests,
            };

            poolInterests.push(poolInterest);
        })
    );

    console.log('poolInterests: ', poolInterests);

    return poolInterests;
}

export async function getUsersInterestById(id: string) {
    const participants = (await GetAllParticipants()) as Participant[];
    const myParticipants = participants.filter((participant) => participant.users?.includes(id));

    const participantIdToFilter = myParticipants
        .map((user) => user.id)
        .map((participantId) => {
            return Filter.where('participantId', '==', participantId);
        })
        .filter((filter) => filter !== undefined);

    const filterForAllParticipants = Filter.or(...participantIdToFilter);

    const userParticipants = (
        await adminDb.collectionGroup('interests').where(filterForAllParticipants).where('interestLevel', '>', 0).get()
    ).docs
        .filter((doc) => doc.ref.path.includes('pool-events'))
        .map((doc) => ({ id: doc.id, ...doc.data() })) as Interest[];

    return userParticipants;
}

export async function migrateInterestsToParticipantInterests() {
    console.log('starting migration');

    const participantsRef = await adminDb.collection('participants').get();
    const participants = participantsRef.docs.map((doc) => ({ id: doc.id, ...doc.data() })) as Participant[];
    const allPoolEventsRef = await adminDb.collection('pool-events').get();
    const allPoolEvents = allPoolEventsRef.docs.map((doc) => ({ id: doc.id, ...doc.data() })) as PoolEvent[];
    participants.forEach(async (participant) => {
        if (participant.id === 'FqxDY3CjlGfLqW5PlfSA') {
            const interestRef = await adminDb
                .collection('participants')
                .doc(participant.id as string)
                .collection('interests')
                .get();
            const interests = interestRef.docs.map((doc) => ({ id: doc.id, ...doc.data() })) as Interest[];

            interests.forEach(async (interest) => {
                //if interestLevet is not of type InterestLevel, set it to NotInterested
                const interestLevelToSet =
                    Object.values(InterestLevel).includes(interest.interestLevel) ?
                        interest.interestLevel
                    :   InterestLevel.NotInterested;

                // console.log('interest: ', interest.id);
                const newInterest = { ...interest };
                delete newInterest.id;
                newInterest.interestLevel = interestLevelToSet;
                newInterest.poolEventTitle = allPoolEvents.find((event) => event.id === interest.poolEventId)?.title;
                newInterest.poolName = allPoolEvents.find((event) => event.id === interest.poolEventId)?.poolName;
                newInterest.updateAt = new Date().toISOString();
                newInterest.updatedBy = 'migrate Interests To Participant Interests';

                // console.log('interestWithouId: ', interestWithouId);

                // await adminDb
                //     .collection('participants')
                //     .doc(participant.id as string)
                //     .collection('interests')
                //     .doc(interest.poolEventId as string)
                //     .set(newInterest);

                await adminDb
                    .collection('participants')
                    .doc(participant.id as string)
                    .collection('participant-interests')
                    .doc(interest.poolEventId as string)
                    .set(newInterest);

                await adminDb
                    .collection('pool-events')
                    .doc(interest.poolEventId)
                    .collection('interests')
                    .doc(participant.id as string)
                    .set(newInterest);
            });
        }
    });
}

export async function getEventById(id: string) {
    const event = (await adminDb.collection('events').doc(id).get()).data() as ConEvent;
    return { ...event, id };
}
export async function getMyEventById(id: string, userId: string) {
    const event = (
        await adminDb.collection('users').doc(userId).collection('my-events').doc(id).get()
    ).data() as ConEvent;
    return { ...event, id };
}

export async function updateEvents() {
    revalidatePath('/', 'page');
}
export async function updateEventById(id: string) {
    revalidatePath(`/event/${id}`, 'page');
}

export async function updateDashboardEvents() {
    revalidatePath('/admin/dashboard/events', 'page');
}

export async function getAllParticipants() {
    const participantsRef = await adminDb.collection('participants').get();
    const participants = participantsRef.docs.map((doc) => ({ id: doc.id, ...doc.data() })) as Participant[];
    return participants;
}
