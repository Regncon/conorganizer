import { NextResponse } from 'next/server';
import { adminDb } from '@/lib/firebaseAdmin';
import { FirebaseCollections } from '@/models/enums';
import { Enrollment, EventTicket, Participant, UserSettings } from '@/models/types';

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

    const eventEnrollments = await getEventEnrollments();

    //console.log('eventEnrollments', eventEnrollments);

    eventEnrollments.forEach(async (enrollmentItem) => {
        //console.log('enrollment', enrollment);
        const participant = await GetParticipantsFromFirebaseUserSettings(enrollmentItem.id).then((participants) => {
            return participants.find((p) => p.isPrimary);
        });
        //console.log('participants', participant);

        const eventParticipantRef = `${FirebaseCollections.events}/${enrollmentItem.eventId}/${FirebaseCollections.Enrollments}/${enrollmentItem.userId}/${FirebaseCollections.EventParticipants}/`;
        //console.log(eventParticipantRef, 'eventParticipantRef');

        const eventEnrollment: Enrollment = {
            choice: enrollmentItem.choice,
            name: participant?.name || '',
            isPrimary: participant?.isPrimary || false,
        };

        //console.log(eventEnrollment, 'eventEnrollment');

        console.log('adding event enrollment', eventParticipantRef, 'and id', enrollmentItem.userId, 'setting with', eventEnrollment);

        await adminDb.collection(eventParticipantRef).doc(enrollmentItem.userId).set(eventEnrollment);
    });

    /*         await adminDb
            .collection(
                `${FirebaseCollections.events}/${enrollment.eventId}/${FirebaseCollections.Enrollments}/${enrollment.userId}/${FirebaseCollections.EventParticipants}/`
            )
            .doc(participant.userId)
            .set(
                choice: enrollment,
                name: participant?.name,
                isPrimary: participant?.isPrimary,
            );

            await adminDb
            .collection(FirebaseCollections.userSetting)
            .doc(userRecord.uid)
            .set({
                name: `${paidUser.crm.first_name} ${paidUser.crm.last_name}`.trim(),
                checkInId: paidUser.id.toString(),
                */
    /*         export const participantEnrollmentsRef = (eventId: string, userId: string, participantId: string) =>
               doc(db, `events/${eventId}`, `/enrollments/${userId}/eventParticipants/${participantId}`);
            }); */

    /*     userSettingsInFirebase.forEach(async (user) => {
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
                }); */

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

async function getEventEnrollments() {
    const eventEnrollmentsRef = await adminDb
        //.collectionGroup(`${FirebaseCollections.events}/${eventId}/${FirebaseCollections.Enrollments}`)
        .collectionGroup(FirebaseCollections.Enrollments)
        .get();

    return eventEnrollmentsRef.docs.map((doc) => {
        //console.log(doc.ref.path, 'doc');
        // events/ydfeJ8VIuQ4tIyKVjuCo/enrollments/t6Y0JPt9kmMUDxxxS1vDFAhHkbJ3
        // split on / and get the 2rd element
        const eventId = doc.ref.path.split('/')[1];
        const userId = doc.ref.path.split('/')[3];

        const data = doc.data();
        data.id = doc.id;
        data.eventId = eventId;
        data.userId = userId;
        return data as Enrollment[];
    });
}
