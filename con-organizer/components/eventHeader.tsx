import { faChessKing, faDiceD20, faHatWizard } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Box, CardHeader, CardMedia, Chip } from '@mui/material';
import { gameType } from '@/lib/enums';
import { ConEvent } from '@/lib/types';

type Props = {
    conEvent: ConEvent;
};

const EventHeader = ({ conEvent }: Props) => {
    return (
        <>
            <CardHeader sx={{ paddingBottom: '0.5rem' }} title={conEvent?.title} subheader={conEvent?.subtitle} />

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
                        {conEvent?.gameType === gameType.roleplaying ? (
                            <Chip
                                icon={<FontAwesomeIcon icon={faDiceD20} />}
                                label="Rollespill"
                                size="small"
                                variant="outlined"
                            />
                        ) : null}
                        {conEvent?.gameType === gameType.boardgame ? (
                            <Chip
                                icon={<FontAwesomeIcon icon={faChessKing} />}
                                label="Brettspill"
                                size="small"
                                variant="outlined"
                            />
                        ) : null}
                        {conEvent?.gameType === gameType.other ? (
                            <Chip
                                icon={<FontAwesomeIcon icon={faHatWizard} />}
                                label="Annet"
                                size="small"
                                variant="outlined"
                            />
                        ) : null}
                    </span>
                    <span>{conEvent.gameSystem} </span>
                    <span>{conEvent.room} </span>
                    <span>{conEvent.pool} </span>
                </Box>
            </Box>
        </>
    );
};

export default EventHeader;
