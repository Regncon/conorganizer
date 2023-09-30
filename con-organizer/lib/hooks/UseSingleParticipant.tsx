import { useEffect, useState } from 'react';
import { Participant } from '../../models/types';
import { singleParticipant$ } from '../observable';

export const useSingleParticipants = (id: string) => {
    const [participant, setParticipant] = useState<Participant>();
    const [loading, setLoading] = useState<boolean>(true);

    useEffect(() => {
        const participantObservable = singleParticipant$(id).subscribe((participant) => {
            setParticipant(participant as Participant);
            setLoading(false);
        });

        return () => {
            participantObservable.unsubscribe();
        };
    }, []);

    return { participant, loading };
};