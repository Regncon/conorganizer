import { useEffect, useState } from 'react';
import { Subscription } from 'rxjs';
import { Participant } from '../../models/types';
import { allParticipants$ } from '../observable';

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
