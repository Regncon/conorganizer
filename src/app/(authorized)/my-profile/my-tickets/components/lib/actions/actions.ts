'use server';
import { getMyUserInfo } from '$app/(authorized)/my-events/lib/actions';
import { GetAllParticipants, GetAllParticipantsSnapshot } from '$app/(public)/components/lib/serverAction';
import { adminDb, getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import { ActionResponse, Participant } from '$lib/types';
import { addDoc, collection, doc, getDocs, query, updateDoc, where } from 'firebase/firestore';
import { revalidatePath } from 'next/cache';

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

export const getParticipantByUser = async () => {
    const { user } = await getAuthorizedAuth();
    if (user === null) return null;

    console.log(
        adminDb
            .collection('participants')
            // .where('userId', 'array-contains', user.uid)
            .get()
            .then((querySnapshot) => {
                querySnapshot.forEach((doc) => {
                    console.log(doc.id, ' => ', doc.data());
                });
            })
            .catch((error) => {
                console.error('Error getting documents: ', error);
            })
    );
};
export const AssignParticipantByEmail = async () => {
    console.log('AssignParticipantByEmail');
    const { db, user } = await getAuthorizedAuth();
    if (db === null || user === null) {
        const response: ActionResponse = {
            type: 'error',
            message: 'Ikke autorisert',
            error: 'getAuthorizedAuth failed',
        };
        console.error(response);
        throw response;
    }
    console.log('user email', user?.email);

    const tickets = await GetTicketsByEmail(user?.email);
    console.log('tickets', tickets);

    // create participant from ticket

    const participants = await GetParticipantsByEmail(user?.email as string);
    const myParticipants = participants.filter((participant) => participant.users?.includes(user.uid));
    if (participants.length === myParticipants.length) {
        console.log('ticket connected do user YOU WIN');
        return participants;
    }

    const userInfo = (await getMyUserInfo(db, user)) ?? { admin: false, participantIds: [] };

    participants.forEach(async (participant) => {
        //TODO: HUSK OG TEST
        if (participant.users?.includes(user?.uid)) {
            console.log('ticket connected do user YOU WIN');
            return;
        }

        if (participant.users?.includes(user.uid) === false || participant?.users === undefined) {
            if (participant.users === undefined) {
                participant.users = [];
            }

            participant.users.push(user.uid);
        }

        if (userInfo.participantIds?.includes(participant.id) === false || userInfo.participantIds === undefined) {
            if (userInfo.participantIds === undefined) {
                userInfo.participantIds = [];
            }
            userInfo.participantIds.push(participant.id);
        }
        // console.log((await test.query.where('id', '==', participant.id).get()).forEach(e => e.ref.update(participant.users ?? [])), 'ASDASDASDASD ==============');
        adminDb.collection('participants').doc(participant.id).update(participant);
    });

    adminDb.collection('users').doc(user.uid).update(userInfo);

    return participants;
};

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

    // const query = `{
    // eventTickets(customer_id: 13446, id: ${ticketId}, onlyCompleted: true) {
    //   id
    //   category
    //   category_id
    //   crm {
    //     first_name
    //     last_name
    //     id
    //     email
    //   }
    //   order_id
    // }
    // }`;
    //
    // const res = await fetch(
    //     `https://app.checkin.no/graphql?client_id=${process.env.CHECKIN_KEY}&client_secret=${process.env.CHECKIN_SECRET}`,
    //     {
    //         method: 'POST',
    //         body: JSON.stringify({ query }),
    //         headers: {
    //             'Content-Type': 'application/json',
    //         },
    //     }
    // );
    // if (res.status !== 200) {
    //     throw new Error('Failed to fetch tickets');
    // }
    // const queryResult: CrmJson | undefined = await res.json();
    // console.log(queryResult, 'queryResult');
    //
    // // const result = await ConvertTicketToParticipant(queryResult?.data.eventTickets[0] as EventTicket);
    // return result;
};

const ConvertTicketToParticipant = async (ticketId: number, tickets: EventTicket[]) => {
    const { db, user } = await getAuthorizedAuth();
    if (db === null || user === null) {
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

    const ticket = tickets.find((ticket) => ticket.id === ticketId);
    if (!ticket) throw new Error('ticket not found');

    const participantRef = collection(db, 'participants');
    const q = query(participantRef, where('ticketId', '==', ticketId));
    const querySnapshot = await getDocs(q);

    if (!querySnapshot.empty) {
        console.log('Participant already exists:', querySnapshot.docs[0].data());
        const response: ActionResponse = {
            type: 'warning',
            message: 'Deltageren er allerede registrert',
        };
        return response;
    }
    const isOver18 = new Date().getFullYear() - new Date(ticket.crm.born).getFullYear() > 18;

    const orderEmails = tickets.filter((t) => t.order_id === ticket.order_id).map((t) => t.crm.email);

    let participant: Partial<Participant> = {
        name: `${ticket.crm.first_name} ${ticket.crm.last_name}`,
        over18: isOver18,
        ticketEmail: ticket.crm.email,
        orderEmails: orderEmails,
        ticketId: ticket.id,
        orderId: ticket.order_id,
        ticketCategory: ticket.category,
        ticketCategoryId: ticket.category_id,
        createdAt: new Date().toISOString(),
        createdBy: user.email as string,
        updateAt: new Date().toISOString(),
        updatedBy: user.email as string,
    };
    console.log(participant, 'participant');

    let participantId = '';
    try {
        const docRef = await addDoc(collection(db, 'participants'), participant);
        console.log('Document written with ID: ', docRef.id);
        participantId = docRef.id;
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
