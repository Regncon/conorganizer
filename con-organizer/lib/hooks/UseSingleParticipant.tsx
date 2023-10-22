import { useEffect, useState } from 'react';
import { doc } from 'firebase/firestore';
import { docData } from 'rxfire/firestore';
import { Subscription } from 'rxjs';
import { Participant } from '../../models/types';
import { db } from '../firebase';
export const participantRef = (userId: string, participantId: string) =>
    doc(db, `usersettings/${userId}/participants/${participantId}`);
export function singleParticipant$(userId: string, participantId: string) {
    return docData(participantRef(userId, participantId), { idField: 'id' });
}
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
