'use client';

import { use, useEffect, useState } from 'react';
import { faChild, faClock, faUserPlus } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Box, Chip, Link } from '@mui/material';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import Divider from '@mui/material/Divider';
import Typography from '@mui/material/Typography';
import parse from 'html-react-parser';
import { useAllParticipants } from '@/lib/hooks/UseAllParticipants';
import { ConEvent, Participant } from '@/models/types';
import { useAuth } from './AuthProvider';
import EventHeader from './EventHeader';
import EventPreference from './EventPreference';

type Props = {
    conEvent: ConEvent | undefined;
};

const EventUi = ({ conEvent }: Props) => {
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

    const [participantList, setParticipantList] = useState<Participant[] | undefined>([]);

    const user = useAuth();
    const { participants, loadingParticipants } = useAllParticipants(user?.uid || '');

    console.log(participants, 'participants');
    useEffect(() => {
        if (participants) {
            setParticipantList(participants);
            console.log(participants, 'participants');
        }
        console.log(user, 'user');
    }, [user, participants]);

    return (
        <>
            <Card sx={{ width: '100%' }}>
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
                    {participantList ? (
                        participantList.map((participant) => (
                            <>
                                <EventPreference conEvent={conEvent} participant={participant} key={participant.id} />
                                <Divider />
                            </>
                        ))
                    ) : (
                        <EventPreference conEvent={conEvent} participant={undefined} />
                    )}
                </CardContent>
                <Divider />
            </Card>
        </>
    );
};

export default EventUi;
