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
const RoomCard = ({ roomName, title, gameMaster, system, imageUri }: Props) => {
    const smallRoomRowX = 2460;
    const roomCoordinates = () => {
        switch (roomName) {
            case RoomName.Klang:
                return { x: smallRoomRowX, y: 450 };
            case RoomName.Sonate:
                return { x: smallRoomRowX, y: 640 };
            default:
                return { x: 0, y: 0 };
        }
    };
    return (
        <Box
            sx={{
                padding: '1rem',
                width: '306px',
                opacity: 0.9,
                background: `linear-gradient( rgba(0, 0, 0, 0.5), rgba(0, 0, 0, 0.5) ), url(${imageUri})`,
                backgroundSize: 'cover',
                borderRadius: '1.75rem',
                position: 'absolute',
                left: roomCoordinates().x,
                top: roomCoordinates().y,
            }}
        >
            <Typography>{title}</Typography>
            <Typography sx={{ fontWeight: 'bold', fontSize: '1.1rem' }}> {gameMaster} </Typography>
            <Typography sx={{ color: 'primary.main' }}>{system}</Typography>
        </Box>
    );
};
export default RoomCard;