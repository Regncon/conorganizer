import { NextResponse } from 'next/server';
import { adminDb } from '@/lib/firebaseAdmin';
import { FirebaseCollections } from '@/models/enums';
import { EventTicket, Participant } from '@/models/types';

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
    //ToDo: authenticate request

    const queryResult: CrmJson | undefined = await GetParticipantsFromCheckIn();

    if (queryResult?.errors || !queryResult?.data?.eventTickets) {
        return NextResponse.json({ errors: queryResult?.errors }, { status: 403 });
    }

    const participantInFirebase = await GetParticipantFromFirebase();

    const newParticipants: Participant[] = [];

    queryResult.data.eventTickets
    .filter((crm) => crm.category_id !== 116907)
    .forEach((crm) => {
        if (participantInFirebase.find((p) => p.externalId === crm.id.toString())) {
            console.log('already exists');
        }
        else {
            const newParticipant: Participant = {
                externalId: crm.id.toString(),
                name: `${crm.crm.first_name} ${crm.crm.last_name}` || '',
                email: crm.crm.email || '',
                connectedUser: '',
                eventTicket: crm,
            };
            newParticipants.push(newParticipant);
            console.log(newParticipant, 'newParticipant');
        }
    });

    await adminDb.collection(FirebaseCollections.Participants).doc().set(newParticipants[0]);

    return NextResponse.json({ newParticipants }, { status: 200 });
};

async function GetParticipantFromFirebase() {
    const participantInFirebaseRef = await adminDb.collection(FirebaseCollections.Participants).get();

    return participantInFirebaseRef.docs.map((doc) => doc.data()) as Participant[];
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
