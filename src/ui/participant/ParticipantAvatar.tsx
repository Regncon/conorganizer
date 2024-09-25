import * as React from 'react';
import Avatar from '@mui/material/Avatar';
import { Box, Typography } from '@mui/material';

type props = {
    firstName: string;
    lastName: string;
    small?: boolean;
};

function stringToColor(string: string) {
    let hash = 0;
    let i;

    /* eslint-disable no-bitwise */
    for (i = 0; i < string.length; i += 1) {
        hash = string.charCodeAt(i) + ((hash << 5) - hash);
    }

    let color = '#';
    const saturation = 0.65; // Saturation factor (0-1)
    const lightness = 0.8; // Lightness factor (0-1)
    for (let i = 0; i < 3; i++) {
        const baseValue = (hash >> (i * 8)) & 0xff;
        const newValue = Math.floor(baseValue * saturation + lightness * 255 * (1 - saturation));
        color += newValue.toString(16).padStart(2, '0').slice(-2);
    }

    return color;
}

function stringAvatar(name: string) {
    return {
        sx: {
            bgcolor: stringToColor(name),
        },
        children: `${name.split(' ')[0][0]}${name.split(' ')[1][0]}`,
    };
}
const ParticipantAvatar = ({ firstName, lastName, small }: props) => {
    const name = `${firstName} ${lastName}`;

    return (
        <Box sx={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
            <Avatar {...stringAvatar(name)} />
            <Typography>{small ? firstName : name}</Typography>
        </Box>
    );
};
export default ParticipantAvatar;
