import type { ConEvent } from '$lib/types';
import { Typography } from '@mui/material';

type Props = {
    day: string;
};

const ScrollToHeader = ({ day }: Props) => {
    return (
        <Typography id={day} sx={{ scrollMarginTop: 'var(--scroll-margin-top)' }} variant="h1">
            {day}
        </Typography>
    );
};

export default ScrollToHeader;
