import { faChessKing, faDiceD20, faPalette } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Alert, Box, CardHeader, CardMedia, Chip } from '@mui/material';
import { GameType } from '@/lib/enums';
import { ConEvent } from '@/models/types';

type Props = {
    conEvent: ConEvent | undefined;
};

const EventHeader = ({ conEvent }: Props) => {
    return (
        <>
            {conEvent?.published === false ? (
                <Alert severity="warning" sx={{ marginBottom: '1rem' }}>
                    Dette arrangementet er ikke publisert enda.
                </Alert>
            ) : null}
            <CardHeader
                title={conEvent?.title}
                subheader={conEvent?.subtitle}
                sx={{
                    backgroundImage: 'url(/placeholder.jpg)',
                    pb: '.5rem',
                    height: '50vh',
                    minHeight: '300px',
                    maxHeight: '500px',
                    backgroundSize: 'cover',
                    color: 'white',
                    alignItems: 'end',
                }}
            />

            <Box className="flex justify-start pb-4">
                <CardMedia
                    className="ml-4"
                    sx={{ width: '40%', maxHeight: '130px' }}
                    component="img"
                    image="/placeholder.jpg"
                    alt={conEvent?.title}
                />
                <Box className="flex flex-col pl-4 pr-4">
                    <span>
                        {conEvent?.gameType === GameType.roleplaying ? (
                            <Chip
                                icon={<FontAwesomeIcon icon={faDiceD20} />}
                                label="Rollespill"
                                size="small"
                                variant="outlined"
                            />
                        ) : null}
                        {conEvent?.gameType === GameType.boardgame ? (
                            <Chip
                                icon={<FontAwesomeIcon icon={faChessKing} />}
                                label="Brettspill"
                                size="small"
                                variant="outlined"
                            />
                        ) : null}
                        {conEvent?.gameType === GameType.other ? (
                            <Chip
                                icon={<FontAwesomeIcon icon={faPalette} />}
                                label="Annet"
                                size="small"
                                variant="outlined"
                            />
                        ) : null}
                    </span>
                    <span>{conEvent?.gameSystem} </span>
                    <span>{conEvent?.room} </span>
                    <span>{conEvent?.pool} </span>
                </Box>
            </Box>
        </>
    );
};

export default EventHeader;
