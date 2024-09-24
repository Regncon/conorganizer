'use server';
import { adminDb } from '$lib/firebase/firebaseAdmin';
import { revalidatePath } from 'next/cache';
import type { ConEvent, Participant, PoolEvent } from '$lib/types';
import { PoolName } from '$lib/enums';
import { createIconOptions } from './helpers/icons';

export async function getAllEvents() {
    const eventRef = await adminDb.collection('events').get();
    const events = eventRef.docs.map((doc) => ({ id: doc.id, ...doc.data() })) as ConEvent[];
    return events;
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
        const poolEvents = getPoolEventsByDay;
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
        poolEvent.possiblyEnglish
    );
    poolEvent.icons = icons;
    return { ...poolEvent, id };
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
