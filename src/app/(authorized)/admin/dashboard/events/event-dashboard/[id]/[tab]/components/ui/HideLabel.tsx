import { Typography } from '@mui/material';
import type { PropsWithChildren } from 'react';

const HideLabel = ({ children }: PropsWithChildren) => {
    return (
        <Typography component="span" sx={{ display: { md: 'block', xs: 'none' }, fontSize: '12px' }}>
            {children}
        </Typography>
    );
};
export default HideLabel;
