'use server';

import { geMyEventByRefPath } from '$app/(authorized)/my-events/lib/actions';
import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import { ConEvent, MyNewEvent } from '$lib/types';
import { addDoc, collection, doc, updateDoc } from 'firebase/firestore';

export async function convertToConEvent(myEventId: string, myEventUserId: string) {
    const { db, user } = await getAuthorizedAuth();
    if (db === null || user === null) {
        return;
    }
    const myEvent: MyNewEvent = await geMyEventByRefPath(db, myEventId, myEventUserId);
    console.log('My event', myEvent);
    if (myEvent.isAccepted) {
        console.log('Event', myEventId, 'is already approved');
        return;
    }
    console.log('Approving event', myEventId, 'for user', myEventUserId);

    const event: ConEvent = {
        name: myEvent.name,
        title: myEvent.title,
        subTitle: myEvent.subTitle,
        description: myEvent.description,
        gameType: myEvent.gameType,
        system: myEvent.system,
        gameMaster: myEvent.name,
        published: false,
        shortDescription: '',
        email: myEvent.email,
        phone: myEvent.phone,
        roomIds: [],
        poolIds: [],
        isSmallCard: false,
        participants: myEvent.participants,
        puljeFridayEvening: false,
        puljeSaturdayMorning: false,
        puljeSaturdayEvening: false,
        puljeSundayMorning: false,
        unwantedFridayEvening: !myEvent.fridayEvening,
        unwantedSaturdayMorning: !myEvent.saturdayMorning,
        unwantedSaturdayEvening: !myEvent.saturdayEvening,
        unwantedSundayMorning: !myEvent.sundayMorning,
        moduleCompetition: myEvent.moduleCompetition,
        childFriendly: myEvent.childFriendly,
        possiblyEnglish: myEvent.possiblyEnglish,
        adultsOnly: myEvent.adultsOnly,
        volunteersPossible: myEvent.volunteersPossible,
        lessThanThreeHours: myEvent.lessThanThreeHours,
        moreThanSixHours: myEvent.moreThanSixHours,
        beginnerFriendly: myEvent.beginnerFriendly,
        additionalComments: myEvent.additionalComments,
        createdAt: Date.now().toString(),
        createdBy: user.uid,
        updateAt: Date.now().toString(),
        updatedBy: user.uid,
    };

    console.log('Event', event);

    try {
        const docRef = await addDoc(collection(db, 'events'), event);
        console.log('Document written with ID: ', docRef.id);
        myEvent.eventDocId = docRef.id;
        myEvent.isAccepted = true;
    } catch (e) {
        console.error('Error adding document: ', e);
        return;
    }

    try {
        await updateDoc(doc(db, 'users', myEventUserId, 'my-events', myEventId), myEvent);
        console.log('Document updated');
    } catch (e) {
        console.error('Error updating document: ', e);
    }
}
