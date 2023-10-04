'use client';

import { useEffect, useState } from 'react';
import { Alert, Dialog, FormControl, FormLabel, Link, Radio, RadioGroup, Typography } from '@mui/material';
import { doc, setDoc,updateDoc } from 'firebase/firestore';
import { db } from '@/lib/firebase';
import { useSingleEnrollment } from '@/lib/hooks/UseEnrollments';
import { EnrollmentChoice } from '@/models/enums';
import { ConEvent, Participant } from '@/models/types';
import { useAuth } from './AuthProvider';
import EnrollmentSelector from './EnrollmentSelector';
import Login from './Login';

type Props = {
    conEvent: ConEvent | undefined;
    participant: Participant | undefined;
};

const EventPreference = ({ conEvent, participant }: Props) => {
    const user = useAuth();
    const [errorMessage, setErrorMessage] = useState<string>();
    const [enrollmentChoice, setEnrollmentChoice] = useState<EnrollmentChoice>(EnrollmentChoice.NotInterested);
    const [openLogin, setOpenLogin] = useState(false);
    const { enrollments: enrollment } = useSingleEnrollment(conEvent?.id || '', user?.uid || '');

    useEffect(() => {
        setEnrollmentChoice(user && conEvent?.id && enrollment ? enrollment.choice : EnrollmentChoice.NotInterested);
    }, [user, conEvent, enrollment]);

    const handleEnrollmentChoiceChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const choice = Number(event.target.value) as EnrollmentChoice;

        if (user && conEvent?.id) {
            setEnrollmentChoice(choice);
            updateEnrollmentInDb(choice);
        }
    };

    async function updateEnrollmentInDb(choice: EnrollmentChoice) {
        console.log('updateEnrollmentInDb', choice);
        console.log('updateEnrollmentInDb', participant);
        try {
            if (!user || !conEvent?.id) {
                return;
            }
            const setEnrollmentRef = doc(db, `events/${conEvent.id}`, `/enrollments/${user.uid}/participants/${participant?.id}`);
            if (enrollment) {
                await updateDoc(setEnrollmentRef, {
                    choice: choice,
                });
            } else {
                await setDoc(setEnrollmentRef, {
                    choice: choice,
                });
            }
        } catch (e) {
            console.error(e);
            const error = e as Error;
            setErrorMessage(error.message);
        }
    }

    return (
        <>
            <FormControl>
                <FormLabel id="demo-row-radio-buttons-group-label">
                    {user ? (
                        <Typography variant="h6">Påmelding: {participant?.name}</Typography>
                    ) : (
                        <Alert severity="info">
                            <Link href="https://www.regncon.no/kjop-billett-til-regncon-xxxi/" color="secondary">
                                Kjøp billett
                            </Link>
                            <span> og </span>
                            <Link
                                sx={{ cursor: 'pointer' }}
                                onClick={() => {
                                    setOpenLogin(true);
                                }}
                                color="secondary"
                            >
                                logg inn
                            </Link>
                            <span> for å melde deg på.</span>
                        </Alert>
                    )}
                </FormLabel>
                <RadioGroup
                    row
                    aria-labelledby="demo-row-radio-buttons-group-label"
                    name="row-radio-buttons-group"
                    defaultValue={EnrollmentChoice.NotInterested}
                    value={enrollmentChoice}
                    sx={{
                        display: 'grid',
                        width: '100%',
                        maxWidth: '1080px',
                        padding: '.2em',
                        gridAutoFlow: 'column',
                        gridAutoColumns: '1fr',
                        placeContent: 'center',
                    }}
                    onChange={(e) => {
                        handleEnrollmentChoiceChange(e);
                    }}
                >
                    <EnrollmentSelector
                        sx={{ display: 'grid', textAlign: 'center', p: '.4em', m: '0' }}
                        value={EnrollmentChoice.NotInterested}
                        disabled={!user}
                        control={<Radio size="small" />}
                        label="Ikke p&aring;meldt"
                    />
                    <EnrollmentSelector
                        value={EnrollmentChoice.IfIHaveTo}
                        disabled={!user}
                        control={<Radio size="small" />}
                        sx={{
                            display: 'grid',
                            backgroundColor: '#00000055',
                            textAlign: 'center',
                            p: '.4em',
                            m: '0',
                        }}
                        label="Litt interessert"
                    />
                    <EnrollmentSelector
                        value={EnrollmentChoice.Interested}
                        disabled={!user}
                        sx={{
                            display: 'grid',
                            backgroundColor: '#00000088',
                            textAlign: 'center',
                            p: '.4em',
                            m: '0',
                        }}
                        control={<Radio size="small" />}
                        label="Ganske interessert"
                    />
                    <EnrollmentSelector
                        value={EnrollmentChoice.VeryInterested}
                        disabled={!user}
                        control={<Radio size="small" />}
                        label="Veldig interessert"
                        sx={{
                            display: 'grid',
                            backgroundColor: '#000000ff',
                            textAlign: 'center',
                            p: '.4em',
                            m: '0',
                        }}
                    />
                </RadioGroup>
            </FormControl>
            <Dialog open={openLogin}>
                <Login setChoice={() => setOpenLogin(false)} />
            </Dialog>
        </>
    );
};

export default EventPreference;
