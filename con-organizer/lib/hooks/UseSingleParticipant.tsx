import { useEffect, useState } from 'react';
import { Subscription } from 'rxjs';
import { Participant } from '../../models/types';
import { singleParticipant$ } from '../observable';

export const useSingleParticipants = (id?: string, participantId?: string) => {
    const [participant, setParticipant] = useState<Participant>();
    const [loadingParticipant, setLoading] = useState<boolean>(true);

    useEffect(() => {
        let participantObservable: Subscription;
        if (id && participantId) {
            singleParticipant$(id, participantId).subscribe((participant) => {
                setParticipant(participant as Participant);
                setLoading(false);
            });
        }

        return () => {
            participantObservable.unsubscribe();
        };
    }, []);

    return { participant, loadingParticipant };
};
