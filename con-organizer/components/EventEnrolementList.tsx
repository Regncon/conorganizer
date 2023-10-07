'use client';

import { Box, Card, CardContent, CardHeader, Checkbox, Divider } from '@mui/material';
import { doc, updateDoc } from 'firebase/firestore';
import { db } from '@/lib/firebase';
import { useAllEnrollmentChoices } from '@/lib/hooks/UseAllEnrollmentChoices';
import { EnrollmentOptions, FirebaseCollections } from '@/models/enums';
import { EnrollmentChoice } from '@/models/types';
import EventEnrollment from './EventEnrollment';

type Props = { id: string };

const EventEnrolementList = ({ id }: Props) => {
    const { enrollmentChoices, loadingEnrollmentChoices } = useAllEnrollmentChoices(id || '');

    console.log(enrollmentChoices);

    const handleChoiceChange = (event: React.ChangeEvent<HTMLInputElement>, enrollmentChoice: EnrollmentChoice) => {
        console.log(event.target.checked, enrollmentChoice);
        updateEnrollmentChoiceOnUser(enrollmentChoice, event.target.checked);
    };

    async function updateEnrollmentChoiceOnUser(enrollmentChoice: EnrollmentChoice, enroll: boolean) {
        const setEnrollmentRef = doc(
            db,
            `${FirebaseCollections.events}/${enrollmentChoice.eventId}`,
            `/${FirebaseCollections.EnrollmentChoices}/${enrollmentChoice.participantId}`
        );
        try {
            await updateDoc(setEnrollmentRef, {
                isEnrolled: enroll,
                enrolledEventId: enroll ? enrollmentChoice.eventId : '',
                enrolledEventTitle: enroll ? enrollmentChoice.eventTitle : '',
            });
        } catch (e) {
            console.error(e);
        }
    }

    return (
        <Card>
            <CardHeader title='Veldig interessert' />
            <CardContent>
                {enrollmentChoices
                    ?.toSorted((a, b) => a.choice - b.choice)
                    .filter((enrollmentChoice) => enrollmentChoice.choice === EnrollmentOptions.VeryInterested)
                    .map((enrollmentChoice) => (
                        <EventEnrollment
                            key={enrollmentChoice.id}
                            enrollmentChoice={enrollmentChoice}
                            handleChoiceChange={handleChoiceChange}
                        />
                    ))}
            </CardContent>
        </Card>
    );
};

export default EventEnrolementList;
