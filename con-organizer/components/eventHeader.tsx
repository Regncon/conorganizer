import { faDiceD20 } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Box, CardHeader, CardMedia, Chip } from '@mui/material';
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
                        <Chip
                            icon={<FontAwesomeIcon icon={faDiceD20} />}
                            label="Rollespill"
                            size="small"
                            variant="outlined"
                        />
                    </span>
                    <span>DnD 5e </span> <span>Rom 222,</span>
                    <span>Søndag: 12:00 - 16:00</span>
                </Box>
            </Box>
        </>
    );
};

export default EventHeader;
