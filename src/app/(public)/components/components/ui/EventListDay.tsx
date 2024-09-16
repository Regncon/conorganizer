'use client';

import type { PoolName } from '$lib/enums';
import { Typography } from '@mui/material';
import { getTranslatedDay } from '../../lib/helpers/translation';

type Props = {
    poolDay: PoolName;
};

const EventListDay = ({ poolDay }: Props) => {
    const translatedDay = getTranslatedDay(poolDay);
    return (
        <Typography
            id={translatedDay}
            sx={{ scrollMarginTop: 'calc(var(--scroll-margin-top) + var(--app-bar-height))' }}
            variant="h1"
        >
            {translatedDay}
        </Typography>
    );
};

export default EventListDay;
