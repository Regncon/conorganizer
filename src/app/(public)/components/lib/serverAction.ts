'use server';
import { adminDb } from '$lib/firebase/firebaseAdmin';
import { revalidatePath } from 'next/cache';
import type { ConEvent, PoolEvent } from '$lib/types';
import { PoolName } from '$lib/enums';
import { cache } from 'react';

export async function getAllEvents() {
    const eventRef = await adminDb.collection('events').get();
    const events = eventRef.docs.map((doc) => ({ id: doc.id, ...doc.data() })) as ConEvent[];
    return events;
}

export type PoolEvents = Awaited<ReturnType<typeof getAllPoolEvents>>;

const initialSortedMap = new Map<PoolName, PoolEvent[]>([
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
    }, new Map<PoolName, PoolEvent[]>(initialSortedMap));

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
    const event = (await adminDb.collection('pool-events').doc(id).get()).data() as PoolEvent;
    return { ...event, id };
}

export async function getEventById(id: string) {
    const event = (await adminDb.collection('events').doc(id).get()).data() as ConEvent;
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
