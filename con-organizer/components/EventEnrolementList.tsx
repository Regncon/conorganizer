'use client';

import { useEffect, useState } from 'react';
import {
    Box,
    Card,
    CardContent,
    CardHeader,
    Divider,
    Paper,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableRow,
} from '@mui/material';
import { doc, updateDoc } from 'firebase/firestore';
import { db } from '@/lib/firebase';
import { useAllEnrollmentChoices } from '@/lib/hooks/UseAllEnrollmentChoices';
import { EnrollmentOptions, FirebaseCollections } from '@/models/enums';
import { EnrollmentChoice } from '@/models/types';
import EventEnrollment from './EventEnrollment';

type Props = { id: string };

const EventEnrolementList = ({ id }: Props) => {
    const { enrollmentChoices, loadingEnrollmentChoices } = useAllEnrollmentChoices(id || '');

    const [enrollmentChoicesState, setEnrollmentChoicesState] = useState<EnrollmentChoice[]>([]);

    useEffect(() => {
        if (enrollmentChoices) {
            setEnrollmentChoicesState(enrollmentChoices);
            console.log(enrollmentChoices);
        }
    }, [enrollmentChoices]);

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
        <Box
        sx={{
            display: 'flex',
            flexDirection: 'column',
            gap: '1rem',
            
            
        }}
        >
            <Card>
                <CardHeader title="Veldig interessert" />
                <CardContent>
                    <TableContainer component={Paper}>
                        <Table sx={{ minWidth: 650 }} aria-label="simple table">
                            <TableBody>
                                {enrollmentChoicesState
                                    ?.toSorted((a, b) => a.choice - b.choice)
                                    .filter(
                                        (enrollmentChoice) =>
                                            enrollmentChoice.choice === EnrollmentOptions.VeryInterested
                                    )
                                    .map((enrollmentChoice) => (
                                        <EventEnrollment
                                            key={enrollmentChoice.id}
                                            enrollmentChoice={enrollmentChoice}
                                            handleChoiceChange={handleChoiceChange}
                                        />
                                    ))}
                            </TableBody>
                        </Table>
                    </TableContainer>
                </CardContent>
            </Card>
            <Card>
                <CardHeader title="Ganske interessert" />
                <CardContent>
                    <TableContainer component={Paper}>
                        <Table sx={{ minWidth: 650 }} aria-label="simple table">
                            <TableBody>
                                {enrollmentChoicesState
                                    ?.toSorted((a, b) => a.choice - b.choice)
                                    .filter(
                                        (enrollmentChoice) => enrollmentChoice.choice === EnrollmentOptions.Interested
                                    )
                                    .map((enrollmentChoice) => (
                                        <EventEnrollment
                                            key={enrollmentChoice.id}
                                            enrollmentChoice={enrollmentChoice}
                                            handleChoiceChange={handleChoiceChange}
                                        />
                                    ))}
                            </TableBody>
                        </Table>
                    </TableContainer>
                </CardContent>
            </Card>
            <Card>
                <CardHeader title="Litt interessert" />
                <CardContent>
                    <TableContainer component={Paper}>
                        <Table sx={{ minWidth: 650 }} aria-label="simple table">
                            <TableBody>
                                {enrollmentChoicesState
                                    ?.toSorted((a, b) => a.choice - b.choice)
                                    .filter(
                                        (enrollmentChoice) => enrollmentChoice.choice === EnrollmentOptions.IfIHaveTo
                                    )
                                    .map((enrollmentChoice) => (
                                        <EventEnrollment
                                            key={enrollmentChoice.id}
                                            enrollmentChoice={enrollmentChoice}
                                            handleChoiceChange={handleChoiceChange}
                                        />
                                    ))}
                            </TableBody>
                        </Table>
                    </TableContainer>
                </CardContent>
            </Card>
            <Card>
                <CardHeader title="Ikke interessert" />
                <CardContent>
                    <TableContainer component={Paper}>
                        <Table sx={{ minWidth: 650 }} aria-label="simple table">
                            <TableBody>
                                {enrollmentChoicesState
                                    ?.toSorted((a, b) => a.choice - b.choice)
                                    .filter(
                                        (enrollmentChoice) =>
                                            enrollmentChoice.choice === EnrollmentOptions.NotInterested
                                    )
                                    .map((enrollmentChoice) => (
                                        <EventEnrollment
                                            key={enrollmentChoice.id}
                                            enrollmentChoice={enrollmentChoice}
                                            handleChoiceChange={handleChoiceChange}
                                        />
                                    ))}
                            </TableBody>
                        </Table>
                    </TableContainer>
                </CardContent>
            </Card>
        </Box>
    );
};

export default EventEnrolementList;
