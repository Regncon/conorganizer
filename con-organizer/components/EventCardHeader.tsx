import { faChessKing, faDiceD20, faPalette } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Alert, Box, Typography } from '@mui/material';
import { useSingleEnrollment } from '@/lib/hooks/UseEnrollments';
import { EnrollmentChoice, GameType } from '@/models/enums';
import { ConEvent } from '@/models/types';
import { useAuth } from './AuthProvider';
type Props = {
    conEvent: ConEvent | undefined;
    listView?: boolean;
};
const EventCardHeader = ({ conEvent, listView = false }: Props) => {
    const user = useAuth();
    const { enrollments: enrollment } = useSingleEnrollment(conEvent?.id || '', user?.uid || '');

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
                    backgroundImage: `url(${conEvent?.imageUrl || '/image/placeholder.jpg'})`,
                    backgroundSize: 'cover',
                }}
            >
                <Box
                    sx={{
                        background: 'linear-gradient(#00000099, #00000066, transparent)',
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
                    padding: '.5em',
                }}
            >
                <span>
                    {conEvent?.gameType === GameType.roleplaying ? (
                        <Box sx={{ display: 'flex', gap: '.3em', placeItems: 'center' }}>
                            <FontAwesomeIcon icon={faDiceD20} fontSize="1em" color="orangered" />
                            <Typography variant="body1">Rollespill</Typography>
                        </Box>
                    ) : null}
                    {conEvent?.gameType === GameType.boardgame ? (
                        <Box sx={{ display: 'flex', gap: '.3em', placeItems: 'center' }}>
                            <FontAwesomeIcon icon={faChessKing} fontSize="1em" color="orangered" />
                            <Typography variant="body1">Brettspill</Typography>
                        </Box>
                    ) : null}
                    {conEvent?.gameType === GameType.other ? (
                        <Box sx={{ display: 'flex', gap: '.3em', placeItems: 'center' }}>
                            <FontAwesomeIcon icon={faPalette} fontSize="1em" color="orangered" />
                            <Typography variant="body1">Annet</Typography>
                        </Box>
                    ) : null}
                </span>
                {/* <span>{conEvent?.gameSystem} </span> */}
                <span>{conEvent?.room} </span>
                <span>{conEvent?.host} </span>
                <div></div>
                <span>{EnrollmentChoice[enrollment?.choice ?? 0]} </span>
            </Box>
        </>
    );
};

export default EventCardHeader;
