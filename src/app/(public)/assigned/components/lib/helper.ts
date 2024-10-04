import { getParticipantByUser } from '$app/(authorized)/my-profile/my-tickets/components/lib/actions/actions';
import { FirebaseCollectionNames, PoolName } from '$lib/enums';
import { adminDb, getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import type { PoolPlayer } from '$lib/types';
import { Filter } from 'firebase-admin/firestore';

export const getAssignedGameByDay = async (participantId: string) => {
    const { user } = await getAuthorizedAuth();
    if (user === null) {
        console.warn('No user logged in');
        return { poolName: null, currentGamesForParticipant: null };
    }
    const myParticipants = await getParticipantByUser();

    const participantFilter = myParticipants
        .filter((participant) => participant.id === participantId)
        ?.map((participant) => Filter.where('participantId', '==', participant.id));
    const currentGamesForParticipantId = (
        await adminDb
            .collection(FirebaseCollectionNames.players)
            .where(Filter.or(...participantFilter))
            .get()
    ).docs.map((doc) => ({ id: doc.id, ...doc.data() }) as PoolPlayer);
    console.warn(currentGamesForParticipantId, 'found currentGamesForParticipantId');
    const now = new Date();
    const day = now.getDay(); // 0 = Sunday, 1 = Monday, ..., 5 = Friday, 6 = Saturday
    const hours = now.getHours();
    const isMorning = hours < 12;

    let poolName: PoolName;

    switch (day) {
        case 5: // Friday
            poolName = PoolName.fridayEvening;
            break;
        case 6: // Saturday
            poolName = isMorning ? PoolName.saturdayMorning : PoolName.saturdayEvening;
            break;
        case 0: // Sunday
            poolName = isMorning ? PoolName.sundayMorning : PoolName.fridayEvening;
            break;
        default:
            poolName = PoolName.fridayEvening;
    }

    const findForPoolName = currentGamesForParticipantId.find((player) => player.poolName === poolName);
    console.warn(
        findForPoolName,
        findForPoolName ?
            'found currentGamesForParticipantId for poolName'
        :   'no currentGamesForParticipantId for poolName'
    );
    return {
        currentGamesForParticipant: findForPoolName,
        poolName,
    };
};
