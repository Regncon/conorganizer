'use client';

import type { PoolName } from '$lib/enums';
import { Typography, type SxProps } from '@mui/material';
import { getTranslatedDayAndTime } from '../../lib/helpers/translation';

type Props = {
    poolDay: PoolName;
    sx?: SxProps;
};

const EventListDay = ({ poolDay, sx }: Props) => {
    const translatedDay = getTranslatedDayAndTime(poolDay);
    return (
        <Typography
            id={translatedDay}
            sx={{ scrollMarginTop: 'calc(var(--scroll-margin-top) + var(--app-bar-height-desktop))', ...sx }}
            variant="h1"
        >
            {translatedDay}
        </Typography>
    );
};

export default EventListDay;
