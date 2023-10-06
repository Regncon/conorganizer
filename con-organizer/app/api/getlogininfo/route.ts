import { type NextRequest, NextResponse } from 'next/server';
import { adminDb, adminUser } from '@/lib/firebaseAdmin';
import { GetParticipantsFromCheckIn } from '@/lib/helpers';
import { FirebaseCollections } from '@/models/enums';
import { CrmJson, Participant } from '@/models/types';

type Payload = {
    email: string;
    password: string;
};
export const POST = async (request: NextRequest) => {
    const payload: Payload = await request.json();
    //console.log(payload, 'userInput');

    try {
        await adminUser.getUserByEmail(payload.email);
    } catch (error) {
        const queryResult: CrmJson | undefined = await GetParticipantsFromCheckIn();

        if (!queryResult || queryResult?.errors) {
            return NextResponse.json({ errors: queryResult?.errors }, { status: 403 });
        }

        const paidUser = queryResult?.data?.eventTickets
            .find((crm) => crm?.crm?.email === payload.email && crm?.category_id !== 116907);

        if (paidUser) {
            const userRecord = await adminUser.createUser({
                email: payload.email,
                password: payload.password,
                displayName: `${paidUser.crm.first_name} ${paidUser.crm.last_name}`,
            });
            const customToken = await adminUser.createCustomToken(userRecord.uid);
            console.log(customToken);

            const primaryParticipantDocument: Participant = {
                externalId: paidUser?.id || 0,
                connectedUser: userRecord.uid,
                isPrimary: true,
                primaryParticipantId: paidUser?.id.toString(),
                name: `${paidUser?.crm.first_name} ${paidUser?.crm.last_name}` || '',
                email: paidUser?.crm.email || '',
                eventTicket: paidUser,
                orderId: paidUser?.order_id || 0,
            };

            console.log('adding user settings to ', userRecord.uid);
            
            await adminDb
                .collection(FirebaseCollections.userSetting)
                .doc(userRecord.uid)
                .set({
                    name: `${paidUser.crm.first_name} ${paidUser.crm.last_name}`.trim(),
                    checkInId: paidUser.id.toString(),
                });


            console.log('adding primary participant', userRecord.uid, 'setting with', primaryParticipantDocument);

            await adminDb
                .collection(`${FirebaseCollections.userSetting}/${userRecord.uid}/${FirebaseCollections.Participants}/`)
                .add(primaryParticipantDocument);

            const secondaryParticipants = queryResult.data.eventTickets
                .filter((crm) => crm.category_id !== 116907)
                .filter((q) => q.crm.email === primaryParticipantDocument?.email)
                .filter((crm) => crm.id !== primaryParticipantDocument?.externalId)

            if (secondaryParticipants.length > 0) {
                secondaryParticipants.forEach(async (secondaryParticipant) => {
                    const secondaryParticipantDocument: Participant = {
                        externalId: secondaryParticipant.id || 0,
                        connectedUser: userRecord.uid,
                        isPrimary: false,
                        primaryParticipantId: primaryParticipantDocument.externalId.toString(),
                        name: `${secondaryParticipant.crm.first_name} ${secondaryParticipant.crm.last_name}` || '',
                        email: secondaryParticipant.crm.email || '',
                        eventTicket: secondaryParticipant,
                        orderId: secondaryParticipant.order_id || 0,
                    };
                    console.log('adding secondary participant', userRecord.uid, 'setting with', secondaryParticipantDocument);
                    await adminDb
                        .collection(
                            `${FirebaseCollections.userSetting}/${userRecord.uid}/${FirebaseCollections.Participants}/`
                        )
                        .add(secondaryParticipantDocument);
                });
            }

            return NextResponse.json({ user: 'created' }, { status: 200 });
        }
        return NextResponse.json({ user: `Don't exist` }, { status: 200 });
    }
    return NextResponse.json({ user: 'Exists' }, { status: 200 });
};
