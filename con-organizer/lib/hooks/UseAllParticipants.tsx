import { useEffect, useState } from 'react';
import { collection } from 'firebase/firestore';
import { collectionData } from 'rxfire/firestore';
import { Subscription } from 'rxjs';
import { Participant } from '../../models/types';
import { db } from '../firebase';
export function allParticipants$(userId: string) {
    return collectionData(participantsRef(userId), { idField: 'id' });
}
export const participantsRef = (userId: string) => collection(db, `usersettings/${userId}/participants/`);
export const useAllParticipants = (userId?: string) => {
    const [participants, setParticipants] = useState<Participant[]>();
    const [loadingParticipants, setLoadingParticipants] = useState<boolean>(true);

    useEffect(() => {
        let allParticipantsObservable: Subscription;
        //console.log(userId, 'userId');
        if (userId)
            allParticipantsObservable = allParticipants$(userId).subscribe((participants) => {
                setParticipants(participants as Participant[] | undefined);
                setLoadingParticipants(false);

                return () => {
                    allParticipantsObservable.unsubscribe();
                };
            });
        return () => {
            null;
        };
    }, []);
    return { participants, loadingParticipants };
};
