import { faChessKing, faDiceD20, faPalette } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Alert, Box, CardMedia, Typography } from '@mui/material';
import { useSingleEnrollment } from '@/lib/hooks/UseEnrollments';
import { GameType } from '@/models/enums';
import { ConEvent } from '@/models/types';
import { useAuth } from './AuthProvider';
type Props = {
    conEvent: ConEvent | undefined;
};
const EventHeader = ({ conEvent }: Props) => {
    const user = useAuth();
    const { enrollments: enrollment } = useSingleEnrollment(conEvent?.id || '', user?.uid || '');
    const image = conEvent?.imageUrl;

    return (
        <>
            <Box>
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
                        background: 'linear-gradient(white, white, #333)',
                        display: 'grid',
                        alignItems: 'end',
                    }}
                >
                    <CardMedia component="img" image={image} sx={{ maxHeight: '50vh', width: '100%' }} />
                    <Box
                        sx={{
                            color: 'white',
                            maxWidth: { xs: '100vw', md: '1080px' },
                            padding: '.7em',
                            position: 'absolute',
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

                <Box
                    sx={{
                        display: 'grid',
                        justifyItems: 'start',
                        alignItems: 'center',
                        gridTemplateColumns: 'auto 1fr auto',
                        gap: '1em',
                        backgroundColor: 'white',
                        color: 'black',
                        padding: '1em',
                    }}
                >
                    <span>
                        {conEvent?.gameType === GameType.roleplaying ? (
                            <Box sx={{ placeItems: 'center', display: 'grid', gap: '.5em' }}>
                                <FontAwesomeIcon icon={faDiceD20} fontSize="2em" color="orangered" />
                                <Typography variant="caption">Rollespill</Typography>
                            </Box>
                        ) : null}
                        {conEvent?.gameType === GameType.boardgame ? (
                            <Box sx={{ placeItems: 'center', display: 'grid', gap: '.5em' }}>
                                <FontAwesomeIcon icon={faChessKing} fontSize="2em" color="orangered" />
                                <Typography variant="caption">Brettspill</Typography>
                            </Box>
                        ) : null}
                        {conEvent?.gameType === GameType.other ? (
                            <Box sx={{ placeItems: 'center', display: 'grid', gap: '.5em' }}>
                                <FontAwesomeIcon icon={faPalette} fontSize="2em" color="orangered" />
                                <Typography variant="caption">Annet</Typography>
                            </Box>
                        ) : null}
                    </span>
                    <Box sx={{ display: 'grid' }}>
                        <span>{conEvent?.gameSystem} </span>
                        <span>{conEvent?.room} </span>
                        <span>{conEvent?.host} </span>
                        {/* <span>{EnrollmentChoice[enrollment?.choice ?? 0]} </span> */}
                    </Box>
                    <Box marginLeft="auto">
 {/*                        {!enrollment?.choice ? (
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
                        ) : null} */}
                    </Box>
                </Box>
            </Box>
        </>
    );
};

export default EventHeader;
