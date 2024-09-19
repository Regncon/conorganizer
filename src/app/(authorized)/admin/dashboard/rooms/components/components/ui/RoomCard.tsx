import { Box, Typography } from '@mui/material';

type Props = {
    title: string;
    gameMaster: string;
    system: string;
    imageUri: string;
};
const RoomCard = ({ title, gameMaster, system, imageUri }: Props) => {
    return (
        <Box
            sx={{
                padding: '1rem',
                width: '306px',
                opacity: 0.9,
                background: `linear-gradient( rgba(0, 0, 0, 0.5), rgba(0, 0, 0, 0.5) ), url(${imageUri})`,
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
