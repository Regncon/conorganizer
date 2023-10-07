import { type NextRequest, NextResponse } from 'next/server';
import { adminDb } from '@/lib/firebaseAdmin';
import { EnrollmentOptions, FirebaseCollections, Pool } from '@/models/enums';
import { EnrollmentChoice } from '@/models/types';

type Payload = {
    pool: Pool;
};

export const POST = async (request: NextRequest) => {
    const payload: Payload = await request.json();

    console.log(payload, 'payload');

    const EnrollmentChoices: EnrollmentChoice[] = await getAllEventParticipants();

    const previousEnrollmentChoices = await getAllEnrollmentChoices();

    console.log(previousEnrollmentChoices);

    //console.log(EnrollmentChoices, 'EnrollmentChoices');

    const userSettings = await GetUserSettingsFromFirebase();

    const participants = await getParticipants();

    //console.log(userSettings, 'userSettings');

    const events = await getEvents();

    //console.log(events, 'events');

    events
        .filter((event) => event.pool === payload.pool)
        .map(async (event) => {
            EnrollmentChoices.filter((enrollmentChoice) => enrollmentChoice.eventId === event.id)
                .filter((enrollmentChoice) =>
                    userSettings.find((userSetting) => userSetting.id === enrollmentChoice.userId)
                )
                .filter((enrollmentChoice) => enrollmentChoice.participantId !== enrollmentChoice.userId)
                .forEach(async (enrollmentChoice) => {
                    console.log('adding enrollmentChoice to event ', event.id, enrollmentChoice.participantId);

                    const participant = participants.find(
                        (participant) => participant.id === enrollmentChoice.participantId
                    );

                    enrollmentChoice.ticketType = participant?.eventTicket?.category ?? '';

                    enrollmentChoice.firstChoiceEventTitle =
                        events?.find((event) => event?.id === enrollmentChoice?.firstChoiceEventId)?.title ?? '';

                    enrollmentChoice.eventTitle = event.title;

                    const firstChoiceEvent =
                        previousEnrollmentChoices.find(
                            (previousEnrollmentChoice) =>
                                previousEnrollmentChoice.participantId === enrollmentChoice.participantId &&
                                previousEnrollmentChoice.choice === EnrollmentOptions.VeryInterested &&
                                previousEnrollmentChoice.isEnrolled === true
                        ) ?? null;

                    if (firstChoiceEvent) {
                        enrollmentChoice.hasGotFirstChoice = true;
                        enrollmentChoice.firstChoiceEventId = firstChoiceEvent.eventId;
                        enrollmentChoice.firstChoiceEventTitle = firstChoiceEvent.eventTitle;
                    }
                    else {
                        enrollmentChoice.hasGotFirstChoice = false;
                        enrollmentChoice.firstChoiceEventId = '';
                        enrollmentChoice.firstChoiceEventTitle = '';
                    }

                    console.log(enrollmentChoice);
                    await adminDb
                        .collection(
                            `${FirebaseCollections.events}/${event.id}/${FirebaseCollections.EnrollmentChoices}`
                        )
                        .doc(enrollmentChoice.participantId)
                        .set(enrollmentChoice);
                });
        });

    return NextResponse.json({}, { status: 200 });
};

async function getAllEventParticipants() {
    const EnrollmentChoices: EnrollmentChoice[] = [];

    const eventEnrollmentsRef = await adminDb.collectionGroup(FirebaseCollections.EventParticipants).get();

    eventEnrollmentsRef.docs.map((doc) => {
        const eventId = doc.ref.path.split('/')[1];
        const userId = doc.ref.path.split('/')[3];

        const data = doc.data();
        data.participantId = doc.id;
        data.eventId = eventId;
        data.userId = userId;

        EnrollmentChoices.push(data as EnrollmentChoice);
    });
    return EnrollmentChoices;
}

async function getAllEnrollmentChoices() {
    const EnrollmentChoices: EnrollmentChoice[] = [];

    const eventEnrollmentsRef = await adminDb.collectionGroup(FirebaseCollections.EnrollmentChoices).get();

    return eventEnrollmentsRef.docs.map((doc) => {
        const data = doc.data();
        data.id = doc.id;
        return data;
    });
}

async function getParticipants() {
    const participants = await adminDb.collectionGroup(FirebaseCollections.Participants).get();

    return participants.docs.map((doc) => {
        const data = doc.data();
        data.id = doc.id;
        return data;
    });
}

async function GetUserSettingsFromFirebase() {
    const userSettingsFirebaseRef = await adminDb.collection(FirebaseCollections.userSetting).get();

    return userSettingsFirebaseRef.docs.map((doc) => {
        const data = doc.data();
        data.id = doc.id;
        return data;
    });
}

async function getEvents() {
    const eventsRef = await adminDb.collection(FirebaseCollections.events).get();

    return eventsRef.docs.map((doc) => {
        const data = doc.data();
        data.id = doc.id;
        return data;
    });
}
