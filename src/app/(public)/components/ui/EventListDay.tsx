'use client';
import type { EventDay, EventDays } from '$app/(public)/page';
import { Typography } from '@mui/material';
import { useEffect, useRef } from 'react';

type Props = {
    eventDay: EventDay;
    intersectionObserver: IntersectionObserver | null;
};

const EventListDay = ({ eventDay, intersectionObserver }: Props) => {
    const ref = useRef<HTMLDivElement>(null);
    // const intersectionObserver = new IntersectionObserver(
    //     (entries) => {
    //         entries.forEach((entry) => {
    //             console.log(entry);
    //         });
    //     },
    //     { root: ref.current, threshold: 1 }
    // );

    useEffect(() => {
        if (ref.current) {
            intersectionObserver?.observe(ref.current);
        }
        return () => {
            if (ref.current) {
                intersectionObserver?.unobserve(ref.current);
            }
        };
    }, [ref, ref.current]);

    return (
        <Typography
            id={eventDay}
            sx={{ scrollMarginTop: 'calc(var(--scroll-margin-top) + var(--app-bar-height))' }}
            variant="h1"
            ref={ref}
        >
            {eventDay}
        </Typography>
    );
};

export default EventListDay;
