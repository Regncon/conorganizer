import { use } from 'react';
import { NextResponse } from 'next/server';
import { adminDb } from '@/lib/firebaseAdmin';
import { FirebaseCollections } from '@/models/enums';
import { EventTicket, Participant, UserSettings } from '@/models/types';

// const url = https://app.checkin.no/graphql?client_id=API_KEY&client_secret=API_SECRET
const query = `{
    eventTickets(customer_id: 13446, id: 58182) {
      id
      category
      category_id
      crm {
        first_name
        last_name
        id
        email
      }
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

    const userSettingsInFirebase = await GetUserSettingsFromFirebase();

    const migratedParticipants: Participant[] = [];

    /*     userSettingsInFirebase.forEach(async (user) =>
    {
        const participants = await GetParticipantsFromFirebaseUserSettings(user.id);

        participants.map((participant) => {
            if (participant) {
                console.log(participant, 'Old participant found');
                migratedParticipants.push(participant);
            }
        });

    });
 */

    queryResult.data.eventTickets
        .filter((crm) => crm.category_id !== 116907)
        .forEach(async (crm) => {
            
            const user = userSettingsInFirebase.find(
                (p) => p.name === `${crm.crm.first_name} ${crm.crm.last_name}`.trim()
            );

            if (user) {
                console.log(user.name, ' already exists');

                const participants = await GetParticipantsFromFirebaseUserSettings(user.id);

                const existingParticipant = participants.find((p) => p.externalId === crm.id);

                if (existingParticipant) {
                    console.log(existingParticipant, 'Old participant found');
                } else {
                    const newParticipant: Participant = {
                        externalId: crm.id.toString(),
                        name: `${crm.crm.first_name} ${crm.crm.last_name}` || '',
                        email: crm.crm.email || '',
                        connectedUser: '',
                        eventTicket: crm,
                    };
                    migratedParticipants.push(newParticipant);
                    console.log(newParticipant, 'Adding new participant to user', user.name);
                }

                //console.log(newParticipant, 'newParticipant');
            }

        });

    //await adminDb.collection(FirebaseCollections.Participants).doc().set(newParticipants[0]);

    //console.log(migratedParticipants, 'migratedParticipants');
    return NextResponse.json({ migratedParticipants }, { status: 200 });
};

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
