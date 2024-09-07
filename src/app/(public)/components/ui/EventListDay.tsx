'use client';
import type { EventDay, EventDays } from '$app/(public)/page';
import { Typography } from '@mui/material';
import { use, useContext, useEffect, useRef } from 'react';
import { IntersectionObserverContext } from '../lib/IntersectionObserverContext';

type Props = {
    eventDay: EventDay;
};

const EventListDay = ({ eventDay }: Props) => {
    return (
        <Typography
            id={eventDay}
            sx={{ scrollMarginTop: 'calc(var(--scroll-margin-top) + var(--app-bar-height))' }}
            variant="h1"
        >
            {eventDay}
        </Typography>
    );
};

export default EventListDay;
