import { FirebaseCollectionNames } from '$lib/enums';
import { getAuthorizedAuth } from '$lib/firebase/firebaseAdmin';
import { PoolPlayer } from '$lib/types';
import { collection, getDocs } from 'firebase/firestore';

export async function getAssignetPlayers(poolId: string) {
    const { db } = await getAuthorizedAuth();
    if (!db) {
        throw new Error('Database is undefined');
    }
    const poolPlayersRef = collection(
        db,
        FirebaseCollectionNames.poolEvents,
        poolId,
        FirebaseCollectionNames.poolPlayers
    );
    const poolPlayersSnapshot = await getDocs(poolPlayersRef);
    const poolPlayers = poolPlayersSnapshot.docs.map((doc) => ({ id: doc.id, ...doc.data() })) as PoolPlayer[];

    return poolPlayers;
}
