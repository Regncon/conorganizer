'use server';
import { adminDb } from '$lib/firebase/firebaseAdmin';
import { revalidatePath } from 'next/cache';
import type { ConEvent, PoolEvent } from '$lib/types';
import type { PoolName } from '$lib/enums';

export async function getAllEvents() {
    const eventRef = await adminDb.collection('events').get();
    const events = eventRef.docs.map((doc) => ({ id: doc.id, ...doc.data() })) as ConEvent[];
    return events;
}

type PoolEventDay = {
    day: PoolName;
    poolEvents: PoolEvent[];
};

export async function getAllPoolEventsSortedByDay() {
    const allPoolEventsRef = await adminDb.collection('pool-events').get();
    const allPoolEvents = allPoolEventsRef.docs.map((doc) => ({ id: doc.id, ...doc.data() })) as PoolEvent[];
    // console.log(allPoolEvents);

    const poolEvents = allPoolEvents.reduce(
        (acc, event) => {
            if (acc[2]?.day && event.poolName === 'fridayEvening') {
                console.log(acc[2]?.poolEvents, event);
            }
            if (acc.some((day) => day.day === event.poolName) === false) {
                return [...acc, { day: [event.poolName], poolEvents: [event] }] as PoolEventDay[];
            }
            const day = acc.find((day) => day.day === event.poolName);
            if (day !== undefined) {
                day.poolEvents.push(event);
            }
            return acc as PoolEventDay[];
        },
        [] as unknown as PoolEventDay[]
    );
    return poolEvents;
}

export async function getEventById(id: string) {
    const event = (await adminDb.collection('events').doc(id).get()).data() as ConEvent;
    return { ...event, id };
}

export async function getPoolEventById(id: string) {
    const event = (await adminDb.collection('pool-events').doc(id).get()).data() as PoolEvent;
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
