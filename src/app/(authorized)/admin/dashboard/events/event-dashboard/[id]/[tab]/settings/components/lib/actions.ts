'use server';
import { getEventById, getPoolEventById } from '$app/(public)/components/lib/serverAction';
import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import { ConEvent, PoolEvent } from '$lib/types';
import { doc, updateDoc } from 'firebase/firestore';

export async function updatePoolEvent(eventId: string, incomingData: Partial<ConEvent>) {
    console.log('incomingData: ', incomingData);
    const { db, user } = await getAuthorizedAuth();
    if (db === null || user === null) {
        return;
    }

    const conEvent: ConEvent = await getEventById(eventId);

    await Promise.all(
        conEvent.poolIds.map(async (pool) => {
            const poolEvent: PoolEvent = await getPoolEventById(pool.id);
            if (!poolEvent) {
                console.error('PoolEvent not found');
                return;
            }
            try {
                await updateDoc(doc(db, 'pool-events', pool.id), incomingData);
                console.log('Pool events updated');
            } catch (e) {
                console.error('Error poolEvent updating document: ', e);
            }
        })
    );
}
