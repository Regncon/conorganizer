import { getParticipantByUser } from '$app/(authorized)/my-profile/my-tickets/components/lib/actions/actions';
import { FirebaseCollectionNames } from '$lib/enums';
import { adminDb, getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import { Filter } from 'firebase-admin/firestore';

export const getAssignedGameByDay = async (participantId: string) => {
    const { user } = await getAuthorizedAuth();
    if (user === null) {
        return null;
    }
    const myParticipants = await getParticipantByUser();

    const participantFilter = myParticipants
        .filter((participant) => participant.id === participantId)
        ?.map((participant) => Filter.where('participantId', '==', participant.id));
    console.log('participantFilter', participantFilter);

    const currentGamesForParticipantId = (
        await adminDb
            .collection(FirebaseCollectionNames.players)
            .where(Filter.or(...participantFilter))
            .get()
    ).docs.map((doc) => doc.data());
    console.log('currentGames', currentGamesForParticipantId);
    return currentGamesForParticipantId;
};
