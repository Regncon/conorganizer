import * as React from 'react';
import Avatar from '@mui/material/Avatar';
import { Box, Typography } from '@mui/material';

type props = {
    name: string;
};

const ParticipantAvatar = ({ name }: props) => {
    function stringToColor(string: string) {
        let hash = 0;
        let i;

        /* eslint-disable no-bitwise */
        for (i = 0; i < string.length; i += 1) {
            hash = string.charCodeAt(i) + ((hash << 5) - hash);
        }

        let color = '#';

        for (i = 0; i < 3; i += 1) {
            const value = (hash >> (i * 8)) & 0xff;
            color += `00${value.toString(16)}`.slice(-2);
        }
        /* eslint-enable no-bitwise */

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

    return (
        <Box sx={{ display: 'flex', alignItems: 'center' }}>
            <Avatar {...stringAvatar(name)} />
            <Typography>{name}</Typography>
        </Box>
    );
};
export default ParticipantAvatar;
