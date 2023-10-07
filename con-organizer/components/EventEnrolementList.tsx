'use client';

import { Box, Card, CardContent, Checkbox, Divider } from '@mui/material';
import { doc, updateDoc } from 'firebase/firestore';
import { db } from '@/lib/firebase';
import { useAllEnrollmentChoices } from '@/lib/hooks/UseAllEnrollmentChoices';
import { EnrollmentOptions, FirebaseCollections } from '@/models/enums';
import { EnrollmentChoice } from '@/models/types';

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
            <CardContent>
                {enrollmentChoices
                    ?.toSorted((a, b) => a.choice - b.choice)
                    .map((enrollmentChoice) => (
                        <Box
                            sx={{
                                display: 'flex',
                                flexDirection: 'row',
                                p: 2,
                                m: 1,
                                borderRadius: 1,
                                border: '1px solid #ccc',
                                width: '100%',
                            }}
                            key={enrollmentChoice.id}
                        >
                            <span>{EnrollmentOptions[enrollmentChoice.choice]}</span>
                            <span>{enrollmentChoice.name}</span>
                            <Checkbox
                                checked={enrollmentChoice.isEnrolled}
                                onChange={(event) => handleChoiceChange(event, enrollmentChoice)}
                            />
                        </Box>
                    ))}
            </CardContent>
        </Card>
    );
};

export default EventEnrolementList;
