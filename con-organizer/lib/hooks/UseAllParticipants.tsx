import { useEffect, useState } from 'react';
import { Subscription } from 'rxjs';
import { Participant } from '../../models/types';
import { allParticipants$ } from '../observable';

export const useAllParticipants = (userId?: string) => {
    const [participants, setParticipants] = useState<Participant[]>();
    const [loadingParticipants, setLoadingParticipants] = useState<boolean>(true);

    useEffect(() => {
        let eventsObservable: Subscription;
        if (userId)
            eventsObservable = allParticipants$(userId).subscribe((participants) => {
                setParticipants(participants as Participant[] | undefined);
                setLoadingParticipants(false);
            });

        return () => {
            eventsObservable.unsubscribe();
        };
    }, []);
    return { participants, loadingParticipants };
};
