'use server';
import { getEventById, getPoolEventById } from '$app/(public)/components/lib/serverAction';
import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import { ConEvent, PoolEvent } from '$lib/types';
import { doc, updateDoc } from 'firebase/firestore';
import { console } from 'inspector';

export async function updateEnventAndPoolEvent(eventId: string, incomingData: Partial<ConEvent>) {
    console.log('incomingData: ', incomingData);
    const { db, user } = await getAuthorizedAuth();
    if (db === null || user === null) {
        return;
    }
    console.log('incomingData: ', incomingData);
    const conEvent: ConEvent = await getEventById(eventId);
    // update the event with the incoming data

    console.log('code is here');
    try {
        await updateDoc(doc(db, 'events', eventId), incomingData);
        console.log('Document updated');
    } catch (e) {
        console.error('Error updating event document: ', e);
    }

    // update all pool events belonging to the event with the incoming data

    // convert incomingData to PoolEvent
    const incomingDataPool: Partial<PoolEvent> = {
        published: incomingData.published,
        title: incomingData.title,
        gameMaster: incomingData.gameMaster,
        system: incomingData.system,
        shortDescription: incomingData.shortDescription,
        description: incomingData.description,
        smallImageURL: incomingData.smallImageURL,
        bigImageURL: incomingData.bigImageURL,
        gameType: incomingData.gameType,
        isSmallCard: incomingData.isSmallCard,
        participants: incomingData.participants,
        childFriendly: incomingData.childFriendly,
        possiblyEnglish: incomingData.possiblyEnglish,
        adultsOnly: incomingData.adultsOnly,
        lessThanThreeHours: incomingData.lessThanThreeHours,
        moreThanSixHours: incomingData.moreThanSixHours,
        beginnerFriendly: incomingData.beginnerFriendly,
        additionalComments: incomingData.additionalComments,
        createdAt: incomingData.createdAt,
        createdBy: incomingData.createdBy,
        updateAt: incomingData.updateAt,
        updatedBy: incomingData.updatedBy,
    };
    // console.log(incomingData, 'incomingDataPool: ');

    await Promise.all(
        conEvent.poolIds.map(async (pool) => {
            const poolEvent: PoolEvent = await getPoolEventById(pool.id);
            if (!poolEvent) {
                console.error('PoolEvent not found');
                return;
            }
            try {
                await updateDoc(doc(db, 'pool-events', pool.id), incomingData);
                console.log('Document updated');
            } catch (e) {
                console.error('Error poolEvent updating document: ', e);
            }
        })
    );
}
