import { faChessKing, faDiceD20, faPalette } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Alert, Box, Typography } from '@mui/material';
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
                    backgroundBlendMode: 'multiply',
                    backgroundColor: '#999',
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
                    padding: '.3em .5em',
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
                {!enrollment?.choice ? (
                    <Typography variant="caption" color="lightgray">
                        ⬤ Ikke p&aring;meldt
                    </Typography>
                ) : null}
                {enrollment?.choice === 1 ? (
                    <Typography variant="caption" color="gray">
                        ⬤ Litt interessert
                    </Typography>
                ) : null}
                {enrollment?.choice === 2 ? (
                    <Typography variant="caption" color="darkgray">
                        ⬤ Ganske interessert
                    </Typography>
                ) : null}
                {enrollment?.choice === 3 ? (
                    <Typography variant="caption" color="black">
                        ⬤ Veldig interessert
                    </Typography>
                ) : null}
            </Box>
        </>
    );
};

export default EventCardHeader;
