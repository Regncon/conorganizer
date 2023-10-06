import { NextResponse } from 'next/server';
import {
    addPrimaryParticipant,
    addSecondaryParticipants,
    GetParticipantsFromCheckIn,
    GetParticipantsFromFirebaseUserSettings,
    GetUserSettingsFromFirebase,
} from '@/lib/apiHelpers';
import { CrmJson, Participant, UserSettings } from '@/models/types';

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
            let primaryParticipantDocument: Participant | undefined;

            if (participants.find((p) => p.externalId.toString() === user.checkInId) === undefined) {
                primaryParticipantDocument = await addPrimaryParticipant(queryResult, user);
            }

            //if (participants) {
            addSecondaryParticipants(queryResult, primaryParticipantDocument, user, participants);
            //}

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
