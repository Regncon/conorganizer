import { useEffect, useState } from 'react';
import { Participant } from '../../models/types';
import { allParticipants$ } from '../observable';

export const useAllParticipants = () => {
    const [participants, setParticipants] = useState<Participant[]>();
    const [loadingParticipants, setLoadingParticipants] = useState<boolean>(true);

    useEffect(() => {
        const eventsObservable = allParticipants$.subscribe((participants) => {
            setParticipants(participants as Participant[] | undefined);
            setLoadingParticipants(false);
        });

        return () => {
            eventsObservable.unsubscribe();
        };
    }, []);
    return { participants, loadingParticipants };
};
