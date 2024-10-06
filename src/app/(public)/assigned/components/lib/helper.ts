import { getParticipantByUser } from '$app/(authorized)/my-profile/my-tickets/components/lib/actions/actions';
import { FirebaseCollectionNames, PoolName } from '$lib/enums';
import { adminDb, getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import type { PoolPlayer } from '$lib/types';
import { Filter } from 'firebase-admin/firestore';

export const getAssignedGameByDay = async (participantId: string) => {
    const { user } = await getAuthorizedAuth();
    if (user === null) {
        console.warn('No user logged in');
        return null;
    }
    const myParticipants = await getParticipantByUser();
    console.log(myParticipants, 'myParticipants');
    if (myParticipants.length === 0) {
        return [];
    }

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

    currentGamesForParticipantId;
    console.warn(
        currentGamesForParticipantId,
        currentGamesForParticipantId ?
            'found currentGamesForParticipantId for poolName'
        :   'no currentGamesForParticipantId for poolName'
    );
    return currentGamesForParticipantId;
};
