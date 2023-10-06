import { NextResponse } from 'next/server';
import { adminDb } from '@/lib/firebaseAdmin';
import { FirebaseCollections } from '@/models/enums';
import { EventTicket, Participant, UserSettings } from '@/models/types';
import { use } from 'react';

// const url = https://app.checkin.no/graphql?client_id=API_KEY&client_secret=API_SECRET
const query = `{
    eventTickets(customer_id: 13446, id: 58182, onlyCompleted: true) {
      id
      category
      category_id
      crm {
        first_name
        last_name
        id
        email
      }
      order_id
    }
  }`;

type CrmJson = {
    data: {
        eventTickets: EventTicket[];
    };
    errors: Error;
};
export const GET = async () => {
    console.log('getParticipants staring');

    const queryResult: CrmJson | undefined = await GetParticipantsFromCheckIn();

    if (queryResult?.errors || !queryResult?.data?.eventTickets) {
        return NextResponse.json({ errors: queryResult?.errors }, { status: 403 });
    }

    const userSettingsInFirebase = ((await GetUserSettingsFromFirebase()) as UserSettings[]) || [];

    userSettingsInFirebase.forEach(async (user) => {
        const participants = await GetParticipantsFromFirebaseUserSettings(user.id);

        if (!participants || participants.length === 0) {
            return;
        }
        // find duplicate participants that have the same externalId
        const duplicateParticipants = participants.filter(
            (p) => participants.filter((p2) => p2.externalId === p.externalId).length > 1
        );
        if (duplicateParticipants.length > 0) {
            console.error(user.name, user.id, ' has duplicate participants');
            console.error(duplicateParticipants);
        }
    });

    return NextResponse.json({}, { status: 200 });
};

async function GetParticipantsFromCheckIn() {
    const res = await fetch(
        `https://app.checkin.no/graphql?client_id=${process.env.CHECKIN_KEY}&client_secret=${process.env.CHECKIN_SECRET}`,
        {
            method: 'POST',
            body: JSON.stringify({ query }),
            headers: {
                'Content-Type': 'application/json',
            },
        }
    );

    const queryResult: CrmJson | undefined = await res.json();

    return queryResult;
}

async function GetUserSettingsFromFirebase() {
    const userSettingsFirebaseRef = await adminDb.collection(FirebaseCollections.userSetting).get();

    return userSettingsFirebaseRef.docs.map((doc) => {
        const data = doc.data();
        data.id = doc.id;
        return data;
    });
}

async function GetParticipantsFromFirebaseUserSettings(userId: string) {
    const participantsFirebaseRef = await adminDb
        .collection(`${FirebaseCollections.userSetting}/${userId}/${FirebaseCollections.Participants}`)
        .get();

    return participantsFirebaseRef.docs.map((doc) => doc.data()) as Participant[];
}
