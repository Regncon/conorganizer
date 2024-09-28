'use server';
import { getMyUserInfo } from '$app/(authorized)/my-events/lib/actions';
import { GetAllParticipants } from '$app/(public)/components/lib/serverAction';
import { adminDb, getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import { ActionResponse, Interest, Participant } from '$lib/types';
import { doc, updateDoc } from 'firebase/firestore';
import { revalidatePath } from 'next/cache';
import { AssignUserToParticipant, generateParticipant, NewTickets } from './Helpers';
import { InterestLevel, PoolName } from '$lib/enums';

export type CrmRecord = {
    id: number;
    firstName: string;
    lastName: string;
    hash: string;
    email: {
        email: string;
        id: number;
    };
};

export type CrmData = {
    records: number;
    data: CrmRecord[];
};

export type CrmJson = {
    data: {
        eventTickets: EventTicket[];
    };
    errors: Error;
};

export type EventTicket = {
    id: number;
    category: string;
    category_id: number;
    order_id: number;
    crm: {
        first_name: string;
        last_name: string;
        id: number;
        email: string;
        born: string;
    };
};

export const updateInterest = async (
    participantId?: string,
    poolEventId?: string,
    poolName?: PoolName,
    poolEventTitle?: string,
    interestLevel?: InterestLevel
) => {
    const { user } = await getAuthorizedAuth();
    if (
        user === null ||
        user.email === undefined ||
        user.uid === undefined ||
        participantId === undefined ||
        poolEventId === undefined ||
        poolEventTitle === undefined ||
        poolName === undefined ||
        interestLevel === undefined
    ) {
        console.error(
            'updateInterest failed with missing parameters',
            'participantId:',
            participantId,
            'poolEventId:',
            poolEventId,
            'poolName:',
            poolName,
            'poolEventTitle:',
            poolEventTitle,
            'interestLevel:',
            interestLevel
        );
        return;
    }

    const participant = (await adminDb.collection('participants').doc(participantId).get()).data() as Participant;
    if (participant === undefined) {
        return;
    }

    if (participant.users?.includes(user.uid) === false) {
        return;
    }

    const newInterest: Interest = {
        interestLevel: interestLevel,
        poolEventId: poolEventId,
        poolEventTitle: poolEventTitle,
        poolName: poolName,
        participantId: participantId,
        participantFirstName: participant.firstName,
        participantLastName: participant.lastName,
        updateAt: new Date().toISOString(),
        updatedBy: user.uid || '',
    };

    await adminDb
        .collection('pool-events')
        .doc(poolEventId)
        .collection('interests')
        .doc(participantId)
        .set(newInterest);

    await adminDb
        .collection('participants')
        .doc(participantId)
        .collection('interests')
        .doc(poolEventId)
        .set(newInterest);

    await adminDb
        .collection('participants')
        .doc(participantId)
        .collection('participant-interests')
        .doc(poolEventId)
        .set(newInterest);
};

export const getInterest = async (participantId?: string, poolEventId?: string) => {
    console.log('getInterest', participantId, poolEventId);

    const { user } = await getAuthorizedAuth();
    if (
        user === null ||
        user.email === undefined ||
        user.uid === undefined ||
        participantId === undefined ||
        poolEventId === undefined
    ) {
        return;
    }

    const interest = (
        await adminDb.collection('pool-events').doc(poolEventId).collection('interests').doc(participantId).get()
    ).data() as Interest;

    console.log('interest', interest);

    if (interest === undefined) {
        return InterestLevel.NotInterested;
    }
    return interest.interestLevel;
};

export const getParticipantByUser = async () => {
    const { user } = await getAuthorizedAuth();
    if (user === null) return;

    // console.log(
    //     adminDb
    //         .collection('participants')
    //         // .where('userId', 'array-contains', user.uid)
    //         .get()
    //         .then((querySnapshot) => {
    //             querySnapshot.forEach((doc) => {
    //                 console.log(doc.id, ' => ', doc.data());
    //             });
    //         })
    //         .catch((error) => {
    //             console.error('Error getting documents: ', error);
    //         })
    // );
    const participants = (await GetAllParticipants()) as Participant[];
    return participants.filter((participant) => participant.users?.includes(user.uid));
};

export const GetMyParticipants = async () => {
    const { user } = await getAuthorizedAuth();
    if (user === null) {
        throw new Error('Failed to get authorized auth');
    }
    const participants = (await GetAllParticipants()) as Participant[];
    const newParticipants = participants.filter((participant) => participant.users?.includes(user.uid));
    return newParticipants.map((participant, i) => {
        return {
            id: participant.id,
            firstName: participant.firstName,
            lastName: participant.lastName,
            isSelected: i === 0 ? true : false,
        };
    });
};

export const AssignParticipantByEmail = async () => {
    const { db, user } = await getAuthorizedAuth();
    if (db === null || user === null) {
        throw new Error('Failed to get authorized auth');
    }

    const tickets = await GetTicketsByEmail(user?.email);
    if (!tickets) {
        throw new Error('Failed to get tickets from Checkin');
    }

    const participants = (await GetAllParticipants()) as Participant[];

    const newParticipants = NewTickets(tickets, participants, user);

    let newParticipantIds: string[] = [];
    newParticipants.forEach(async (newParticipant) => {
        try {
            const docRef = await adminDb.collection('participants').add(newParticipant);
            newParticipantIds.push(docRef.id);
            console.log('Participant written with ID: ', docRef.id);
        } catch (e) {
            throw e;
        }
        console.log('newParticipantIds', newParticipantIds);
    });

    const updatedParticipants = AssignUserToParticipant(participants, user);

    if (newParticipantIds.length === 0 && updatedParticipants.length === 0) {
        console.log('No participants to update');
        return participants.filter((participant) => participant.users?.includes(user.uid));
    }

    updatedParticipants.forEach(async (participant) => {
        const participantId = participants.find((p) => p.ticketId === participant.ticketId)?.id;
        try {
            await adminDb
                .collection('participants')
                .doc(participantId as string)
                .update(participant);
        } catch (e) {
            throw e;
        }
    });

    const myUserInfo = await getMyUserInfo(db, user);
    let myUserInfoToBeUpdated = myUserInfo ? { ...myUserInfo } : { admin: false, participantIds: [] };

    if (myUserInfoToBeUpdated.participantIds === undefined) {
        myUserInfoToBeUpdated.participantIds = [];
    }

    if (newParticipantIds.length > 0) {
        myUserInfoToBeUpdated.participantIds = myUserInfoToBeUpdated.participantIds.concat(newParticipantIds);
    }

    if (updatedParticipants.length > 0) {
        updatedParticipants.forEach((participant) => {
            if (myUserInfoToBeUpdated.participantIds?.includes(participant.id as string) === false) {
                myUserInfoToBeUpdated.participantIds.push(participant.id as string);
            }
        });
    }

    if (myUserInfo) {
        adminDb.collection('users').doc(user.uid).update(myUserInfoToBeUpdated);
    }
    if (!myUserInfo) {
        adminDb.collection('users').doc(user.uid).set(myUserInfoToBeUpdated);
    }
    const participantsBelongingToUser = participants
        .concat(newParticipants)
        .concat(participants.filter((p) => updatedParticipants.some((up) => up.id === p.id)));

    return participantsBelongingToUser;
};

// let newParticipantIds: string[] = [];
// newTickets.forEach(async (newTicket) => {
//     let newParticipant = generateParticipant(newTicket.id, tickets, user.email as string);
//     newParticipant = { ...newParticipant, users: [user.uid] };
//
//     try {
//         const docRef = await adminDb.collection('participants').add(newParticipant);
//         newParticipantIds.push(docRef.id);
//         console.log('Participant written with ID: ', docRef.id);
//     } catch (e) {
//         throw e;
//     }
//     console.log('newParticipantIds', newParticipantIds);
// });
//
// const myParticipants = participants.filter((participant) => participant.users?.includes(user.uid));
// if (participants.length === myParticipants.length && newParticipantIds.length === 0) {
//     console.log('ticket connected do user YOU WIN');
//     return participants;
// }
//
// const myUserInfo = await getMyUserInfo(db, user);
// let myUserInfoToBeUpdated = myUserInfo ? { ...myUserInfo } : { admin: false, participantIds: [] };
//
// if (myUserInfoToBeUpdated.participantIds === undefined) {
//     myUserInfoToBeUpdated.participantIds = [];
// }
//
// if (newParticipantIds.length > 0) {
//     myUserInfoToBeUpdated.participantIds = myUserInfoToBeUpdated.participantIds.concat(newParticipantIds);
// }
//
// participants.forEach(async (participant) => {
//     if (participant.users?.includes(user?.uid)) {
//         console.log('ticket connected do user YOU WIN');
//         return;
//     }
//
//     if (participant.users?.includes(user.uid) === false || participant?.users === undefined) {
//         if (participant.users === undefined) {
//             participant.users = [];
//         }
//
//         participant.users.push(user.uid);
//     }
//
//     if (myUserInfoToBeUpdated.participantIds?.includes(participant.id) === false) {
//         myUserInfoToBeUpdated.participantIds.push(participant.id);
//     }
// });
//
// if (myUserInfo) {
//     adminDb.collection('users').doc(user.uid).update(myUserInfoToBeUpdated);
// }
// if (!myUserInfo) {
//     adminDb.collection('users').doc(user.uid).set(myUserInfoToBeUpdated);
// }
// console.log('completed assignParticipantByEmail');
// return participants;
// };

export const GetTicketsByEmail = async (email: string | null | undefined) => {
    if ((typeof email === 'string') == false) {
        return;
    }

    const allTickets = await GetTicketsFromCheckIn();

    const emailTickets = allTickets?.filter((ticket) => ticket.crm.email === email);
    const ticketsWithOrderNumberFromEmail = allTickets?.filter((ticket) =>
        emailTickets?.some((emailTicket) => emailTicket.order_id === ticket.order_id)
    );

    return ticketsWithOrderNumberFromEmail;
};

export const GetParticipantsByEmail = async (email: string) => {
    const allParticipants = (await GetAllParticipants()) as Participant[];

    const participants = allParticipants.filter(
        (participant) =>
            participant.ticketEmail === email ||
            participant.orderEmails.includes(email) ||
            (participant.connectedEmails && participant.connectedEmails.includes(email))
    );
    return participants;
};

export const ConvertTicketIdToParticipant = async (ticketId: number) => {
    console.log('ConvertTicketToParticipant ', ticketId);
    const tickets = await GetTicketsFromCheckIn();

    if (!tickets) {
        const response: ActionResponse = {
            type: 'error',
            message: 'Checkin feilet',
        };
        console.error(response);
        return response;
    }

    const ticket = tickets?.find((ticket) => ticket.id === ticketId);
    console.log(ticket, 'ticket');

    if (!ticket) {
        const response: ActionResponse = {
            type: 'error',
            message: 'Fant ikke billetten',
        };
        console.error(response);
        return response;
    }

    const result = await ConvertTicketToParticipant(ticket.id, tickets);
    return result;
};

const ConvertTicketToParticipant = async (ticketId: number, tickets: EventTicket[]) => {
    const { user } = await getAuthorizedAuth();
    if (user === null) {
        // return new ActionResponse = {
        //     type: "error"
        //     , message: 'Ikke autorisert'
        //     error: 'getAuthorizedAuth failed'
        const response: ActionResponse = {
            type: 'error',
            message: 'Ikke autorisert',
            error: 'getAuthorizedAuth failed',
        };
        console.error(response);
        return response;
    }

    const participantRef = adminDb.collection('participants').where('ticketId', '==', ticketId);
    const querySnapshot = await participantRef.get();

    if (!querySnapshot.empty) {
        console.log('Participant already exists:', querySnapshot.docs[0].data());
        const response: ActionResponse = {
            type: 'warning',
            message: 'Deltageren er allerede registrert',
        };
        return response;
    }

    const participant = generateParticipant(ticketId, tickets, user.email as string);

    try {
        const docRef = await adminDb.collection('participants').add(participant);
        console.log('Document written with ID: ', docRef.id);
    } catch (e) {
        const response: ActionResponse = {
            type: 'error',
            message: 'Feil ved lagring av deltager',
            error: e as string,
        };
        console.error(response);
        return response;
    }
    const response: ActionResponse = {
        type: 'success',
        message: 'Deltager lagret',
    };
    revalidatePath('/admin/dashboard/participants', 'page');
    return response;
};

export const GetTicketsFromCheckIn = async () => {
    const query = `{
    eventTickets(customer_id: 13446, id: 73685, onlyCompleted: true) {
      id
      category
      category_id
      crm {
        first_name
        last_name
        id
        email
        born
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
    //console.log(queryResult, 'queryResult');

    return queryResult?.data.eventTickets;
};

export const ConnectEmailToParticipant = async (participantId: string, email: string) => {
    console.log('ConnectEmailToParticipant ', participantId, email);
    const { db } = await getAuthorizedAuth();
    if (db === null) {
        const response: ActionResponse = {
            type: 'error',
            message: 'Ikke autorisert',
            error: 'getAuthorizedAuth failed',
        };
        console.error(response);
        return response;
    }

    const participant = (await adminDb.collection('participants').doc(participantId).get()).data() as Participant;

    if (participant === undefined) {
        const response: ActionResponse = {
            type: 'error',
            message: 'Fant ikke deltager',
        };
        console.error(response);
        return response;
    }
    if (participant.connectedEmails === undefined) {
        participant.connectedEmails = [];
    }

    if (participant.connectedEmails.includes(email)) {
        const response: ActionResponse = {
            type: 'warning',
            message: 'E-posten er allerede koblet til deltageren',
        };
        return response;
    }

    participant.connectedEmails.push(email);

    try {
        await updateDoc(doc(db, 'participants', participantId), participant);
    } catch (e) {
        const response: ActionResponse = {
            type: 'error',
            message: 'Feil ved lagring av deltager',
            error: e as string,
        };
        console.error(response);
        return response;
    }
    const response: ActionResponse = {
        type: 'success',
        message: 'Epost lagt til',
    };
    revalidatePath('/admin/dashboard/participants', 'page');
    return response;
};

export const DeleteConnectedEmail = async (participantId: string, email: string) => {
    console.log('DeleteConnectedEmail ', participantId, email);
    const { db } = await getAuthorizedAuth();
    if (db === null) {
        const response: ActionResponse = {
            type: 'error',
            message: 'Ikke autorisert',
            error: 'getAuthorizedAuth failed',
        };
        console.error(response);
        return response;
    }

    const participant = (await adminDb.collection('participants').doc(participantId).get()).data() as Participant;

    if (participant === undefined) {
        const response: ActionResponse = {
            type: 'error',
            message: 'Fant ikke deltager',
        };
        console.error(response);
        return response;
    }
    if (participant.connectedEmails === undefined) {
        participant.connectedEmails = [];
    }

    if (!participant.connectedEmails.includes(email)) {
        const response: ActionResponse = {
            type: 'warning',
            message: 'Eposten er ikke koblet til deltageren',
        };
        return response;
    }

    participant.connectedEmails = participant.connectedEmails.filter((e) => e !== email);

    try {
        await updateDoc(doc(db, 'participants', participantId), participant);
    } catch (e) {
        const response: ActionResponse = {
            type: 'error',
            message: 'Feil ved lagring av deltager',
            error: e as string,
        };
        console.error(response);
        return response;
    }
    const response: ActionResponse = {
        type: 'success',
        message: 'Epost fjernet',
    };
    revalidatePath('/admin/dashboard/participants', 'page');
    return response;
};
