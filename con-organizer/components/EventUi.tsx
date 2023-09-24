'use client';

import { useEffect, useState } from 'react';
import { Box, Dialog, Link } from '@mui/material';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import Divider from '@mui/material/Divider';
import FormControl from '@mui/material/FormControl';
import FormLabel from '@mui/material/FormLabel';
import Radio from '@mui/material/Radio';
import RadioGroup from '@mui/material/RadioGroup';
import Typography from '@mui/material/Typography';
import { doc, setDoc, updateDoc } from 'firebase/firestore';
import parse from 'html-react-parser';
import { db } from '@/lib/firebase';
import { useSingleEnrollment } from '@/lib/hooks/UseEnrollments';
import { EnrollmentChoice } from '@/models/enums';
import { ConEvent } from '@/models/types';
import { useAuth } from './AuthProvider';
import EnrollmentSelector from './EnrollmentSelector';
import EventHeader from './EventHeader';
import Login from './Login';

type Props = {
    conEvent: ConEvent | undefined;
};

const EventUi = ({ conEvent }: Props) => {
    const user = useAuth();
    const { enrollments: enrollment } = useSingleEnrollment(conEvent?.id || '', user?.uid || '');
    const [errorMessage, setErrorMessage] = useState<string>();
    const [enrollmentChoice, setEnrollmentChoice] = useState<EnrollmentChoice>(EnrollmentChoice.NotInterested);
    const [openLogin, setOpenLogin] = useState(false);

    const [description, setDescription] = useState('');
    useEffect(() => {
        if (conEvent) {
            const tmp: string = conEvent?.description.replaceAll('<p>&nbsp;', '');
            setDescription(tmp);
        }
    }, [conEvent]);

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
        try {
            if (!user || !conEvent?.id) {
                return;
            }
            const setEnrollmentRef = doc(db, `events/${conEvent.id}`, `/enrollments/${user.uid}`);
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
        <Card>
            <EventHeader conEvent={conEvent} />
            <Divider />
            <Box className="p-4" sx={{ minHeight: '7rem', display: 'grid', gap: '.5rem' }}>
                <Typography variant="body1">{parse(description || '')}</Typography>
                <span>
                    <em>Tidspunkt: {conEvent?.pool}</em>
                </span>
            </Box>

                <Divider />
                <CardContent>
                    <FormControl>
                        <FormLabel id="demo-row-radio-buttons-group-label">
                            <Typography variant="h6">
                                Påmelding
                                {user ? (
                                    ''
                                ) : (
                                    <>
                                        :{' '}
                                        <Link href="https://www.regncon.no/kjop-billett-til-regncon-xxxi/">
                                            Kjøp bilett
                                        </Link>
                                        <span> og </span>
                                        <Link
                                            component="button"
                                            onClick={() => {
                                                setOpenLogin(true);
                                            }}
                                        >
                                            logg inn
                                        </Link>
                                        <span> for å melde deg på.</span>
                                    </>
                                )}
                            </Typography>
                        </FormLabel>
                        <RadioGroup
                            row
                            aria-labelledby="demo-row-radio-buttons-group-label"
                            name="row-radio-buttons-group"
                            defaultValue={EnrollmentChoice.NotInterested}
                            value={enrollmentChoice}
                            sx={{
                                display: 'grid',
                                width: '100vw',
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
                                sx={{ display: 'grid', textAlign: 'center', p: '.4em' }}
                                value={EnrollmentChoice.NotInterested}
                                disabled={!user}
                                control={<Radio size="small" />}
                                label="Ikke interessert"
                            />
                            <EnrollmentSelector
                                value={EnrollmentChoice.IfIHaveTo}
                                disabled={!user}
                                sx={{ display: 'grid', backgroundColor: '#00000055', textAlign: 'center', p: '.4em' }}
                                control={<Radio size="small" />}
                                label="Hvis jeg må"
                            />
                            <EnrollmentSelector
                                value={EnrollmentChoice.Interested}
                                disabled={!user}
                                sx={{ display: 'grid', backgroundColor: '#000000aa', textAlign: 'center', p: '.4em' }}
                                control={<Radio size="small" />}
                                label="Har lyst"
                            />
                            <EnrollmentSelector
                                value={EnrollmentChoice.VeryInterested}
                                disabled={!user}
                                control={<Radio size="small" />}
                                label="Har veldig lyst"
                                sx={{ display: 'grid', backgroundColor: '#000000ff', textAlign: 'center', p: '.4em' }}
                            />
                        </RadioGroup>
                    </FormControl>
                </CardContent>
                <Divider />
            </Card>
            <Dialog open={openLogin}>
                <Login setChoice={ () => setOpenLogin(false) } />
            </Dialog>
        </>
    );
};

export default EventUi;
