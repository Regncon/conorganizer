import { NextResponse } from 'next/server';
import { adminDb } from '@/lib/firebaseAdmin';
import { FirebaseCollections } from '@/models/enums';
import { EventTicket, Participant, UserSettings } from '@/models/types';

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

    /*     queryResult.data.eventTickets
        .filter((crm) => crm.category_id !== 116907)
        .forEach(async (crm) => {
            const newParticipant: Participant = {
                externalId: crm.id.toString(),
                name: `${crm.crm.first_name} ${crm.crm.last_name}` || '',
                email: crm.crm.email || '',
                connectedUser: '',
                eventTicket: crm,
            };
            console.log(newParticipant);
        }); */

    if (queryResult?.errors || !queryResult?.data?.eventTickets) {
        return NextResponse.json({ errors: queryResult?.errors }, { status: 403 });
    }

    const userSettingsInFirebase = ((await GetUserSettingsFromFirebase()) as UserSettings[]) || [];

    const migratedParticipants: Participant[] = [];

    userSettingsInFirebase.forEach(async (user) => {
        if (!user.name) {
            console.error(user.name, user.id, ' has no name');
            return;
        }

        if (user.checkInId) {
            //console.log(user.name, user.id, ' already has checkIn id', user.checkInId);
            const participants = await GetParticipantsFromFirebaseUserSettings(user.id);

            if (participants.find((p) => p.externalId.toString() === user.checkInId) === undefined) {
                await addPrimaryParticipant(queryResult, user);
            }

            if (participants) {
                const secondaryParticipants = queryResult.data.eventTickets
                    .filter((crm) => crm.category_id !== 116907)
                    .filter((q) => q.crm.email === participants.find((p) => p.isPrimary)?.email)
                    .filter((crm) => crm.id.toString() !== user.checkInId.toString())
                    .filter((p) => !participants.find((p2) => p2.externalId.toString() === p.id.toString()));

                //console.log(secondaryParticipants, 'secondaryParticipants');

                //make sure there are no duplicates in secondaryParticipants where id is the same as externalId in participants
                

                if (secondaryParticipants.length > 0) {
                    console.log(secondaryParticipants, 'secondaryParticipants');
                    secondaryParticipants.forEach(async (secondaryParticipant) => {
                        const secondaryParticipantDocument: Participant = {
                            externalId: secondaryParticipant.id || 0,
                            connectedUser: user.id,
                            isPrimary: false,
                            primaryParticipantId: user.checkInId,
                            name: `${secondaryParticipant.crm.first_name} ${secondaryParticipant.crm.last_name}` || '',
                            email: secondaryParticipant.crm.email || '',
                            eventTicket: secondaryParticipant,
                            orderId: secondaryParticipant.order_id || 0,
                        };
                        console.log('updating', user.id, 'setting with', secondaryParticipantDocument);
                                                                         await adminDb
                             .collection(
                                `${FirebaseCollections.userSetting}/${user.id}/${FirebaseCollections.Participants}/`
                            )
                            .add(secondaryParticipantDocument);  
                    });
                }
            }

            return;
        }

        const checkInTicket = queryResult.data.eventTickets.find(
            (crm) => user.name === `${crm.crm.first_name} ${crm.crm.last_name}`
        );
        if (checkInTicket) {
            console.log(user.name, ' updating setting with checking id', checkInTicket.id.toString());
            /*             await adminDb.doc(`${FirebaseCollections.userSetting}/${user.id}`).update({
                checkInId: checkInTicket.id.toString(),
            });
            return; */
        }

        console.error(user.name, user.id, ' not found in checkin');
        return;
    });

    /*         const participants = await GetParticipantsFromFirebaseUserSettings(user.id);

        participants.map((participant) => {
            if (participant) {
                console.log(participant, 'Old participant found');
                migratedParticipants.push(participant);
            }
        }); */

    /*             console.log(user.name, ' already exists');

            const participants = await GetParticipantsFromFirebaseUserSettings(user.id);

            const existingParticipant = participants.find((p) => p.externalId === checkInTicket.id.toString());

            if (existingParticipant) {
                console.log(existingParticipant, 'Old participant found');
            } else {
                const newParticipant: Participant = {
                    externalId: checkInTicket.id.toString(),
                    name: `${checkInTicket.crm.first_name} ${checkInTicket.crm.last_name}` || '',
                    email: checkInTicket.crm.email || '',
                    connectedUser: '',
                    eventTicket: checkInTicket,
                };
                migratedParticipants.push(newParticipant);
                console.log(newParticipant, 'Adding new participant to user', user.name);
            }
            */
    //console.log(newParticipant, 'newParticipant');

    /*     queryResult.data.eventTickets
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

                console.log(newParticipant, 'newParticipant');
            }
        });
 */
    //await adminDb.collection(FirebaseCollections.Participants).doc().set(newParticipants[0]);

    //console.log(migratedParticipants, 'migratedParticipants');
    return NextResponse.json({ migratedParticipants }, { status: 200 });
};

async function addPrimaryParticipant(queryResult: CrmJson, user: UserSettings) {
    const primaryParticipant = queryResult.data.eventTickets
        .filter((crm) => crm.category_id !== 116907)
        .find((crm) => crm.id.toString() === user.checkInId) as EventTicket;

    const primaryParticipantDocument: Participant = {
        externalId: primaryParticipant?.id || 0,
        connectedUser: user.id,
        isPrimary: true,
        name: `${primaryParticipant?.crm.first_name} ${primaryParticipant?.crm.last_name}` || '',
        email: primaryParticipant?.crm.email || '',
        eventTicket: primaryParticipant,
        orderId: primaryParticipant?.order_id || 0,
    };

    console.log('updating', user.id, 'setting with', primaryParticipantDocument);
         await adminDb
        .collection(`${FirebaseCollections.userSetting}/${user.id}/${FirebaseCollections.Participants}/`)
        .add(primaryParticipantDocument);
    return primaryParticipantDocument; 
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
