'use client';

import { useEffect, useState } from 'react';
import { Box } from '@mui/material';
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
import { EnrollmentChoice } from '@/lib/enums';
import { db } from '@/lib/firebase';
import { useSingleEnrollment } from '@/lib/hooks/UseEnrollments';
import { ConEvent } from '@/models/types';
import { useAuth } from './AuthProvider';
import EnrollmentSelector from './EnrollmentSelector';
import EventHeader from './EventHeader';

type Props = {
    conEvent: ConEvent | undefined;
};

const EventUi = ({ conEvent }: Props) => {
    const user = useAuth();
    const { enrollments: enrollment } = useSingleEnrollment(conEvent?.id || '', user?.uid || '');
    const [errorMessage, setErrorMessage] = useState<string>();
    const [enrollmentChoice, setEnrollmentChoice] = useState<EnrollmentChoice>(EnrollmentChoice.NotInterested);

    const [description, setDescription] = useState('');
    useEffect(() => {
        if (conEvent) {
            const tmp: string = conEvent?.description;
            // tmp = tmp.replace(/\n/g, '</p><p>');
            setDescription(tmp);
        }
    }, [conEvent]);

    useEffect(() => {
         if (user && conEvent?.id, enrollment) {
            setEnrollmentChoice(enrollment?.choice || EnrollmentChoice.NotInterested);
        } 
    }, [user, conEvent, enrollment]);

    const handleEnrollmentChoiceChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const choice = Number(event.target.value) as EnrollmentChoice;
        setEnrollmentChoice(choice);
    };

    useEffect(() => {
        if (user && conEvent?.id) {
            setEnrollment();
        }
    }
    , [enrollmentChoice]);

    async function setEnrollment() {
        try {
            if (!user || !conEvent?.id) {
                return;
            }
            const setEnrollmentRef = doc(db, `events/${conEvent.id}`, `/enrollments/${user.uid}`);
            if (enrollment) {
                await updateDoc(setEnrollmentRef, {
                    choice: enrollmentChoice,
                });
            } else {
                await setDoc(setEnrollmentRef, {
                    choice: enrollmentChoice,
                });
            }
        } catch (e) {
            console.error(e);
            const error = e as Error;
            setErrorMessage(error.message);
        }
    }
    // throw new Error(
    //     'lorem Ipsum error in conAuthor authorization dialog box - invalid Lorem, ipsum dolor sit amet consectetur adipisicing elit. Ea quia in blanditiis mollitia exercitationem, asperiores nam quidem commodi nulla illum laborum, distinctio magnam debitis vitae rerum, maiores maxime sapiente! Quia! Lorem, ipsum dolor sit amet consectetur adipisicing elit. Ea quia in blanditiis mollitia exercitationem, asperiores nam quidem commodi nulla illum laborum, distinctio magnam debitis vitae rerum, maiores maxime sapiente! Quia!'
    // );
    return (
        <Card>
            <EventHeader conEvent={conEvent} />
            <Divider />
            <Box className="p-4" sx={{ minHeight: '7rem', display: 'grid', gap: '.5rem' }}>
                {parse(description || '')}
            </Box>

            <Divider />
            <CardContent>
                <FormControl>
                    <FormLabel id="demo-row-radio-buttons-group-label">
                        <Typography variant="h6">
                            Påmelding
                            {user ? '' : ' (Kjøp bilett og logg inn for å melde deg på)'}
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
    );
};

export default EventUi;
