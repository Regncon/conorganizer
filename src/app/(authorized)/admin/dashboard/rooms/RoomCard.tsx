import { Box, Typography } from '@mui/material';
import { PoolName, RoomName } from '$lib/enums';

type Props = {
    poolName: PoolName;
    roomName: RoomName;
    title: string;
    gameMaster: string;
    system: string;
    imageUri: string;
};
const RoomCard = ({ roomPosition, title, gameMaster, system, imageUri }: Props) => {
    return (
        <Box
            sx={{
                padding: '1rem',
                width: '306px',
                backgroundImage: `url(${imageUri})`,
                backgroundSize: 'cover',
                borderRadius: '1.75rem',
            }}
        >
            <Typography>{title}</Typography>
            <Typography sx={{ fontWeight: 'bold', fontSize: '1.1rem' }}> {gameMaster} </Typography>
            <Typography sx={{ color: 'primary.main' }}>{system}</Typography>
        </Box>
    );
};
export default RoomCard;
