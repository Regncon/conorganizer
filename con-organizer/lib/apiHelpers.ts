import { FirebaseCollections } from '@/models/enums';
import { CrmJson, EventTicket, Participant, UserSettings } from '@/models/types';
import { adminDb } from './firebaseAdmin';

export async function addPrimaryParticipant(queryResult: CrmJson, user: UserSettings) {
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

    console.log('adding primary participant', user.id, 'setting with', primaryParticipantDocument);
    await adminDb
        .collection(`${FirebaseCollections.userSetting}/${user.id}/${FirebaseCollections.Participants}/`)
        .add(primaryParticipantDocument);
    return primaryParticipantDocument;
}

export async function GetUserSettingsFromFirebase() {
    const userSettingsFirebaseRef = await adminDb.collection(FirebaseCollections.userSetting).get();

    return userSettingsFirebaseRef.docs.map((doc) => {
        const data = doc.data();
        data.id = doc.id;
        return data;
    });
}

export async function GetParticipantsFromFirebaseUserSettings(userId: string) {
    const participantsFirebaseRef = await adminDb
        .collection(`${FirebaseCollections.userSetting}/${userId}/${FirebaseCollections.Participants}`)
        .get();

    return participantsFirebaseRef.docs.map((doc) => doc.data()) as Participant[];
}

export async function GetParticipantsFromCheckIn() {
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

export function addSecondaryParticipants(
    queryResult: CrmJson,
    primaryParticipantDocument: Participant | undefined,
    user: UserSettings,
    participants: any
) {
    const secondaryParticipants = queryResult.data.eventTickets
        .filter((crm) => crm.category_id !== 116907)
        .filter((q) => q.crm.email === primaryParticipantDocument?.email)
        .filter((crm) => crm.id !== primaryParticipantDocument?.externalId)
        //.filter((q) => q.crm.email === participants.find((p) => p.isPrimary)?.email)
        .filter((crm) => crm.id.toString() !== user.checkInId.toString())
        .filter((p) => !participants.find((p2) => p2.externalId.toString() === p.id.toString()));

    //console.log(secondaryParticipants, 'secondaryParticipants');
    //make sure there are no duplicates in secondaryParticipants where id is the same as externalId in participants
    if (secondaryParticipants.length > 0) {
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
            console.log('adding secondary participant', user.id, 'setting with', secondaryParticipantDocument);
            await adminDb
                .collection(`${FirebaseCollections.userSetting}/${user.id}/${FirebaseCollections.Participants}/`)
                .add(secondaryParticipantDocument);
        });
    }
}