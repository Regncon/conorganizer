import { useEffect, useState } from 'react';
import { faChessKing, faChild, faClock, faDiceD20, faPalette, faUserPlus } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Alert, Box, Tooltip, Typography } from '@mui/material';
import { useSingleEnrollment } from '@/lib/hooks/UseEnrollments';
import { GameType } from '@/models/enums';
import { ConEvent } from '@/models/types';
import { useAuth } from './AuthProvider';
type Props = {
    conEvent: ConEvent | undefined;
    listView?: boolean;
};
const EventCardHeader = ({ conEvent }: Props) => {
    const user = useAuth();
    const { enrollments: enrollment } = useSingleEnrollment(conEvent?.id || '', user?.uid || '');
    const [childFriendly, setChildFriendly] = useState<boolean>(conEvent?.childFriendly || false);
    const [possiblyEnglish, setPossiblyEnglish] = useState(conEvent?.possiblyEnglish || false);
    const [adultsOnly, setAdultsOnly] = useState(conEvent?.adultsOnly || false);
    const [volunteersPossible, setVolunteersPossible] = useState(conEvent?.volunteersPossible || false);
    const [lessThanThreeHours, setLessThanThreeHours] = useState(conEvent?.lessThanThreeHours || false);
    const [moreThanSixHours, setMoreThanSixHours] = useState(conEvent?.moreThanSixHours || false);
    const [beginnerFriendly, setBeginnerFriendly] = useState(conEvent?.beginnerFriendly || false);

    useEffect(() => {
        setChildFriendly(conEvent?.childFriendly || false);
        setPossiblyEnglish(conEvent?.possiblyEnglish || false);
        setAdultsOnly(conEvent?.adultsOnly || false);
        setVolunteersPossible(conEvent?.volunteersPossible || false);
        setLessThanThreeHours(conEvent?.lessThanThreeHours || false);
        setMoreThanSixHours(conEvent?.moreThanSixHours || false);
        setBeginnerFriendly(conEvent?.beginnerFriendly || false);
    }, [conEvent]);
    return (
        <>
            {conEvent?.published === false ? (
                <Alert
                    severity="warning"
                    sx={{
                        marginBottom: '1rem',
                        position: 'absolute',
                        maxWidth: '500px',
                        margin: '.2em',
                        background: 'linear-gradient(to right, black, transparent)',
                    }}
                >
                    Dette arrangementet er ikke publisert enda.
                </Alert>
            ) : null}
            <Box
                sx={{
                    backgroundImage: `url(${
                        'image/programm/listview' + conEvent?.imageUrl || '/image/placeholder.jpg'
                    })`,
                    backgroundSize: 'cover',
                    backgroundBlendMode: 'multiply',
                    backgroundPositionY: 'bottom',
                    backgroundColor: '#ccc',
                    display: 'flex',
                    alignItems: 'space-between',
                    flexDirection: 'column',
                }}
            >
                <Box
                    sx={{
                        minHeight: '5em',
                        display: 'flex',
                        alignItems: 'start',
                    }}
                >
                    <Box
                        sx={{
                            color: 'white',
                            maxWidth: { xs: '100vw', md: '1080px' },
                            padding: '.7em',
                        }}
                    >
                        <Typography variant="h3" sx={{ textWrap: 'balance' }}>
                            {conEvent?.title}
                        </Typography>
                        <Typography variant="h4" sx={{ textWrap: 'balance' }}>
                            {conEvent?.subtitle}
                        </Typography>
                    </Box>
                </Box>
            </Box>
            <Box
                sx={{
                    display: 'grid',
                    gridAutoFlow: 'column',
                    justifyItems: 'start',
                    alignItems: 'center',
                    gridTemplateColumns: 'auto auto 1fr',
                    gap: '.5em',
                    color: 'black',
                    backgroundColor: 'white',
                    padding: '.3em .5em',
                }}
            >
                <span>
                    <Box sx={{ display: 'flex', gap: '.3em', placeItems: 'center' }}>
                        <Tooltip
                            arrow
                            sx={{ padding: '2em', backgroundColor: 'black' }}
                            title={
                                <Box
                                    sx={{
                                        display: 'grid',
                                        gridTemplateColumns: 'auto',
                                        alignItems: 'center',
                                        justifyContent: 'center',
                                        gap: '.3em',
                                    }}
                                >
                                    {conEvent?.gameType === GameType.roleplaying ? (
                                        <Typography>
                                            <FontAwesomeIcon icon={faDiceD20} fontSize="1em" color="orangered" /> :
                                            Rollespill
                                        </Typography>
                                    ) : null}
                                    {conEvent?.gameType === GameType.boardgame ? (
                                        <Typography>
                                            <FontAwesomeIcon icon={faChessKing} fontSize="1em" color="orangered" /> :
                                            Brettspill
                                        </Typography>
                                    ) : null}
                                    {conEvent?.gameType === GameType.other ? (
                                        <Typography>
                                            <FontAwesomeIcon icon={faPalette} fontSize="1em" color="orangered" /> :
                                            Annet
                                        </Typography>
                                    ) : null}
                                    {childFriendly && (
                                        <Typography>
                                            <FontAwesomeIcon icon={faChild} fontSize="1em" color="#00aaff" /> :
                                            Barnevennlig
                                        </Typography>
                                    )}
                                    {possiblyEnglish && <Typography>🇬🇧 : Can be run in English</Typography>}
                                    {adultsOnly && <Typography fontSize=".9rem">🔞 : Kun for voksne</Typography>}
                                    {beginnerFriendly && <Typography>😅 : Nybegynnervennlig</Typography>}
                                    {volunteersPossible && (
                                        <Typography>
                                            <FontAwesomeIcon icon={faUserPlus} fontSize=".8em" color="#55cc99" /> : Kan
                                            kjøres av frivillige
                                        </Typography>
                                    )}
                                    {moreThanSixHours && (
                                        <Typography color="white">
                                            <FontAwesomeIcon icon={faClock} fontSize=".8em" color="red" /> : Kan vare
                                            lengre enn seks timer
                                        </Typography>
                                    )}
                                    {lessThanThreeHours && (
                                        <Typography color="white">
                                            <FontAwesomeIcon icon={faClock} fontSize=".8em" color="lightgreen" /> : Kan
                                            vare mindre enn tre timer
                                        </Typography>
                                    )}
                                </Box>
                            }
                        >
                            <Box display="flex" gap=".4em" sx={{ placeItems: 'center' }}>
                                {conEvent?.gameType === GameType.roleplaying ? (
                                    <FontAwesomeIcon icon={faDiceD20} fontSize="1em" color="orangered" />
                                ) : null}
                                {conEvent?.gameType === GameType.boardgame ? (
                                    <FontAwesomeIcon icon={faChessKing} fontSize="1em" color="orangered" />
                                ) : null}
                                {conEvent?.gameType === GameType.other ? (
                                    <FontAwesomeIcon icon={faPalette} fontSize="1em" color="orangered" />
                                ) : null}
                                {childFriendly && <FontAwesomeIcon icon={faChild} fontSize="1em" color="#00aaff" />}
                                {possiblyEnglish && <Typography>🇬🇧</Typography>}
                                {adultsOnly && <Typography fontSize=".9rem">🔞</Typography>}
                                {beginnerFriendly && <Typography fontSize=".8em">😅</Typography>}
                                {volunteersPossible && (
                                    <FontAwesomeIcon icon={faUserPlus} fontSize=".8em" color="#55cc99" />
                                )}
                                {moreThanSixHours && (
                                    <Typography color="red">
                                        <FontAwesomeIcon icon={faClock} fontSize=".8em" />+
                                    </Typography>
                                )}
                                {lessThanThreeHours && (
                                    <Typography color="lightgreen">
                                        <FontAwesomeIcon icon={faClock} fontSize=".8em" />÷
                                    </Typography>
                                )}
                            </Box>
                        </Tooltip>
                        {/* <Typography variant="body1">Rollespill</Typography> */}
                    </Box>
                </span>
                <Typography variant="caption" fontWeight="700">
                    {conEvent?.room}
                </Typography>
                <Typography variant="caption" sx={{ fontStyle: 'italic' }}>
                    {conEvent?.host}
                </Typography>
                <div></div>
                {/* {!enrollment?.choice ? (
                    <Typography variant="caption" color="lightgray">
                        ⬤ Ikke p&aring;meldt
                    </Typography>
                ) : null}
                {enrollment?.choice === 1 ? (
                    <Typography variant="caption" color="darkgray">
                        ⬤ Litt interessert
                    </Typography>
                ) : null}
                {enrollment?.choice === 2 ? (
                    <Typography variant="caption" color="gray">
                        ⬤ Ganske interessert
                    </Typography>
                ) : null}
                {enrollment?.choice === 3 ? (
                    <Typography variant="caption" color="black">
                        ⬤ Veldig interessert
                    </Typography>
                ) : null}
 */}            </Box>
        </>
    );
};

export default EventCardHeader;
