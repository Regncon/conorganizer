'use client';

import { useEffect, useState } from 'react';
import { faChild, faClock, faUserPlus } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Alert, Box, Chip, Dialog, Link } from '@mui/material';
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
    const [childFriendly, setChildFriendly] = useState<boolean>(conEvent?.childFriendly || false);
    const [possiblyEnglish, setPossiblyEnglish] = useState(conEvent?.possiblyEnglish || false);
    const [adultsOnly, setAdultsOnly] = useState(conEvent?.adultsOnly || false);
    const [volunteersPossible, setVolunteersPossible] = useState(conEvent?.volunteersPossible || false);
    const [lessThanThreeHours, setLessThanThreeHours] = useState(conEvent?.lessThanThreeHours || false);
    const [moreThanSixHours, setMoreThanSixHours] = useState(conEvent?.moreThanSixHours || false);
    const [beginnerFriendly, setBeginnerFriendly] = useState(conEvent?.beginnerFriendly || false);
    const [description, setDescription] = useState('');
    useEffect(() => {
        if (conEvent && conEvent.description) {
            const tmp: string = conEvent.description.replaceAll('<p>&nbsp;', '');
            setDescription(tmp);
        }
        setChildFriendly(conEvent?.childFriendly || false);
        setPossiblyEnglish(conEvent?.possiblyEnglish || false);
        setAdultsOnly(conEvent?.adultsOnly || false);
        setVolunteersPossible(conEvent?.volunteersPossible || false);
        setLessThanThreeHours(conEvent?.lessThanThreeHours || false);
        setMoreThanSixHours(conEvent?.moreThanSixHours || false);
        setBeginnerFriendly(conEvent?.beginnerFriendly || false);
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
            <Card sx={{ background: 'linear-gradient(to left, black, #222, #222, black)', width: '100%' }}>
                <EventHeader conEvent={conEvent} />
                <Divider />
                <Box
                    className="p-4"
                    sx={{
                        minHeight: '7rem',
                        display: 'grid',
                        gap: '.5rem',
                        maxWidth: '45em',
                        margin: { xs: '0 auto 1em auto', md: '2em auto 3em auto' },
                    }}
                >
                    <Typography variant="body1">{parse(description || '')}</Typography>
                    <span>
                        <em>Tidspunkt: {conEvent?.pool}</em>
                    </span>
                    <Box display="flex" gap=".5em" mt="1em" flexWrap="wrap">
                        {childFriendly && (
                            <Chip
                                icon={
                                    <p>
                                        &nbsp;
                                        <FontAwesomeIcon icon={faChild} fontSize="1em" color="#00aaff" />
                                    </p>
                                }
                                label="Barnevennlig"
                                variant="outlined"
                            />
                        )}
                        {possiblyEnglish && (
                            <Chip
                                icon={<Typography>&nbsp;🇬🇧</Typography>}
                                label="Can be run in English"
                                variant="outlined"
                            />
                        )}
                        {adultsOnly && (
                            <Chip
                                icon={<Typography fontSize=".9rem">&nbsp;🔞</Typography>}
                                label="Kun for voksne"
                                variant="outlined"
                            />
                        )}
                        {beginnerFriendly && (
                            <Chip
                                icon={<Typography>&nbsp;😅</Typography>}
                                label="Nybegynnervennlig"
                                variant="outlined"
                            />
                        )}
                        {volunteersPossible && (
                            <Chip
                                icon={
                                    <p>
                                        &nbsp;
                                        <FontAwesomeIcon icon={faUserPlus} fontSize=".8em" color="#55cc99" />
                                    </p>
                                }
                                label={
                                    <Typography variant="caption">
                                        Kan kjøres av frivillige.{' '}
                                        <Link href="mailto:arrangere@regncon.no" color="secondary">
                                            Kontakt styret for å melde deg!
                                        </Link>
                                    </Typography>
                                }
                                variant="outlined"
                            />
                        )}
                        {moreThanSixHours && (
                            <Chip
                                icon={
                                    <Typography color="red">
                                        &nbsp;
                                        <FontAwesomeIcon icon={faClock} fontSize=".8em" />+
                                    </Typography>
                                }
                                label="Kan vare lengre enn seks timer"
                                variant="outlined"
                            />
                        )}
                        {lessThanThreeHours && (
                            <Chip
                                icon={
                                    <Typography color="lightgreen">
                                        &nbsp;
                                        <FontAwesomeIcon icon={faClock} fontSize=".8em" />÷
                                    </Typography>
                                }
                                label="Kan vare mindre enn tre timer"
                                variant="outlined"
                            />
                        )}
                    </Box>
                </Box>

                <Divider />
                <CardContent
                    sx={
                        conEvent?.hideEnrollment
                            ? { display: 'none' }
                            : { backgroundColor: '#181818', borderRadius: '0', width: '100%' }
                    }
                >
                    <FormControl>
                        <FormLabel id="demo-row-radio-buttons-group-label">
                            {user ? (
                                <Typography variant="h6">Påmelding</Typography>
                            ) : (
                                <Alert severity="info">
                                    <Link
                                        href="https://www.regncon.no/kjop-billett-til-regncon-xxxi/"
                                        color="secondary"
                                    >
                                        Kjøp billett
                                    </Link>
                                    <span> og </span>
                                    <Link
                                        component="button"
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
                </CardContent>
                <Divider />
            </Card>
            <Dialog open={openLogin}>
                <Login setChoice={() => setOpenLogin(false)} />
            </Dialog>
        </>
    );
};

export default EventUi;
