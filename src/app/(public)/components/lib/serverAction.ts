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

const initialSortedMap = new Map<PoolName, Set<PoolEvent>>([
    [PoolName.fridayEvening, new Set()],
    [PoolName.saturdayMorning, new Set()],
    [PoolName.saturdayEvening, new Set()],
    [PoolName.sundayMorning, new Set()],
]);

export async function getAllPoolEvents() {
    const allPoolEventsRef = await adminDb.collection('pool-events').get();
    const allPoolEvents = allPoolEventsRef.docs.map((doc) => ({ id: doc.id, ...doc.data() })) as PoolEvent[];

    const sortedPoolEventDay = allPoolEvents.reduce((acc, poolEvent) => {
        const currentAcc = acc.get(poolEvent.poolName);
        if (currentAcc !== undefined) {
            const existPoolEvent = [...currentAcc].find((event) => event.id === poolEvent.id);
            if (existPoolEvent === undefined) {
                acc.set(poolEvent.poolName, currentAcc.add(poolEvent));
            }
        }
        return acc;
    }, new Map<PoolName, Set<PoolEvent>>(initialSortedMap));
    return sortedPoolEventDay;
}
export async function getAdjacentPoolEventsById(id: string, day: PoolName) {
    const poolDayEvents = await getAllPoolEvents();
    const poolEventSet = poolDayEvents.get(day);

    if (poolEventSet) {
        const poolEvents = [...poolEventSet];
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
