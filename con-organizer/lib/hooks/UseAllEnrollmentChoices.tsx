import { useEffect, useState } from 'react';
import { Subscription } from 'rxjs';
import { EnrollmentChoice } from '@/models/types';
import { allEnrollmentChoices$ } from '../observable';

export const useAllEnrollmentChoices = (eventId?: string) => {
    const [enrollmentChoices, setEnrollmentChoices] = useState<EnrollmentChoice[]>();
    const [loadingEnrollmentChoices, setLoadingEnrollmentChoices] = useState<boolean>(true);

    useEffect(() => {
        let allEnrollmentChoicesObservable: Subscription;
        if (eventId)
            allEnrollmentChoicesObservable = allEnrollmentChoices$(eventId).subscribe((enrollmentChoices) => {
                setEnrollmentChoices(enrollmentChoices as EnrollmentChoice[] | undefined);
                setLoadingEnrollmentChoices(false);

                return () => {
                    allEnrollmentChoicesObservable.unsubscribe();
                };
            });
        return () => {
            null;
        };
    }, []);
    return { enrollmentChoices, loadingEnrollmentChoices };
};
