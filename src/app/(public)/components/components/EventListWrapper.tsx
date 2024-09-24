'use client';
import { Box } from '@mui/material';
import { useRef, type PropsWithChildren } from 'react';
import { useObserveIntersectionObserver } from '../lib/hooks/useObserveIntersectionObserver';

type Props = {
    day: string;
};

const EventListWrapper = ({ day, children }: PropsWithChildren<Props>) => {
    const ref = useRef<HTMLDivElement>(null);
    useObserveIntersectionObserver(ref);
    return (
        <Box key={day} ref={ref}>
            {children}
        </Box>
    );
};

export default EventListWrapper;
