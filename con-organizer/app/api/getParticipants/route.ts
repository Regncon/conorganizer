import { NextResponse } from 'next/server';
import { skip } from 'node:test';
import { adminDb } from '@/lib/firebaseAdmin';
import { FirebaseCollections } from '@/models/enums';
import { CrmJson, Participant } from '@/models/types';

// const url = https://app.checkin.no/graphql?client_id=API_KEY&client_secret=API_SECRET
const query = `{
  allCrms(customerId: 13446) {
    records
    data {
      id
      firstName
      lastName
      email {
        email
        id
      }
    }
  }
}`;

export const GET = async () => {
    //ToDo: authenticate request

    const queryResult: CrmJson | undefined = await GetParticipantsFromCheckIn();

    if (queryResult?.errors || !queryResult?.data?.allCrms) {
        return NextResponse.json({ errors: queryResult?.errors }, { status: 403 });
    }

    const participantInFirebase = await GetParticipantFromFirebase();
    console.log(participantInFirebase, 'participantInFirebase');

    const newParticipants: Participant[] = [];

    queryResult.data.allCrms.data.forEach((crm) => {
        if (participantInFirebase.find((p) => p.externalId === crm.id.toString())) {
            console.log('already exists');
        }
        else {
            const newParticipant: Participant = {
                externalId: crm.id.toString(),
                name: `${crm.firstName} ${crm.lastName}` || '',
                email: crm.email?.email || '',
                connectedUser: '',
            };
            newParticipants.push(newParticipant);
            console.log(newParticipant, 'newParticipant');
        }
    });

    //await adminDb.collection(FirebaseCollections.Participants).doc().set(newParticipants[0]);

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
